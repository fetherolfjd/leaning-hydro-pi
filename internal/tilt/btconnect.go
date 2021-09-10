package tilt

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/google/logger"
)

const redUUID string = "A495BB10C5B14B44B5121370F02D74DE"
const greenUUID string = "A495BB20C5B14B44B5121370F02D74DE"
const blackUUID string = "A495BB30C5B14B44B5121370F02D74DE"
const purpleUUID string = "A495BB40C5B14B44B5121370F02D74DE"
const orangeUUID string = "A495BB50C5B14B44B5121370F02D74DE"
const blueUUID string = "A495BB60C5B14B44B5121370F02D74DE"
const yellowUUID string = "A495BB70C5B14B44B5121370F02D74DE"
const pinkUUID string = "A495BB80C5B14B44B5121370F02D74DE"

var uuidMap = map[string]string{
	redUUID:    "RED",
	greenUUID:  "GREEN",
	blackUUID:  "BLACK",
	purpleUUID: "PURPLE",
	orangeUUID: "ORANGE",
	blueUUID:   "BLUE",
	yellowUUID: "YELLOW",
	pinkUUID:   "PINK",
}

type TiltDataPoint struct {
	UUID            string
	SpecificGravity float32
	Temperature     float32
	RSSI            int
	TXPower         int
	Address         string
	Color           string
	Timestamp       time.Time
}

func (tdp *TiltDataPoint) StrVals() []string {
	return []string{
		fmt.Sprintf("%d", tdp.Timestamp.UnixMilli()),
		tdp.Color,
		tdp.UUID,
		tdp.Address,
		fmt.Sprintf("%1.4f", tdp.SpecificGravity),
		fmt.Sprintf("%.1f", tdp.Temperature),
		fmt.Sprintf("%d", tdp.RSSI),
		fmt.Sprintf("%d", tdp.TXPower),
	}
}

func (tdp *TiltDataPoint) ToJson() ([]byte, error) {
	return json.Marshal(tdp)
}

var newBtDevice = linux.NewDevice
var setBtDevice = ble.SetDefaultDevice
var btScan = ble.Scan

var dataPointChan chan *TiltDataPoint

func Connect(ctx context.Context, tdpChan chan *TiltDataPoint) error {
	logger.Info("Setting up BT connection")
	if tdpChan != nil {
		dataPointChan = tdpChan
	} else {
		return fmt.Errorf("BLE scanning requires you provide a channel to pass data back")
	}
	device, err := newBtDevice()
	if err != nil {
		return fmt.Errorf("unable to get new device: %v", err)
	}
	setBtDevice(device)
	logger.Info("Scanning for BLE advertisements...")

	return btScan(ctx, true, advHandler, advFilter)
}

func advHandler(adv ble.Advertisement) {
	if dataPointChan != nil {
		now := time.Now()
		addr := adv.Addr().String()
		rssi := adv.RSSI()

		manData := adv.ManufacturerData()

		uuid, devUUIDErr := decodeDeviceUUID(manData)
		if devUUIDErr != nil {
			logger.Errorf("%v", devUUIDErr)
		}

		temp, tempErr := decodeTemperature(manData)
		if tempErr != nil {
			logger.Errorf("%v", tempErr)
		}

		sg, sgErr := decodeSpecificGravity(manData)
		if sgErr != nil {
			logger.Errorf("%v", sgErr)
		}

		txPow, txErr := decodeTransmitPower(manData)
		if txErr != nil {
			logger.Errorf("%v", txErr)
		}

		color := uuidMap[uuid]

		tdp := TiltDataPoint{
			Timestamp:       now,
			Address:         addr,
			RSSI:            rssi,
			UUID:            uuid,
			Temperature:     temp,
			SpecificGravity: sg,
			TXPower:         txPow,
			Color:           color,
		}

		dataPointChan <- &tdp
	} else {
		logger.Info("Advertisement data skipped; no data channel specified!")
	}

}

func advFilter(adv ble.Advertisement) bool {
	uuid, err := decodeDeviceUUID(adv.ManufacturerData())
	if err != nil {
		return false
	} else if _, ok := uuidMap[uuid]; ok {
		return true
	} else {
		return false
	}
}
