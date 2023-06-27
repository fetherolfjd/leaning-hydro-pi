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
	d.Logger.Debug("attempting to access temperature from manufacturer data of length %d", dataLen)
	if dataLen < tempEndByte {
		return 0.0, fmt.Errorf("unable to decode temperature from data of length %d", dataLen)
	}
	tempBytes := data[tempStartByte:tempEndByte]
	d.Logger.Debug("decoding bytes from %X", tempBytes)
	temp := float32(binary.BigEndian.Uint16(tempBytes))
	d.Logger.Debug("decoded temperature of %.2f", temp)
	return temp, nil
}

func (d Decoder) SpecificGravity(data []byte) (float32, error) {
	dataLen := len(data)
	d.Logger.Debug("attempting to access specific gravity from manufacturer data of length %d", dataLen)
	if dataLen < sgEndByte {
		return 0.0, fmt.Errorf("unable to decode specific gravity from data of length %d", dataLen)
	}
	sg := float32(binary.BigEndian.Uint16(data[sgStartByte:sgEndByte])) * 0.001
	d.Logger.Debug("decoded specific gravity of %.2f", sg)
	return sg, nil
}

func (d Decoder) TransmitPower(data []byte) (int, error) {
	dataLen := len(data)
	d.Logger.Debug("attempting to access transmit power from manufacturer data of length %d", dataLen)
	if dataLen < transmitDataByte {
		return 0, fmt.Errorf("unable to decode transmit power from data of length %d", dataLen)
	}
	tx := int(data[transmitDataByte])
	d.Logger.Debug("decoded transmit power of %d", tx)
	return tx, nil
}

func (d Decoder) DeviceUUID(data []byte) (string, error) {
	dataLen := len(data)
	d.Logger.Debug("attempting to access device UUID from manufacturer data of length %d", dataLen)
	if dataLen < deviceUUIDEndByte {
		return "", fmt.Errorf("unable to decode device UUID from data of length %d", dataLen)
	}
	uuid := fmt.Sprintf("%X", data[deviceUUIDStartByte:deviceUUIDEndByte])
	d.Logger.Debug("decoded device UUID of %s", uuid)
	return uuid, nil
}
