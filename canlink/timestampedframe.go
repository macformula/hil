package canlink

import (
	"time"

	"go.einride.tech/can"
)

// TimestampedFrame contains a single CAN frame along with the time it was received.
type TimestampedFrame struct {
	Frame can.Frame
	Time  time.Time
}
