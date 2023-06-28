package processor

import (
	"context"
	"time"

	"github.com/fetherolfjd/leaning-hydro-pi/internal/model"
	"github.com/hashicorp/go-hclog"
)

type Publisher interface {
	Publish(model.TiltReading) error
	PublishAll([]model.TiltReading) error
}

type HydroReadings struct {
	Logger      hclog.Logger
	Publisher   Publisher
	PublishTick time.Duration
	readings    map[string][]model.TiltReading
}

func (p HydroReadings) Process(ctx context.Context, dataCh <-chan model.TiltReading) {
	p.readings = make(map[string][]model.TiltReading, 16)

	go func() {
		tkr := time.NewTicker(p.PublishTick)
		defer tkr.Stop()

		for {
			select {
			case <-ctx.Done():
				p.Logger.Info("processor context complete")
				p.doPublish()
				return
			case <-tkr.C:
				p.doPublish()
			case r := <-dataCh:
				data, ok := p.readings[r.UUID]
				if !ok {
					data = make([]model.TiltReading, 0, 1000)
				}
				data = append(data, r)
				p.readings[r.UUID] = data
			}
		}
	}()

}

func (p HydroReadings) doPublish() {
	for uuid, datas := range p.readings {
		len := len(datas)
		if len == 0 {
			p.Logger.Info("no readings to publish", "uuid", uuid)
			continue
		} else if len == 1 {
			p.Logger.Debug("publishing single reading")
			if err := p.Publisher.Publish(datas[0]); err != nil {
				panic(err)
			}
			continue
		}
		avgSg := float32(0.0)
		avgTemp := float32(0.0)
		avgTx := 0
		avgRSSI := 0
		for _, d := range datas {
			avgSg += d.SpecificGravity
			avgTemp += d.Temperature
			avgTx += d.TXPower
			avgRSSI += d.RSSI
		}
		avgSg = avgSg / float32(len)
		avgTemp = avgTemp / float32(len)
		avgTx = avgTx / len
		avgRSSI = avgRSSI / len

		r := datas[len-1]
		r.SpecificGravity = avgSg
		r.Temperature = avgTemp
		r.RSSI = avgRSSI
		r.TXPower = avgTx
		p.Logger.Debug("publishing averaged data", "num_readings", len)
		if err := p.Publisher.Publish(r); err != nil {
			panic(err)
		}
	}
	p.readings = make(map[string][]model.TiltReading, 16)
}
