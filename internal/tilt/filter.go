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
		f.Logger.Debug("failed to decode UUID", "error", err)
		return false
	}
	f.Logger.Debug("decode successful", "uuid", uuid)
	_, ok := uuidToColor[uuid]
	return ok
}
