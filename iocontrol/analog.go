package iocontrol

import (
	"context"
	"github.com/ethereum/go-ethereum/event"
)

type DataPoint struct {
	Value float64
	err   error
}

type AnalogPin interface {
	SetDirection(context.Context, Direction) error
	Read(context.Context) (float64, error)
	StartReadStream(sampleRate int) (chan DataPoint, error)
	SubscribeToStream(chan DataPoint) (event.Subscription, error)
	StopReadStream()
	Write(context.Context, float64) error
}
