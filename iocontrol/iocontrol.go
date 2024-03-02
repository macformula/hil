package iocontrol

import (
	"go.uber.org/zap"

	"github.com/macformula/hil/iocontrol/raspi"
	"github.com/macformula/hil/iocontrol/speedgoat"
)

const (
	_loggerName = "iocontrol"
)

// IOControlOption is a type for functions operating on IOControl
type IOControlOption func(*IOControl)

// IOControl contains I/O controllers
type IOControl struct {
	sg *speedgoat.Controller
	rp *raspi.Controller

	l *zap.Logger
}

// NewIOControl returns a new IOControl
func NewIOControl(
	l *zap.Logger,
	opts ...IOControlOption) *IOControl {
	io := &IOControl{
		sg: nil,
		rp: nil,

		l: l.Named(_loggerName),
	}

	for _, o := range opts {
		o(io)
	}

	return io
}

// WithSpeedgoat initializes the iocontroller with a speedgoat device
func WithSpeedgoat(sg *speedgoat.Controller) IOControlOption {
	return func(i *IOControl) {
		i.sg = sg
	}
}

// WithRaspi WithRapsi initializes the iocontroller with a raspi device
func WithRaspi(rp *raspi.Controller) IOControlOption {
	return func(i *IOControl) {
		i.rp = rp
	}
}

// SetDigital sets an output digital pin for a specified pin
func (I *IOControl) SetDigital(output DigitalPin, b bool) error {
	switch t := output.(type) {
	case *speedgoat.DigitalPin:
		if I.sg != nil {
			err := I.sg.SetDigital(t, b)
			if err != nil {
				return err
			}
		}
	case *raspi.DigitalPin:
		if I.rp != nil {
			err := I.rp.SetDigital(t, b)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ReadDigital reads an input digital pin for a specified pin
func (I *IOControl) ReadDigital(output DigitalPin) (bool, error) {

	var lvl bool

	switch t := output.(type) {
	case *speedgoat.DigitalPin:
		if I.sg != nil {
			lvl, err := I.sg.ReadDigital(t)
			if err != nil {
				return lvl, err
			}
		}
	case *raspi.DigitalPin:
		if I.rp != nil {
			lvl, err := I.rp.ReadDigital(t)
			if err != nil {
				return lvl, err
			}
		}
	}
	return lvl, nil
}

// WriteVoltage sets a voltage for a specified output analog pin
func (I *IOControl) WriteVoltage(output AnalogPin, voltage float64) error {
	switch t := output.(type) {
	case *speedgoat.AnalogPin:
		if I.sg != nil {
			err := I.sg.WriteVoltage(t, voltage)
			if err != nil {
				return err
			}
		}
	case *raspi.AnalogPin:
		if I.rp != nil {
			err := I.rp.WriteVoltage(t, voltage)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ReadVoltage returns the voltage of a specified input analog pin
func (I *IOControl) ReadVoltage(output AnalogPin) (float64, error) {

	var voltage float64

	switch t := output.(type) {
	case *speedgoat.AnalogPin:
		if I.sg != nil {
			voltage, err := I.sg.ReadVoltage(t)
			if err != nil {
				return voltage, err
			}
		}
	case *raspi.AnalogPin:
		if I.rp != nil {
			voltage, err := I.rp.ReadVoltage(t)
			if err != nil {
				return voltage, err
			}
		}
	}
	return voltage, nil
}

// WriteCurrent sets the current of a specified output analog pin
func (I *IOControl) WriteCurrent(output AnalogPin, current float64) error {
	switch t := output.(type) {
	case *speedgoat.AnalogPin:
		if I.sg != nil {
			err := I.sg.WriteCurrent(t, current)
			if err != nil {
				return err
			}
		}
	case *raspi.AnalogPin:
		if I.rp != nil {
			err := I.rp.WriteCurrent(t, current)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ReadCurrent returns the current of a specified input analog pin
func (I *IOControl) ReadCurrent(output AnalogPin) (float64, error) {

	var current float64

	switch t := output.(type) {
	case *speedgoat.AnalogPin:
		if I.sg != nil {
			current, err := I.sg.ReadCurrent(t)
			if err != nil {
				return current, err
			}
		}
	case *raspi.AnalogPin:
		if I.rp != nil {
			current, err := I.rp.ReadCurrent(t)
			if err != nil {
				return current, err
			}
		}
	}
	return current, nil
}
