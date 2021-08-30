package main

import (
	// 	"context"
	// 	"flag"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"leaning-hydro-pi/internal/tilt"
	"log"
	"os"
	"time"

	"github.com/google/logger"
	// 	"log"
	// 	"net/http"
	// 	"os"
	// 	"os/signal"
	// 	"time"
	// 	"internal/dashboard"
	// 	"internal/measurement"
	// 	"github.com/gorilla/mux"
)

// func main() {
// 	fmt.Println("Craptastic")
// }

// var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

var verbose = flag.Bool("verbose", false, "print info level logs to stdout")

func main() {
	var wait time.Duration
	// device := flag.String("device", "default", "Bluetooth device to use")
	flag.DurationVar(&wait, "shutdown-timeout", time.Second*15, "the duration for which the server will gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	defer logger.Init("LeaningHydroPiLogger", *verbose, false, ioutil.Discard).Close()

	ctx := context.Background()
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

	connErr := tilt.Connect(ctx, tdpChan)
	if connErr != nil {
		logger.Fatalf("Error setting up connection: %v", connErr)
	}

	// router := mux.NewRouter()

	// dashboard.Configure(router)

	// var mChan chan measurement.HydrometerMeasurement

	// var wConn *websocket.Conn

	// go func() {
	// 	log.Println("Configuring measurment from main")
	// 	mChan = measurement.Configure()
	// 	log.Println("Measurement configured from main")
	// 	for m := range mChan {
	// 		fmt.Printf("Got measurement: %v\n", m)
	// 		if wConn != nil {
	// 			writeErr := wConn.WriteJSON(&m)
	// 			if writeErr != nil {
	// 				log.Fatalf("Error writing: %v", writeErr)
	// 			}
	// 		}
	// 	}
	// }()

	// router.HandleFunc("/", serveHome)
	// router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	conn, upgradeErr := upgrader.Upgrade(w, r, nil)
	// 	wConn = conn
	// 	if upgradeErr != nil {
	// 		log.Fatalf("Error upgrading socket: %v", upgradeErr)
	// 	}

	// })

	// // server := &http.Server{Addr: "0.0.0.0:8080", WriteTimeout: time.Second * 15, ReadTimeout: time.Second * 15, IdleTimeout: time.Second * 60, Handler: nil}
	// server := &http.Server{Addr: "0.0.0.0:8080", WriteTimeout: time.Second * 15, ReadTimeout: time.Second * 15, IdleTimeout: time.Second * 60, Handler: router}

	// go func() {
	// 	if err := server.ListenAndServe(); err != http.ErrServerClosed {
	// 		log.Fatalf("HTTP server ListenAndServe error: %v", err)
	// 	}
	// }()

	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, os.Interrupt)

	// <-sigChan
	// ctx, cancel := context.WithTimeout(context.Background(), wait)
	// defer cancel()

	// if err := server.Shutdown(ctx); err != nil {
	// 	log.Printf("HTTP server Shutdown: %v", err)
	// }

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

// func serveHome(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "web/templates/dashboard2.html")
// }
