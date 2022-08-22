package thek

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

// StationController holds the station array
type StationController struct {
  Stations []*Station
}

// Station contains informations about the station
type Station struct {
  Name string
  StreamURL string
  EpgEntries []*EpgEntry
}

// EpgEntry contains informations about a specific program
type EpgEntry struct {
  Name string
  SubTitle string
  StartTime time.Time
  EndTime time.Time
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
func (sc *StationController) Tidy() error {
  now := time.Now()

  for _, sm := range sc.Stations {
    cleanedEntries := []*EpgEntry{}

    for _, epg := range sm.EpgEntries {
      if now.Before(epg.EndTime) {
        cleanedEntries = append(cleanedEntries, epg)
      }
    }

    sm.EpgEntries = cleanedEntries
  }

  return nil
}

// UpdateStationsByDates updates the stations with data from the given dates
func (sc *StationController) UpdateStationsByDates(date []time.Time) error {
  for _, fetchDate := range date {
    doc, err := FetchDocument(fmt.Sprintf(EpgURL, fetchDate.Format(DateFormat)))
    if err != nil {
      log.Fatal(err)
    }

    doc.Find("section.b-epg-timeline").Each(func(i int, s *goquery.Selection) {
      currentStation := &Station{}

      stationName := strings.Fields(s.Find("h3").First().Text())[0]

      // Check if we know this station already
      for _, station := range sc.Stations {
        if station.Name == stationName {
          currentStation = station
          break
        }
      }

      if currentStation.Name == "" {
        lowerStationName := strings.ToLower(stationName)

        _, ok := CurrentConfig.StationURLS[lowerStationName]
        if !ok {
          log.Warnf("Can't find stream url for %s in config. Have to skip this one...\n", lowerStationName)
          return
        }

        currentStation = &Station{
          Name: stationName,
          StreamURL: CurrentConfig.StationURLS[lowerStationName],
          EpgEntries: []*EpgEntry{},
        }

        sc.Stations = append(sc.Stations, currentStation)
      }

      newEpgEntries, err := parseEPG(fetchDate, s)
      if err != nil {
        log.Errorf("An error occured while parsing the epg data for %s: %s\n", stationName, err)
        return
      }

      for _, newEntry := range newEpgEntries {
        currentStation.EpgEntries = append(currentStation.EpgEntries, newEntry)
      }
    })
  }

  return nil
}

func setTimeToDate(t string, d time.Time) (time.Time, error) {
    hhmm := strings.Split(t, ":")
    timeHour, err := strconv.Atoi(hhmm[0])
    if err != nil {
      return d, err
    }
    timeMinute, err := strconv.Atoi(hhmm[1])
    if err != nil {
      return d, err
    }
    return time.Date(d.Year(), d.Month(), d.Day(), timeHour, timeMinute, 0, 0, d.Location()), nil
}

func parseEPG(date time.Time, s *goquery.Selection) ([]*EpgEntry, error) {
  var entries []*EpgEntry

  s.Find("li").Each(func(i int, v *goquery.Selection) {
    timeFields := strings.Fields(v.Find("span.time").First().Text())
    if len(timeFields) != 3 {
      return
    }
    startTimeString, endTimeString := timeFields[0], timeFields[2]

    startTime, err := setTimeToDate(startTimeString, date)
    if err != nil {
      return
    }

    endTime, err := setTimeToDate(endTimeString, date)
    if err != nil {
      return
    }

    videoName, found := v.Find("a").First().Attr("aria-label")
    if !found {
      return
    }

    dialog, found := v.Find("a").First().Attr("data-dialog")
    if !found {
      return
    }

    var dialogData map[string]interface{}
    err = json.Unmarshal([]byte(dialog), &dialogData)
    if err != nil {
      return
    }

    doc, err := FetchDocument(fmt.Sprintf(EpgDescriptionURL, dialogData["contentUrl"]))
    if err != nil {
      log.Fatal(err)
    }

    subtitle := doc.Find("h3.overlay-subtitle").Text()

    entry := &EpgEntry{
      Name: videoName,
      SubTitle: subtitle,
      StartTime: startTime,
      EndTime: endTime,
    }
    entries = append(entries, entry)
  })

  return entries, nil
}
