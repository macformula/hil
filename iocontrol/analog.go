package iocontrol

import (
	"context"
	"github.com/ethereum/go-ethereum/event"
)

// DataPoint contains relevant data to be transferred over a channel
type DataPoint struct {
	Value float64
	err   error
}

// AnalogPin defines methods for analog pin control
type AnalogPin interface {
	SetDirection(context.Context, Direction) error
	Read(context.Context) (float64, error)
	StartReadStream(sampleRate int) (chan DataPoint, error)
	SubscribeToStream(chan DataPoint) (event.Subscription, error)
	StopReadStream()
	Write(context.Context, float64) error
}
