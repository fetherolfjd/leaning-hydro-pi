package tilt_test

import (
	"encoding/hex"
	"testing"
)

func getTestData(t *testing.T) []byte {
	s := "4C000215A495BB10C5B14B44B5121370F02D74DE004403F8C5C7"
	data, decodeErr := hex.DecodeString(s)
	if decodeErr != nil {
		t.Fatal("Failed to decode test hex string")
	}
	return data
}

func getBadTestData(t *testing.T) []byte {
	s := "4C0C0215A495BC10C5B14B44C5121370F02D74DCCC4403F8CCC7"
	data, decodeErr := hex.DecodeString(s)
	if decodeErr != nil {
		t.Fatal("Failed to decode test hex string")
	}
	return data
}
