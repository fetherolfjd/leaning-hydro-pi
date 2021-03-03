package tiltconn

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
)

const redUUID string = "A495BB10C5B14B44B5121370F02D74DE"
const greenUUID string = "A495BB20C5B14B44B5121370F02D74DE"
const blackUUID string = "A495BB30C5B14B44B5121370F02D74DE"
const purpleUUID string = "A495BB40C5B14B44B5121370F02D74DE"
const orangeUUID string = "A495BB50C5B14B44B5121370F02D74DE"
const blueUUID string = "A495BB60C5B14B44B5121370F02D74DE"
const yellowUUID string = "A495BB70C5B14B44B5121370F02D74DE"
const pinkUUID string = "A495BB80C5B14B44B5121370F02D74DE"

var uuidMap = make(map[string]bool)

type TiltReading struct {
	UUID            string
	SpecificGravity float32
	Temperature     float32
	RSSI            int
	TXPower         int
	Address         string
}

var readingChan chan *TiltReading

func Configure() chan *TiltReading {
	readingChan = make(chan *TiltReading)
	go doThings()
	return readingChan
}

func doThings() {
	uuidMap[redUUID] = true
	uuidMap[greenUUID] = true
	uuidMap[blackUUID] = true
	uuidMap[purpleUUID] = true
	uuidMap[orangeUUID] = true
	uuidMap[blueUUID] = true
	uuidMap[yellowUUID] = true
	uuidMap[pinkUUID] = true
	d, err := dev.NewDevice("default")
	if err != nil {
		log.Fatalf("Unable to get new device: %v", err)
	}
	ble.SetDefaultDevice(d)
	log.Println("Scanning for BLE advertisements...")
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), 120*time.Second))
	checkError(ble.Scan(ctx, true, advHandler, nil))
}

func advHandler(adv ble.Advertisement) {
	if len(adv.ManufacturerData()) > 20 {
		reading, _ := decodeManufacturerData(adv)
		if reading != nil && readingChan != nil {
			log.Println("Sending reading data!")
			readingChan <- reading
		}
	}
}

func checkError(err error) {
	log.Fatalf("Scan error: %v", err)
}

func decodeManufacturerData(adv ble.Advertisement) (*TiltReading, error) {
	addr := fmt.Sprintf("%s", adv.Addr())

	manData := adv.ManufacturerData()

	uuid := fmt.Sprintf("%X", manData[4:20])
	log.Printf("Checking UUID: %s", uuid)

	if _, ok := uuidMap[uuid]; ok {
		temp := float32(binary.BigEndian.Uint16(manData[20:22])) * 1.0
		sg := float32(binary.BigEndian.Uint16(manData[22:24])) * 0.001

		tx := int(manData[24])

		return &TiltReading{Address: addr, RSSI: adv.RSSI(), UUID: uuid, Temperature: temp, SpecificGravity: sg, TXPower: tx}, nil
	} else {
		log.Printf("UUID %s not in map %v", uuid, uuidMap)
		return nil, errors.New("Not Tilt device")
	}

}
