package canlink

import (
	"time"

	"go.einride.tech/can"
)

type TimestampedFrame struct {
	Frame can.Frame
	Time  time.Time
}
