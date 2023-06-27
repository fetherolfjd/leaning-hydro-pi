package main

import (
	"context"
	"path"

	"github.com/hashicorp/go-hclog"

	"github.com/fetherolfjd/leaning-hydro-pi/internal/bt"
	"github.com/fetherolfjd/leaning-hydro-pi/internal/config"
	"github.com/fetherolfjd/leaning-hydro-pi/internal/model"
	"github.com/fetherolfjd/leaning-hydro-pi/internal/processor"
	"github.com/fetherolfjd/leaning-hydro-pi/internal/publish"
	"github.com/fetherolfjd/leaning-hydro-pi/internal/signal"
	"github.com/fetherolfjd/leaning-hydro-pi/internal/tilt"
)

func main() {
	cfg := config.Parse()

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "leaning-hydro-pi",
		Level: hclog.LevelFromString(cfg.LogLevel),
	})

	topCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx := signal.Handle(topCtx)

	csvPub, err := publish.NewCsvPublisher(path.Join(cfg.HydrometerLogDir, cfg.HydrometerLogFile), logger)
	if err != nil {
		panic(err)
	}
	defer func() {
		if cErr := csvPub.Close(); cErr != nil {
			logger.Error("error closing CSV file: %v", cErr)
		}
	}()

	readingChan := make(chan model.TiltReading)

	p := processor.HydroReadings{
		Logger:      logger,
		Publisher:   csvPub,
		PublishTick: cfg.PublishTick,
	}

	p.Process(ctx, readingChan)

	dataDecoder := tilt.Decoder{
		Logger: logger,
	}

	pktFilter := tilt.FilterTiltDevices{
		Logger:  logger,
		Decoder: dataDecoder,
	}

	pktHdlr := tilt.TransformTiltData{
		DataCh:  readingChan,
		Logger:  logger,
		Decoder: dataDecoder,
	}

	err = bt.Scan(ctx, pktHdlr, pktFilter)
	cancel()
	if err != nil {
		logger.Error("error scanning: %v", err)
	}

	logger.Info("Shutting down")
}
