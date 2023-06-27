package model

import "time"

type DataPacket struct {
	Address          string
	RSSI             int
	ManufacturerData []byte
}

type TiltReading struct {
	UUID            string
	SpecificGravity float32
	Temperature     float32
	RSSI            int
	TXPower         int
	Address         string
	Color           string
	Timestamp       time.Time
}
