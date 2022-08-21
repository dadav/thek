package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
  // Version is the current version of this tool
  Version = "0.4.0"
)

func main() {
  var err error

  ParseArguments()

  if PrintVersion {
    fmt.Println(Version)
    os.Exit(0)
  }

  SetupLogging(EnableDebugging, EnableJSON)

  log.Debugf("Loading config file: %s", ConfigPath)
  err = LoadConfigFromFile(ConfigPath, &config)
  if err != nil {
    log.Fatal(err)
  }

  log.Debugf("Check if %s is available", RecorderBinary)
  err = CheckRecorder()
  if err != nil {
    log.Fatal(err)
  }

  WatchStreams()
}
