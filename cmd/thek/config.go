package thek

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var (
  // ConfigPath is the path to the config file to use
  ConfigPath string
  // EnableDebugging enables debugging
  EnableDebugging bool
  // EnableJSON enables json output
  EnableJSON bool
  // PrintVersion enables the version printing
  PrintVersion bool
)

var config Config

// Config represents the config structure
type Config struct {
  // Defaults contains all the default values which can be overwritten in other places
  Defaults struct {
    // DefaultOutputDir is the place where the videos will be saved to
    DefaultOutputDir string `yaml:"output_directory"`
    // DefaultSafetyDuration is the time before and after a video
    DefaultSafetyDuration string `yaml:"safety_duration"`
    // DefaultFileExistAction is the action in case a file already exist
    DefaultFileExistAction string `yaml:"file_exist_action"`
  } `yaml:"defaults"`

  // StationURLS is a station to stream mapping
  StationURLS map[string]string `yaml:"stations"`

  // RecorderSchedule contains all the recording tasks
  RecordingTasks []struct {
    // Station is the tv station to follow
    Station string `yaml:"station,omitempty"`
    // ShowKeywords are the words to look out for (all must match)
    ShowKeywords string `yaml:"show_keywords"`
    // SafetyDuration overwrites Defaults.DefaultSafetyDuration
    SafetyDuration string `yaml:"safety_duration,omitempty"`
    // OutputDir overwrites Defaults.DefaultOutputDir
    OutputDir string `yaml:"output_directory,omitempty"`
    // FileExistAction overwrites Defaults.DefaultFileExistAction
    FileExistAction string `yaml:"file_exist_action"`
  } `yaml:"recording_tasks"`
}

// LoadConfigFromFile loads the schedule plan from a file
func LoadConfigFromFile(path string, config *Config) error {
  yamlFile, err := os.Open(path)
  if err != nil {
    return err
  }
  defer yamlFile.Close()
  byteValues, err := ioutil.ReadAll(yamlFile)
  if err != nil {
    return err
  }
  return yaml.Unmarshal(byteValues, config)
}
