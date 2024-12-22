package writer

import (
	"go.einride.tech/can"
	"time"
)

type TimestampedFrame struct {
	Frame can.Frame
	Time  time.Time
}
