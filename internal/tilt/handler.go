package tilt

import (
	"time"

	"github.com/fetherolfjd/leaning-hydro-pi/internal/model"
	"github.com/hashicorp/go-hclog"
)

type DataDecoder interface {
	DeviceUUID([]byte) (string, error)
	TransmitPower([]byte) (int, error)
	SpecificGravity([]byte) (float32, error)
	Temperature([]byte) (float32, error)
}

type TransformTiltData struct {
	DataCh  chan<- model.TiltReading
	Logger  hclog.Logger
	Decoder DataDecoder
}

func (t TransformTiltData) Handle(pkt model.DataPacket) {
	now := time.Now()
	addr := pkt.Address
	rssi := pkt.RSSI

	manData := pkt.ManufacturerData

	uuid, devUUIDErr := t.Decoder.DeviceUUID(manData)
	if devUUIDErr != nil {
		t.Logger.Error("failed to get device ID: %v", devUUIDErr)
		return
	}

	temp, tempErr := t.Decoder.Temperature(manData)
	if tempErr != nil {
		t.Logger.Error("failed to get temperature: %v", tempErr)
		return
	}

	sg, sgErr := t.Decoder.SpecificGravity(manData)
	if sgErr != nil {
		t.Logger.Error("failed to get specific gravity: %v", sgErr)
		return
	}

	txPow, txErr := t.Decoder.TransmitPower(manData)
	if txErr != nil {
		t.Logger.Error("failed to get transmit power: %v", txErr)
	}

	color, ok := uuidToColor[uuid]
	if !ok {
		t.Logger.Error("failed to get color for UUID: %s", uuid)
	}

	t.DataCh <- model.TiltReading{
		Timestamp:       now,
		Address:         addr,
		RSSI:            rssi,
		UUID:            uuid,
		Temperature:     temp,
		SpecificGravity: sg,
		TXPower:         txPow,
		Color:           color,
	}
}
