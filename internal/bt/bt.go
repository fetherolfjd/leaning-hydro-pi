package bt

import (
	"context"
	"fmt"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"

	"github.com/fetherolfjd/leaning-hydro-pi/internal/model"
)

type Filterer interface {
	Filter(pkt model.DataPacket) bool
}

type Handler interface {
	Handle(pkt model.DataPacket)
}

var newBtDevice = linux.NewDevice
var setBtDevice = ble.SetDefaultDevice
var scan = ble.Scan

func Scan(ctx context.Context, hdlr Handler, filter Filterer) error {
	dev, err := newBtDevice()
	if err != nil {
		return fmt.Errorf("unable to create new bluetooth device: %w", err)
	}
	setBtDevice(dev)

	hdlrFunc := func(adv ble.Advertisement) {
		hdlr.Handle(transformAdv(adv))
	}

	filterFunc := func(adv ble.Advertisement) bool {
		return filter.Filter(transformAdv(adv))
	}

	return scan(ctx, true, hdlrFunc, filterFunc)
}

func transformAdv(adv ble.Advertisement) model.DataPacket {
	return model.DataPacket{
		Address:          adv.Addr().String(),
		RSSI:             adv.RSSI(),
		ManufacturerData: adv.ManufacturerData(),
	}
}
