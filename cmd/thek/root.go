package thek

import (
	"fmt"
	"os"

	"github.com/dadav/thek/internal/thek"
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
	  Version: thek.Version,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			thek.SetupLogging(thek.EnableDebugging, thek.EnableJSON)

			log.Debugf("Loading config file: %s", thek.ConfigPath)
			err = thek.LoadConfigFromFile(thek.ConfigPath, &thek.CurrentConfig)
			if err != nil {
				log.Fatal(err)
			}

			log.Debugf("Check if %s is available", thek.RecorderBinary)
			err = thek.CheckRecorder()
			if err != nil {
				log.Fatal(err)
			}

			thek.WatchStreams()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&thek.ConfigPath, "config", "c", "config.yaml", "Path to the config")
	rootCmd.PersistentFlags().BoolVarP(&thek.EnableDebugging, "debug", "d", false, "Enable debug mode")
	rootCmd.PersistentFlags().BoolVarP(&thek.EnableJSON, "json", "j", false, "Enable json output")
}

// Execute runs the program
func Execute() {
	if err := rootCmd.Execute(); err != nil {

		fmt.Fprintf(os.Stderr, "Whoops, there was an error: %s\n", err)
		os.Exit(1)
	}
}
