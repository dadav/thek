package thek

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "thek",
		Short: "thek - a simple tool to record things from the mediathek",
		Long: `thek is a tool to schedule recordings for shows of some german streams.

#########################################
The tool requires ffmpeg to work properly
#########################################`,
	  Version: Version,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
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
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&ConfigPath, "config", "c", "config.yaml", "Path to the config")
	rootCmd.PersistentFlags().BoolVarP(&EnableDebugging, "debug", "d", false, "Enable debug mode")
	rootCmd.PersistentFlags().BoolVarP(&EnableJSON, "json", "j", false, "Enable json output")
}

// Execute runs the program
func Execute() {
	if err := rootCmd.Execute(); err != nil {

		fmt.Fprintf(os.Stderr, "Whoops, there was an error: %s\n", err)
		os.Exit(1)
	}
}
