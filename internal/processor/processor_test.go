package processor_test

import (
	"context"
	"testing"
	"time"

	"github.com/fetherolfjd/leaning-hydro-pi/internal/model"
	"github.com/fetherolfjd/leaning-hydro-pi/internal/processor"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessor(t *testing.T) {

	uuid1 := "A495BB10C5B14B44B5121370F02D74DE"
	uuid2 := "A495BB10C5B14B44B5121370F02D74DF"
	temp := float32(98.6)
	sg := float32(1.211)
	tx := 2
	color := "RED"
	addr := "banana123"
	rssi := 42
	startTime := time.Now()

	t.Run("publishes averages every time interval", func(t *testing.T) {
		readings := []model.TiltReading{}
		for i := 0; i < 20; i++ {
			readings = append(readings, model.TiltReading{
				UUID:            uuid1,
				SpecificGravity: sg + (0.01 * float32(i)),
				Temperature:     temp + (0.1 * float32(i)),
				RSSI:            rssi,
				TXPower:         tx,
				Address:         addr,
				Color:           color,
				Timestamp:       startTime.Add((time.Duration)(i+1) * time.Second),
			})
			readings = append(readings, model.TiltReading{
				UUID:            uuid2,
				SpecificGravity: sg + (0.02 * float32(i)),
				Temperature:     temp + (0.2 * float32(i)),
				RSSI:            rssi,
				TXPower:         tx,
				Address:         addr,
				Color:           color,
				Timestamp:       startTime.Add((time.Duration)(i+1) * time.Second),
			})
		}

		dataCh := make(chan model.TiltReading, 100)
		defer close(dataCh)
		for _, r := range readings {
			dataCh <- r
		}

		rxData := make(chan model.TiltReading, 1)
		defer close(rxData)
		mp := &mockPblsh{}
		mp.On("Publish", mock.Anything).Run(func(args mock.Arguments) {
			rxData <- args.Get(0).(model.TiltReading)
		}).Return(nil).Twice()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p := processor.HydroReadings{
			Logger:      hclog.NewNullLogger(),
			Publisher:   mp,
			PublishTick: 250 * time.Millisecond,
		}

		p.Process(ctx, dataCh)

		pr1 := <-rxData
		pr2 := <-rxData

		cancel()

		checkData := func(r model.TiltReading, sg, temp float32) {
			assert.Equal(t, sg, r.SpecificGravity)
			assert.Equal(t, temp, r.Temperature)
		}

		if pr1.UUID == uuid1 {
			checkData(pr1, float32(1.306), float32(99.55))
			checkData(pr2, float32(1.4009999), float32(100.5))
		} else {
			checkData(pr1, float32(1.4009999), float32(100.5))
			checkData(pr2, float32(1.306), float32(99.55))
		}

		mp.AssertExpectations(t)
	})

}

type mockPblsh struct {
	mock.Mock
}

func (m *mockPblsh) Publish(r model.TiltReading) error {
	return m.Called(r).Error(0)
}

func (m *mockPblsh) PublishAll(rs []model.TiltReading) error {
	return m.Called(rs).Error(0)
}
