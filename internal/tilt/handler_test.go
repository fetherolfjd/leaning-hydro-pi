package tilt_test

import (
	"errors"
	"testing"

	"github.com/fetherolfjd/leaning-hydro-pi/internal/model"
	"github.com/fetherolfjd/leaning-hydro-pi/internal/tilt"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransformTiltData(t *testing.T) {
	testData := getTestData(t)
	uuid := "A495BB10C5B14B44B5121370F02D74DE"
	temp := float32(98.6)
	sg := float32(1.211)
	tx := 2
	color := "RED"
	addr := "banana123"
	rssi := 42

	getPacket := func() model.DataPacket {
		return model.DataPacket{
			Address:          addr,
			RSSI:             rssi,
			ManufacturerData: testData,
		}
	}

	addUUIDCall := func(m *mockDataDec) *mockDataDec {
		m.On("DeviceUUID", testData).
			Return(uuid, nil)
		return m
	}

	addTempCall := func(m *mockDataDec) *mockDataDec {
		m.On("Temperature", testData).
			Return(temp, nil)
		return m
	}

	addSgCall := func(m *mockDataDec) *mockDataDec {
		m.On("SpecificGravity", testData).
			Return(sg, nil)
		return m
	}

	addTxCall := func(m *mockDataDec) *mockDataDec {
		m.On("TransmitPower", testData).
			Return(tx, nil)
		return m
	}

	testCases := []struct {
		name       string
		dataCh     chan model.TiltReading
		getDecoder func() tilt.DataDecoder
		expect     model.TiltReading
	}{
		{
			name: "no publish when uuid fails",
			getDecoder: func() tilt.DataDecoder {
				mockD := &mockDataDec{}
				mockD.On("DeviceUUID", testData).
					Return("", errors.New("test error"))
				return mockD
			},
		},
		{
			name: "no publish when temp parsing fails",
			getDecoder: func() tilt.DataDecoder {
				mockD := &mockDataDec{}
				addUUIDCall(mockD).
					On("Temperature", testData).
					Return(float32(0.0), errors.New("test error"))
				return mockD
			},
		},
		{
			name: "no publish when gravity parsing fails",
			getDecoder: func() tilt.DataDecoder {
				mockD := &mockDataDec{}
				addUUIDCall(addTempCall(mockD)).
					On("SpecificGravity", testData).
					Return(float32(0.0), errors.New("test error"))
				return mockD
			},
		},
		{
			name: "publish when power parsing fails",
			getDecoder: func() tilt.DataDecoder {
				mockD := &mockDataDec{}
				addUUIDCall(addTempCall(addSgCall(mockD))).
					On("TransmitPower", testData).
					Return(0, errors.New("test error"))
				return mockD
			},
			dataCh: make(chan model.TiltReading, 1),
			expect: model.TiltReading{
				UUID:            uuid,
				SpecificGravity: sg,
				Temperature:     temp,
				RSSI:            rssi,
				TXPower:         0,
				Address:         addr,
				Color:           color,
			},
		},
		{
			name: "publish when no color mapped",
			getDecoder: func() tilt.DataDecoder {
				mockD := &mockDataDec{}
				addTxCall(addTempCall(addSgCall(mockD))).
					On("DeviceUUID", testData).
					Return("peppa pig", nil)
				return mockD
			},
			dataCh: make(chan model.TiltReading, 1),
			expect: model.TiltReading{
				UUID:            "peppa pig",
				SpecificGravity: sg,
				Temperature:     temp,
				RSSI:            rssi,
				TXPower:         tx,
				Address:         addr,
				Color:           "",
			},
		},
		{
			name: "publish when happy",
			getDecoder: func() tilt.DataDecoder {
				mockD := &mockDataDec{}
				addUUIDCall(addTempCall(addSgCall(addTxCall(mockD))))
				return mockD
			},
			dataCh: make(chan model.TiltReading, 1),
			expect: model.TiltReading{
				UUID:            uuid,
				SpecificGravity: sg,
				Temperature:     temp,
				RSSI:            rssi,
				TXPower:         tx,
				Address:         addr,
				Color:           color,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDecoder := tc.getDecoder()
			tfmr := tilt.TransformTiltData{
				DataCh:  tc.dataCh,
				Logger:  hclog.NewNullLogger(),
				Decoder: mockDecoder,
			}
			tfmr.Handle(getPacket())
			if tc.dataCh != nil {
				reading := <-tc.dataCh
				assert.Equal(t, tc.expect.Address, reading.Address)
				assert.Equal(t, tc.expect.RSSI, reading.RSSI)
				assert.Equal(t, tc.expect.UUID, reading.UUID)
				assert.Equal(t, tc.expect.Temperature, reading.Temperature)
				assert.Equal(t, tc.expect.SpecificGravity, reading.SpecificGravity)
				assert.Equal(t, tc.expect.TXPower, reading.TXPower)
				assert.Equal(t, tc.expect.Color, reading.Color)
			}

		})
	}
}

type mockDataDec struct {
	mock.Mock
}

func (m *mockDataDec) DeviceUUID(data []byte) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}

func (m *mockDataDec) TransmitPower(data []byte) (int, error) {
	args := m.Called(data)
	return args.Int(0), args.Error(1)
}

func (m *mockDataDec) SpecificGravity(data []byte) (float32, error) {
	args := m.Called(data)
	return args.Get(0).(float32), args.Error(1)
}

func (m *mockDataDec) Temperature(data []byte) (float32, error) {
	args := m.Called(data)
	return args.Get(0).(float32), args.Error(1)
}
