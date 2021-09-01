package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"leaning-hydro-pi/internal/tilt"
	"log"
	"os"
	"time"

	"os/signal"

	"github.com/google/logger"
)

var verbose = flag.Bool("verbose", false, "print info level logs to stdout")

func main() {
	flag.Parse()

	defer logger.Init("LeaningHydroPiLogger", *verbose, false, ioutil.Discard).Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	tdpChan := make(chan *tilt.TiltDataPoint)

	go func() {
		f, err := os.Create("tilt.csv")
		if err != nil {
			logger.Fatalf("unable to create CSV file: %v", err)
		}
		defer f.Close()
		w := csv.NewWriter(f)
		defer w.Flush()
		logger.Info("Waiting for data points...")
		for dataPoint := range tdpChan {
			// logger.Errorf("Received tilt data point: %v", dataPoint)
			csvValues := convertToCsv(dataPoint)
			logger.Errorf("Made CSV values: %v", csvValues)
			if writeErr := w.Write(csvValues); writeErr != nil {
				log.Fatalf("Error writing to file: %v", writeErr)
			}
			w.Flush()
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

func convertToCsv(tdp *tilt.TiltDataPoint) []string {
	return []string{
		fmt.Sprintf("%d", tdp.Timestamp.UnixMilli()),
		fmt.Sprintf("%f", tdp.SpecificGravity),
		fmt.Sprintf("%f", tdp.Temperature),
	}
}
