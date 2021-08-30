package tilt

import (
	"encoding/hex"
	"testing"
)

func TestDecodeTemperatureWithArrayTooShort(t *testing.T) {
	temp, err := decodeTemperature(nil)
	if temp != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}

	var testSlice []byte
	temp, err = decodeTemperature(testSlice)
	if temp != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}

	testSlice = make([]byte, 0)
	temp, err = decodeTemperature(testSlice)
	if temp != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}

	testSlice = []byte{}
	temp, err = decodeTemperature(testSlice)
	if temp != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}

	s := "4C000215"
	testSlice, decodeErr := hex.DecodeString(s)
	if decodeErr != nil {
		t.Fatal("Failed to decode test hex string")
	}
	temp, err = decodeTemperature(testSlice)
	if temp != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}
}

func TestDecodeTemperature(t *testing.T) {
	data := getTestData()
	temp, err := decodeTemperature(data)
	expectedTemp := 68.0
	if temp != float32(expectedTemp) || err != nil {
		t.Fatalf("Extracted temperature %f not match expected temperature %f", temp, expectedTemp)
	}
}

func TestDecodeSpecificGravityWithArrayTooShort(t *testing.T) {
	sg, err := decodeSpecificGravity(nil)
	if sg != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}

	var testSlice []byte
	sg, err = decodeSpecificGravity(testSlice)
	if sg != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}

	testSlice = make([]byte, 0)
	sg, err = decodeSpecificGravity(testSlice)
	if sg != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}

	testSlice = []byte{}
	sg, err = decodeSpecificGravity(testSlice)
	if sg != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}

	s := "4C000215"
	testSlice, decodeErr := hex.DecodeString(s)
	if decodeErr != nil {
		t.Fatal("Failed to decode test hex string")
	}
	sg, err = decodeSpecificGravity(testSlice)
	if sg != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}
}

func TestDecodeSpecificGravity(t *testing.T) {
	data := getTestData()
	expectedSg := 1.016
	sg, err := decodeSpecificGravity(data)
	if sg != float32(expectedSg) || err != nil {
		t.Fatalf("Extracted specific gravity %f does not match expected specific gravity %f", sg, expectedSg)
	}
}

func TestDecodeTransmitPowerWithArrayTooShort(t *testing.T) {
	tx, err := decodeTransmitPower(nil)
	if tx != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}

	var testSlice []byte
	tx, err = decodeTransmitPower(testSlice)
	if tx != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}

	testSlice = make([]byte, 0)
	tx, err = decodeTransmitPower(testSlice)
	if tx != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}

	testSlice = []byte{}
	tx, err = decodeTransmitPower(testSlice)
	if tx != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}

	s := "4C000215"
	testSlice, decodeErr := hex.DecodeString(s)
	if decodeErr != nil {
		t.Fatal("Failed to decode test hex string")
	}
	tx, err = decodeTransmitPower(testSlice)
	if tx != 0.0 || err == nil {
		t.Fatalf("Empty slice did not fail")
	}
}

func TestDecodeTransmitPower(t *testing.T) {
	data := getTestData()
	expectedTx := 197
	tx, err := decodeTransmitPower(data)
	if tx != expectedTx || err != nil {
		t.Fatalf("Extracted transmit power %d does not match expected transmit power %d", tx, expectedTx)
	}
}

func TestDecodeTiltUUID(t *testing.T) {
	data := getTestData()
	want := "A495BB10C5B14B44B5121370F02D74DE"
	uuid, err := decodeDeviceUUID(data)
	if uuid != want || err != nil {
		t.Fatalf("Extracted device UUID %s does not match expected device UUID %s", uuid, want)
	}
}
