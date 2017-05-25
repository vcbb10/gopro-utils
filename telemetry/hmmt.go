package telemetry

import (
	"encoding/binary"
	"errors"
	"time"
)

// Time in Milliseconds
type HMMT struct {
	Time time.Time
}

func (time *HMMT) Parse(bytes []byte) error {
	if 4 != len(bytes) {
		return errors.New("Invalid length HMMT packet")
	}

	bits := time.Parse("060102150405", string(bytes))

	time.Time = t

	return nil
}
