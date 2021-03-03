package measurement

import (
	"internal/tiltconn"
	"log"
)

const redUUID string = "A495BB10C5B14B44B5121370F02D74DE"
const greenUUID string = "A495BB20C5B14B44B5121370F02D74DE"
const blackUUID string = "A495BB30C5B14B44B5121370F02D74DE"
const purpleUUID string = "A495BB40C5B14B44B5121370F02D74DE"
const orangeUUID string = "A495BB50C5B14B44B5121370F02D74DE"
const blueUUID string = "A495BB60C5B14B44B5121370F02D74DE"
const yellowUUID string = "A495BB70C5B14B44B5121370F02D74DE"
const pinkUUID string = "A495BB80C5B14B44B5121370F02D74DE"

type HydrometerMeasurement struct {
	Temperature     float32
	SpecificGravity float32
}

func Configure() chan HydrometerMeasurement {
	readingChan := tiltconn.Configure()
	ch := make(chan HydrometerMeasurement, 5)
	log.Println("Configuring measurements")
	go genData(ch, readingChan)
	return ch
}

func genData(c chan HydrometerMeasurement, inChan chan *tiltconn.TiltReading) {
	log.Println("Waiting for data...")
	for true {
		reading := <-inChan
		log.Println("Received reading data; converting and forwarding")
		c <- HydrometerMeasurement{Temperature: float32(reading.Temperature), SpecificGravity: float32(reading.SpecificGravity)}

	}
}
