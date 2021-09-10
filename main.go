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

var verbose = flag.Bool("verbose", false, "print info level logs to stdout")

func main() {
	flag.Parse()

	defer logger.Init("LeaningHydroPiLogger", *verbose, false, ioutil.Discard).Close()

	ctx, cancel := context.WithCancel(context.Background())
	tdpChan := make(chan *tilt.TiltDataPoint)

	publisher := publish.NewCsvPublisher(".", "tilt-log")
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
