package thek

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

const (
  // RecorderBinary contains the name of the ffmpeg binary
  RecorderBinary = "ffmpeg"
)

// Streams contains the available streams
var Streams map[string]string

// LoadStreamsFromFile loads the streams from the given file
func LoadStreamsFromFile(path string, streams *map[string]string) error {
  jsonFile, err := os.Open(path)
  if err != nil {
    return err
  }
  defer jsonFile.Close()

  byteValue, err := ioutil.ReadAll(jsonFile)
  if err != nil {
    return err
  }
  return json.Unmarshal(byteValue, streams)
}

// CheckRecorder checks if the Recorder binary is available
func CheckRecorder() error {
  _, err := exec.LookPath(RecorderBinary)
  return err
}

// RecordVideo records the given stream for the given duration and saves it to the given output
func RecordVideo(stream, output string, duration time.Duration) error {
  ctx, cancel := context.WithTimeout(context.Background(), duration)
  defer cancel()
  cmd := exec.CommandContext(ctx, RecorderBinary, []string{
    "-i",
    stream,
    "-c",
    "copy",
    output,
  }...)
  return cmd.Run()
}
