package tilt

import (
	"encoding/binary"
	"fmt"

	"github.com/google/logger"
)

const tempStartByte int = 20
const tempEndByte int = 22
const sgStartByte int = 22
const sgEndByte int = 24
const transmitDataByte int = 24
const deviceUUIDStartByte int = 4
const deviceUUIDEndByte int = 20

func decodeTemperature(manufacturerData []byte) (float32, error) {
	manDataLen := len(manufacturerData)
	logger.Infof("Attempting to access temperature from manufacturer data of length %d", manDataLen)
	if manDataLen < tempEndByte {
		return 0.0, fmt.Errorf("unable to decode temperature from data of length %d", manDataLen)
	}
	tempBytes := manufacturerData[tempStartByte:tempEndByte]
	logger.Infof("Decoding bytes from %X", tempBytes)
	temp := float32(binary.BigEndian.Uint16(tempBytes))
	logger.Infof("Decoded temperature of %f", temp)
	return temp, nil
}

func decodeSpecificGravity(manufacturerData []byte) (float32, error) {
	manDataLen := len(manufacturerData)
	logger.Infof("Attempting to access specific gravity from manufacturer data of length %d", manDataLen)
	if manDataLen < sgEndByte {
		return 0.0, fmt.Errorf("unable to decode specific gravity from data of length %d", manDataLen)
	}
	sg := float32(binary.BigEndian.Uint16(manufacturerData[sgStartByte:sgEndByte])) * 0.001
	logger.Infof("Decoded specific gravity of %f", sg)
	return sg, nil
}

func decodeTransmitPower(manufacturerData []byte) (int, error) {
	manDataLen := len(manufacturerData)
	logger.Infof("Attempting to access transmit power from manufacturer data of length %d", manDataLen)
	if manDataLen < transmitDataByte {
		return 0, fmt.Errorf("unable to decode transmit power from data of length %d", manDataLen)
	}
	tx := int(manufacturerData[transmitDataByte])
	logger.Infof("Decoded transmit power of %d", tx)
	return tx, nil
}

func decodeDeviceUUID(manufacturerData []byte) (string, error) {
	manDataLen := len(manufacturerData)
	logger.Infof("Attempting to access device UUID from manufacturer data of length %d", manDataLen)
	if manDataLen < deviceUUIDEndByte {
		return "", fmt.Errorf("unable to decode device UUID from data of length %d", manDataLen)
	}
	uuid := fmt.Sprintf("%X", manufacturerData[deviceUUIDStartByte:deviceUUIDEndByte])
	logger.Infof("Decoded device UUID of %s", uuid)
	return uuid, nil
}
