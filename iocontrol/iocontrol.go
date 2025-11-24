package iocontrol

import (
	"context"
)

type IOController interface {
	Open(context.Context) error
	Close() error
	WriteDigital(DigitalPin, bool) error
	ReadDigital(DigitalPin) (bool, error)
	WriteVoltage(AnalogPin, float64) error
	ReadVoltage(AnalogPin) (float64, error)
}
