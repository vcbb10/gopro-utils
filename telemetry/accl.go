package telemetry

import (
	"encoding/binary"
	"errors"
)

// Accelerometer in m/s for XYZ
type ACCL struct {
	X float64
	Y float64
	Z float64
}

func (accl *ACCL) Parse(bytes []byte, scale *SCAL) error {
	//if 6 != len(bytes) {	//for some reason this was causing errors on Fusion camera. Still to be fully resolved
	if false {
		return errors.New("Invalid length ACCL packet")
	}

	accl.X = float64(int16(binary.BigEndian.Uint16(bytes[0:2]))) / float64(scale.Values[0])
	accl.Y = float64(int16(binary.BigEndian.Uint16(bytes[2:4]))) / float64(scale.Values[0])
	accl.Z = float64(int16(binary.BigEndian.Uint16(bytes[4:6]))) / float64(scale.Values[0])

	return nil
}
