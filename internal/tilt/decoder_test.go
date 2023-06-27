package tilt_test

import (
	"testing"

	"github.com/fetherolfjd/leaning-hydro-pi/internal/tilt"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestDecoder(t *testing.T) {

	t.Run("Temperature", func(t *testing.T) {
		d := tilt.Decoder{
			Logger: hclog.NewNullLogger(),
		}
		temp, err := d.Temperature(nil)
		assert.Equal(t, float32(0.0), temp)
		assert.Error(t, err)

		badData := getBadTestData(t)
		temp, err = d.Temperature(badData[:6])
		assert.Equal(t, float32(0.0), temp)
		assert.Error(t, err)

		goodData := getTestData(t)
		temp, err = d.Temperature(goodData)
		assert.Equal(t, float32(68.0), temp)
		assert.NoError(t, err)
	})

	t.Run("SpecificGravity", func(t *testing.T) {
		d := tilt.Decoder{
			Logger: hclog.NewNullLogger(),
		}
		sg, err := d.SpecificGravity(nil)
		assert.Equal(t, float32(0.0), sg)
		assert.Error(t, err)

		badData := getBadTestData(t)
		sg, err = d.SpecificGravity(badData[:6])
		assert.Equal(t, float32(0.0), sg)
		assert.Error(t, err)

		goodData := getTestData(t)
		sg, err = d.SpecificGravity(goodData)
		assert.Equal(t, float32(1.016), sg)
		assert.NoError(t, err)
	})

	t.Run("TransmitPower", func(t *testing.T) {
		d := tilt.Decoder{
			Logger: hclog.NewNullLogger(),
		}
		tx, err := d.TransmitPower(nil)
		assert.Equal(t, 0, tx)
		assert.Error(t, err)

		badData := getBadTestData(t)
		tx, err = d.TransmitPower(badData[:6])
		assert.Equal(t, 0, tx)
		assert.Error(t, err)

		goodData := getTestData(t)
		tx, err = d.TransmitPower(goodData)
		assert.Equal(t, 197, tx)
		assert.NoError(t, err)
	})

	t.Run("DeviceUUID", func(t *testing.T) {
		d := tilt.Decoder{
			Logger: hclog.NewNullLogger(),
		}
		uuid, err := d.DeviceUUID(nil)
		assert.Equal(t, "", uuid)
		assert.Error(t, err)

		badData := getBadTestData(t)
		uuid, err = d.DeviceUUID(badData[:6])
		assert.Equal(t, "", uuid)
		assert.Error(t, err)

		goodData := getTestData(t)
		uuid, err = d.DeviceUUID(goodData)
		assert.Equal(t, "A495BB10C5B14B44B5121370F02D74DE", uuid)
		assert.NoError(t, err)
	})

}
