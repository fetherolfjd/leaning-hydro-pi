package tilt

import (
	"github.com/fetherolfjd/leaning-hydro-pi/internal/model"
	"github.com/hashicorp/go-hclog"
)

type UUIDDecoder interface {
	DeviceUUID([]byte) (string, error)
}

type FilterTiltDevices struct {
	Logger  hclog.Logger
	Decoder UUIDDecoder
}

func (f FilterTiltDevices) Filter(pkt model.DataPacket) bool {
	uuid, err := f.Decoder.DeviceUUID(pkt.ManufacturerData)
	if err != nil {
		f.Logger.Error("failed to decode UUID: %v", err)
		return false
	}
	f.Logger.Debug("decoded device UUID of: %s", uuid)
	_, ok := uuidToColor[uuid]
	return ok
}
