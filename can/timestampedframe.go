package can

import (
	"go.einride.tech/can"
	"time"
)

const (
	_tracerFormat   = "15:04:05.0000"
	_nameTimeFormat = "15.04.05"
	_nameDateFormat = "2006.01.02"
)

type TimestampedFrame struct {
	Frame can.Frame
	Time  time.Time
}
