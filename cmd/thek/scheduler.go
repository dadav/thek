package thek

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gosimple/slug"
)

// DateFormat is the used format...douh
var DateFormat = "2006-01-02"

// WatchStreams starts the daemon
func WatchStreams() error {
  var epgTime time.Time
  var err error

  alreadyChecked := make(map[string]bool)
  controller := &StationController{}

  log.Info("Starting the main loop")

  for {
    now := time.Now()

    if epgTime.Day() != now.Day() {
      tomorrow := now.Add(time.Hour * 24)
      log.Infof("Updating epg data for %s and %s\n", now.Format(DateFormat), tomorrow.Format(DateFormat))
      err = controller.UpdateStationsByDates([]time.Time{now, tomorrow})
      if err != nil {
        return err
      }
      epgTime = now
      log.Info("Cleaning up some old data (if there is any)")
      controller.Tidy()
    }

    // Check if we need to schedule some recording for today
    for _, req := range config.RecordingTasks {
      for _, station := range controller.Stations {
        lowerReqStationName := strings.ToLower(req.Station)
        lowerStationName := strings.ToLower(station.Name)

        if lowerReqStationName != "" && (lowerReqStationName != lowerStationName) {
          continue
        }

          for _, show := range station.EpgEntries {
            lowerShowName := strings.ToLower(show.Name)
            lowShowSubtitle := strings.ToLower(show.SubTitle)
            lowerReqShowKeywords := strings.ToLower(req.ShowKeywords)

            if strings.Contains(lowerShowName, lowerReqShowKeywords) || strings.Contains(lowShowSubtitle, lowerReqShowKeywords) {
              uniqueKey := fmt.Sprintf("%s_-_%s", show.Name, show.SubTitle)
              if _, ok := alreadyChecked[uniqueKey]; !ok {
                alreadyChecked[uniqueKey] = true

                safetyStr := config.Defaults.DefaultSafetyDuration
                if req.SafetyDuration != "" {
                  safetyStr = req.SafetyDuration
                }
                safety, err := time.ParseDuration(safetyStr)
                if err != nil {
                  return err
                }

                showStart := show.StartTime.Add(-safety)
                showEnd := show.EndTime.Add(safety)
                showDuration := showEnd.Sub(showStart)

                outDir := config.Defaults.DefaultOutputDir
                if req.OutputDir != "" {
                  outDir = req.OutputDir
                }
                err = os.MkdirAll(outDir, os.ModePerm)
                if err != nil {
                  return err
                }

                out := filepath.Join(outDir, fmt.Sprintf("%s_-_%s.mkv", slug.Make(show.Name), slug.Make(show.SubTitle)))
                if _, err := os.Stat(out); err == nil {
                  existAction := config.Defaults.DefaultFileExistAction
                  if req.FileExistAction != "" {
                    existAction = req.FileExistAction
                  }

                  if existAction == "skip" {
                    log.Warnf("%s already exists, skipping\n", out)
                    continue
                  } else if existAction == "rename" {
                    old := out
                    out = filepath.Join(outDir, fmt.Sprintf("%s_-_%s_%d.mkv", slug.Make(show.Name), slug.Make(show.SubTitle), time.Now().Unix()))
                    log.Warnf("%s already exists, using %s instead\n", old, out)
                  } else if existAction == "replace" {
                    log.Warnf("%s already exists, but will be replaced\n", out)
                  } else {
                    log.Warnf("%s already exists, but an unknown action was specified. Skipping\n", out)
                    continue
                  }
                }

                if now.After(showStart) && now.Before(showEnd) {
                  // Currently Running; but not yet started to record
                  showDuration = showEnd.Sub(now)
                  go func(url, outFile string, duration time.Duration) {
                    RecordVideo(url, outFile, duration)
                  }(station.StreamURL, out, showDuration)
                  log.Infof("Started recording of \"%s (%s)\"\n", show.Name, show.SubTitle)
                } else if now.Before(showStart) {
                  // There is still some time left, sleep, then record
                  sleepTime := showStart.Sub(now)
                  go func(url, outFile, showName, subTitle string, duration time.Duration) {
                    time.Sleep(sleepTime)
                    log.Infof("Starting recording of \"%s (%s)\"\n", showName, subTitle)
                    RecordVideo(url, outFile, duration)
                  }(station.StreamURL, out, show.Name, show.SubTitle, showDuration)
                  log.Infof("Scheduled recording of \"%s (%s)\". Starts in %s\n", show.Name, show.SubTitle, sleepTime)
                }
              }
            }
          }
      }
    }
    nextDaySleepTime := time.Date(now.Year(),  now.Month(), now.Day(), 0, 1, 0, 0, now.Location()).Add(time.Hour * 24).Sub(now)
    log.Infof("Sleeping for %s, see ya!", nextDaySleepTime)
    time.Sleep(nextDaySleepTime)
  }
}
