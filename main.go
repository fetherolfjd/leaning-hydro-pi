package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"internal/dashboard"
	"internal/measurement"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func main() {
	var wait time.Duration
	// device := flag.String("device", "default", "Bluetooth device to use")
	flag.DurationVar(&wait, "shutdown-timeout", time.Second*15, "the duration for which the server will gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	router := mux.NewRouter()

	dashboard.Configure(router)

	var mChan chan measurement.HydrometerMeasurement

	var wConn *websocket.Conn

	go func() {
		log.Println("Configuring measurment from main")
		mChan = measurement.Configure()
		log.Println("Measurement configured from main")
		for m := range mChan {
			fmt.Printf("Got measurement: %v\n", m)
			if wConn != nil {
				writeErr := wConn.WriteJSON(&m)
				if writeErr != nil {
					log.Fatalf("Error writing: %v", writeErr)
				}
			}
		}
	}()

	router.HandleFunc("/", serveHome)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, upgradeErr := upgrader.Upgrade(w, r, nil)
		wConn = conn
		if upgradeErr != nil {
			log.Fatalf("Error upgrading socket: %v", upgradeErr)
		}

	})

	// server := &http.Server{Addr: "0.0.0.0:8080", WriteTimeout: time.Second * 15, ReadTimeout: time.Second * 15, IdleTimeout: time.Second * 60, Handler: nil}
	server := &http.Server{Addr: "0.0.0.0:8080", WriteTimeout: time.Second * 15, ReadTimeout: time.Second * 15, IdleTimeout: time.Second * 60, Handler: router}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	log.Println("Shutting down")
	os.Exit(0)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/dashboard2.html")
}
