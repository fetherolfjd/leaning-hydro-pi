package config

import (
	"flag"
	"time"

	"github.com/hashicorp/go-hclog"
)

type LHPConfig struct {
	// LogLevel level to log at
	LogLevel string
	// HydrometerLogFile is the name of file to write readings to
	HydrometerLogFile string
	// HydrometerLogDir is the directory to write the readings file to
	HydrometerLogDir string
	// PublishTick is the time to gather hydrometer readings before publishing the average
	PublishTick time.Duration
}

func Parse() LHPConfig {
	var hydroLogFile string
	var hydroLogDir string
	var logLevel string
	var tickDur time.Duration

	flag.StringVar(&logLevel, "log-level", hclog.Info.String(), "Level to log information to stdout; default is INFO")
	flag.StringVar(&hydroLogFile, "readings-file", "", "Name of file to write hydrometer data to")
	flag.StringVar(&hydroLogDir, "readings-dir", ".", "Path of location to write hydrometer data file. Default is current directory")
	flag.DurationVar(&tickDur, "publish-tick", 1*time.Minute, "The duration over which readings from the hydrometer will be gathered, and then averaged for a single published hydrometer reading; default is 1 minute")
	flag.Parse()

	if hydroLogFile == "" {
		panic("Must provide value for flag '-readings-file'")
	}

	return LHPConfig{
		LogLevel:          logLevel,
		HydrometerLogFile: hydroLogFile,
		HydrometerLogDir:  hydroLogDir,
		PublishTick:       tickDur,
	}
}
