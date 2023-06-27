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

func TestFilterTiltDevices(t *testing.T) {
	t.Run("returns false when UUID parsing fails", func(t *testing.T) {
		mockDec := &mockDecoder{}
		mockDec.On("DeviceUUID", mock.Anything).
			Return("", errors.New("test error"))

		f := tilt.FilterTiltDevices{
			Logger:  hclog.NewNullLogger(),
			Decoder: mockDec,
		}
		assert.False(t, f.Filter(model.DataPacket{
			ManufacturerData: []byte{},
		}))
	})

	t.Run("returns false when UUID is not expected", func(t *testing.T) {
		mockDec := &mockDecoder{}
		mockDec.On("DeviceUUID", mock.Anything).
			Return("banana32", nil)

		f := tilt.FilterTiltDevices{
			Logger:  hclog.NewNullLogger(),
			Decoder: mockDec,
		}
		assert.False(t, f.Filter(model.DataPacket{
			ManufacturerData: []byte{},
		}))
	})

	t.Run("returns true when UUID is expected", func(t *testing.T) {
		mockDec := &mockDecoder{}
		mockDec.On("DeviceUUID", mock.Anything).
			Return("A495BB10C5B14B44B5121370F02D74DE", nil)

		f := tilt.FilterTiltDevices{
			Logger:  hclog.NewNullLogger(),
			Decoder: mockDec,
		}
		assert.True(t, f.Filter(model.DataPacket{
			ManufacturerData: []byte{},
		}))
	})
}

type mockDecoder struct {
	mock.Mock
}

func (m *mockDecoder) DeviceUUID(data []byte) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}
