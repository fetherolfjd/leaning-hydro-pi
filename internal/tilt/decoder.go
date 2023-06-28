package tilt

import (
	"encoding/binary"
	"fmt"

	"github.com/hashicorp/go-hclog"
)

type Decoder struct {
	Logger hclog.Logger
}

func (d Decoder) Temperature(data []byte) (float32, error) {
	dataLen := len(data)
	d.Logger.Trace("attempting to access temperature from manufacturer data", "length", dataLen)
	if dataLen < tempEndByte {
		return 0.0, fmt.Errorf("unable to decode temperature from data of length %d", dataLen)
	}
	tempBytes := data[tempStartByte:tempEndByte]
	temp := float32(binary.BigEndian.Uint16(tempBytes))
	d.Logger.Debug("decode successful", "temperature", temp)
	return temp, nil
}

func (d Decoder) SpecificGravity(data []byte) (float32, error) {
	dataLen := len(data)
	d.Logger.Trace("attempting to access specific gravity from manufacturer data", "length", dataLen)
	if dataLen < sgEndByte {
		return 0.0, fmt.Errorf("unable to decode specific gravity from data of length %d", dataLen)
	}
	sg := float32(binary.BigEndian.Uint16(data[sgStartByte:sgEndByte])) * 0.001
	d.Logger.Debug("decode successful", "specific_gravity", sg)
	return sg, nil
}

func (d Decoder) TransmitPower(data []byte) (int, error) {
	dataLen := len(data)
	d.Logger.Trace("attempting to access transmit power from manufacturer data", "length", dataLen)
	if dataLen < transmitDataByte {
		return 0, fmt.Errorf("unable to decode transmit power from data of length %d", dataLen)
	}
	tx := int(data[transmitDataByte])
	d.Logger.Debug("decode successful", "transmit_power", tx)
	return tx, nil
}

func (d Decoder) DeviceUUID(data []byte) (string, error) {
	dataLen := len(data)
	d.Logger.Trace("attempting to access device UUID from manufacturer data", "length", dataLen)
	if dataLen < deviceUUIDEndByte {
		return "", fmt.Errorf("unable to decode device UUID from data of length %d", dataLen)
	}
	uuid := fmt.Sprintf("%X", data[deviceUUIDStartByte:deviceUUIDEndByte])
	d.Logger.Debug("decode successful", "uuid", uuid)
	return uuid, nil
}
