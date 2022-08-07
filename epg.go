package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/PuerkitoBio/goquery"
)

const (
  // EpgURL contains the query url for the epg informations
  // %s is the placeholder for the date in the format YYYY-MM-dd
  EpgURL = "https://www.zdf.de/live-tv?airtimeDate=%s"
  // EpgDescriptionURL is the base path to the show description
  EpgDescriptionURL = "https://www.zdf.de%s"
)

// StationMap maps the stations to the coresponding data
type StationMap struct {
  Stations map[string]StationData
  Dates map[string]bool
}

// StationData contains informations about the station
type StationData struct {
  StreamURL string
  EpgByDays map[string] []EpgEntry
}

// EpgEntry contains informations about a specific program
type EpgEntry struct {
  Name string
  SubTitle string
  StartTime string
  EndTime string
}

// FetchDocument fetches the website and returns a goquery Document
func FetchDocument(url string) (*goquery.Document, error) {
    res, err := http.Get(url)
    if err != nil {
      log.Fatal(err)
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
      return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
    }

    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
      log.Fatal(err)
    }
    return doc, nil
}

// Tidy cleans up old data
func (sm *StationMap) Tidy() error {
  now := time.Now()
  del := []string{}
  for date := range sm.Dates {
    then, err := time.Parse(DateFormat, date)
    if err != nil {
      return err
    }
    if now.After(then.Add(time.Hour * 24.0)) {
      del = append(del, date)
    }
  }

  for _, d := range del {
    delete(sm.Dates, d)
    for _, stationData := range sm.Stations {
      delete(stationData.EpgByDays, d)
    }
  }
  return nil
}

// LoadEpgByDate fetches the data for a given date
func (sm *StationMap) LoadEpgByDate(date []string) error {
  for _, currentDate := range date {
    _, alreadyHasData := sm.Dates[currentDate]
    if alreadyHasData {
      continue
    }
    sm.Dates[currentDate] = true

    doc, err := FetchDocument(fmt.Sprintf(EpgURL, currentDate))
    if err != nil {
      log.Fatal(err)
    }

    doc.Find("section.b-epg-timeline").Each(func(i int, s *goquery.Selection) {
      stationName := strings.Fields(s.Find("h3").First().Text())[0]
      lowerStationName := strings.ToLower(stationName)

      _, ok := config.StationURLS[lowerStationName]
      if !ok {
        log.Warnf("Can't find stream url for %s in config. Have to skip this one...\n", lowerStationName)
        return
      }

      epg, err := parseEPG(s)
      if err != nil {
        log.Fatal(err)
      }

      stationData, ok := sm.Stations[stationName]
      if !ok {
        stationData = StationData{
          StreamURL: config.StationURLS[lowerStationName],
          EpgByDays: map[string] []EpgEntry{},
        }
        sm.Stations[stationName] = stationData
      }

      stationData.EpgByDays[currentDate] = epg
    })
  }
  return nil
}

func parseEPG(s *goquery.Selection) ([]EpgEntry, error) {
  var entries []EpgEntry

  s.Find("li").Each(func(i int, v *goquery.Selection) {
    timeFields := strings.Fields(v.Find("span.time").First().Text())
    if len(timeFields) != 3 {
      return
    }
    startTime, endTime := timeFields[0], timeFields[2]

    videoName, found := v.Find("a").First().Attr("aria-label")
    if !found {
      return
    }

    dialog, found := v.Find("a").First().Attr("data-dialog")
    if !found {
      return
    }

    var dialogData map[string]interface{}
    err := json.Unmarshal([]byte(dialog), &dialogData)
    if err != nil {
      return
    }

    doc, err := FetchDocument(fmt.Sprintf(EpgDescriptionURL, dialogData["contentUrl"]))
    if err != nil {
      log.Fatal(err)
    }

    subtitle := doc.Find("h3.overlay-subtitle").Text()

    entry := EpgEntry{
      SubTitle: subtitle,
      StartTime: startTime,
      EndTime: endTime,
      Name: videoName,
    }
    entries = append(entries, entry)
  })

  return entries, nil
}
