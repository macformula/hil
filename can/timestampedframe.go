package can

import (
	"go.einride.tech/can"
	"time"
)

const (
	_timeFormat = "15:04:05.0000"
)

type TimestampedFrame struct {
	Frame can.Frame
	Time  time.Time
}
