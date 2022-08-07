package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
  // EpgURL contains the query url for the epg informations
  // %s is the placeholder for the date in the format YYYY-MM-dd
  EpgURL = "https://www.zdf.de/live-tv?airtimeDate=%s"
)

// Sender contains informations about the sender
type Sender struct {
  Name string `json:"name"`
  Program []EpgData `json:"epg"`
}

// EpgData contains informations about the scheduled program
type EpgData struct {
  Date string `json:"date"`
  Entries []EpgEntry `json:"entries"`
}

// EpgEntry contains informations about a specific program
type EpgEntry struct {
  StartTime string `json:"start"`
  EndTime string `json:"end"`
  Name string `json:"name"`
}

// LoadData fetches the data for a given date
func LoadData(date string) ([]Sender, error) {
  res, err := http.Get(fmt.Sprintf(EpgURL, date))
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()

  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

  var result []Sender

  doc.Find("section.b-epg-timeline").Each(func(i int, s *goquery.Selection) {
    senderTitle := strings.Fields(s.Find("h3").First().Text())[0]
    epg, err := parseEPG(s, date)
    if err != nil {
      log.Fatal(err)
    }
    sender := Sender{
      Name: senderTitle,
      Program: []EpgData{epg},
    }
    result = append(result, sender)
  })

  return result, nil
}

func parseEPG(s *goquery.Selection, date string) (EpgData, error) {
  var entries []EpgEntry

  s.Find("li").Each(func(i int, v *goquery.Selection) {
    time := strings.Fields(v.Find("span.time").First().Text())
    if len(time) != 3 {
      return
    }
    startTime, endTime := time[0], time[2]

    videoName, found := v.Find("a").First().Attr("aria-label")
    if !found {
      return
    }

    entry := EpgEntry{
      StartTime: startTime,
      EndTime: endTime,
      Name: videoName,
    }
    entries = append(entries, entry)
  })

  result := EpgData{
    Date: date,
    Entries: entries,
  }

  return result, nil
}

func main() {

  s, err := LoadData("2022-08-07")
  if err != nil {
    log.Fatal(err)
  }

  for _, v := range s {
    fmt.Printf("Sender: %s\n", v.Name)
    for _, d := range v.Program {
      fmt.Printf("Date: %s\n", d.Date)
      for _, x := range d.Entries {
        fmt.Printf("Program: %s\n", x.Name)
      }
    }
  }

}
