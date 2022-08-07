package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// SetupLogging creates a logging instance
func SetupLogging(debug, json bool) {
  if debug {
    log.SetLevel(log.DebugLevel)
  }
  log.SetOutput(os.Stdout)
  if json {
    log.SetFormatter(&log.JSONFormatter{})
  }else {
    log.SetFormatter(&log.TextFormatter{
      FullTimestamp: true,
    })
  }
}
