package main

import (
	"context"
	"flag"
	"io/ioutil"
	"os"

	"github.com/fetherolfjd/leaning-hydro-pi/internal/publish"
	"github.com/fetherolfjd/leaning-hydro-pi/internal/tilt"

	"os/signal"

	"github.com/google/logger"
)

func main() {
	var logFile string
	var logPath string
	var verbose bool

	flag.BoolVar(&verbose, "verbose", false, "Print 'Info' level logs to stdout")
	flag.StringVar(&logFile, "logfile", "", "Name of logfile to write data to")
	flag.StringVar(&logPath, "logpath", ".", "Path of location to write log file to. Default is current directory")
	flag.Parse()

	if logFile == "" {
		logger.Fatal("Must provide value for flag '-logfile'")
	}

	defer logger.Init("LeaningHydroPiLogger", verbose, false, ioutil.Discard).Close()

	ctx, cancel := context.WithCancel(context.Background())
	tdpChan := make(chan *tilt.TiltDataPoint)

	publisher := publish.NewCsvPublisher(logPath, logFile)
	defer publisher.Close()

	go func() {
		logger.Info("Waiting for data points...")
		dps := make([]*tilt.TiltDataPoint, 0)
		for dataPoint := range tdpChan {
			dps = append(dps, dataPoint)
			if len(dps) > 100 {
				publisher.PublishAll(dps)
				dps = make([]*tilt.TiltDataPoint, 0)
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		logger.Info("Starting connection")
		connErr := tilt.Connect(ctx, tdpChan)
		if connErr != nil {
			logger.Errorf("Error setting up connection: %v", connErr)
		}
		logger.Info("Connection ended, sending interrupt")
		sigChan <- os.Interrupt
	}()

	logger.Info("Waiting for signal")
	<-sigChan
	logger.Info("Preparing shutdown")
	cancel()

	logger.Info("Shutting down")
	os.Exit(0)
}
