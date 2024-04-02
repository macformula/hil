package iocontrol

import (
	"github.com/pkg/errors"
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
		l:  l.Named(_loggerName),
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
func (io *IOControl) SetDigital(output DigitalPin, b bool) error {
	switch pin := output.(type) {
	case *speedgoat.DigitalPin:
		if io.sg != nil {
			io.sg.SetDigital(pin, b)
		} else {
			return errors.New("speedgoat target is nil")
		}
	case *raspi.DigitalPin:
		if io.rp != nil {
			err := io.rp.SetDigital(pin, b)
			if err != nil {
				return errors.Wrap(err, "set digital")
			}
		} else {
			return errors.New("raspi target is nil")
		}
	}
	return nil
}

// ReadDigital reads an input digital pin for a specified pin
func (io *IOControl) ReadDigital(input DigitalPin) (bool, error) {
	var lvl bool

	switch pin := input.(type) {
	case *speedgoat.DigitalPin:
		if io.sg != nil {
			lvl = io.sg.ReadDigital(pin)
		} else {
			return lvl, errors.New("speedgoat target is nil")
		}
	case *raspi.DigitalPin:
		if io.rp != nil {
			lvl, err := io.rp.ReadDigital(pin)
			if err != nil {
				return lvl, errors.Wrap(err, "read digital")
			}
		} else {
			return lvl, errors.New("raspi target is nil")
		}
	}
	return lvl, nil
}

// WriteVoltage sets a voltage for a specified output analog pin
func (io *IOControl) WriteVoltage(output AnalogPin, voltage float64) error {
	switch pin := output.(type) {
	case *speedgoat.AnalogPin:
		if io.sg != nil {
			io.sg.WriteVoltage(pin, voltage)
		} else {
			return errors.New("speedgoat target is nil")
		}
	case *raspi.AnalogPin:
		if io.rp != nil {
			err := io.rp.WriteVoltage(pin, voltage)
			if err != nil {
				return errors.Wrap(err, "write voltage")
			}
		} else {
			return errors.New("raspi target is nil")
		}
	}
	return nil
}

// ReadVoltage returns the voltage of a specified input analog pin
func (io *IOControl) ReadVoltage(input AnalogPin) (float64, error) {
	var voltage float64

	switch pin := input.(type) {
	case *speedgoat.AnalogPin:
		if io.sg != nil {
			voltage = io.sg.ReadVoltage(pin)
		} else {
			return voltage, errors.New("speedgoat target is nil")
		}
	case *raspi.AnalogPin:
		if io.rp != nil {
			voltage, err := io.rp.ReadVoltage(pin)
			if err != nil {
				return voltage, errors.Wrap(err, "read voltage")
			}
		} else {
			return voltage, errors.New("raspi target is nil")
		}
	}
	return voltage, nil
}

// WriteCurrent sets the current of a specified output analog pin
func (io *IOControl) WriteCurrent(output AnalogPin, current float64) error {
	switch pin := output.(type) {
	case *speedgoat.AnalogPin:
		if io.sg != nil {
			err := io.sg.WriteCurrent(pin, current)
			if err != nil {
				return errors.Wrap(err, "write current")
			}
		} else {
			return errors.New("speedgoat target is nil")
		}
	case *raspi.AnalogPin:
		if io.rp != nil {
			err := io.rp.WriteCurrent(pin, current)
			if err != nil {
				return errors.Wrap(err, "write current")
			}
		} else {
			errors.New("raspi target is nil")
		}
	}
	return nil
}

// ReadCurrent returns the current of a specified input analog pin
func (io *IOControl) ReadCurrent(input AnalogPin) (float64, error) {
	var current float64

	switch pin := input.(type) {
	case *speedgoat.AnalogPin:
		if io.sg != nil {
			current, err := io.sg.ReadCurrent(pin)
			if err != nil {
				return current, errors.Wrap(err, "read current")
			}
		} else {
			return current, errors.New("speedgoat target is nil")
		}
	case *raspi.AnalogPin:
		if io.rp != nil {
			current, err := io.rp.ReadCurrent(pin)
			if err != nil {
				return current, errors.Wrap(err, "read current")
			}
		} else {
			return current, errors.New("raspi target is nil")
		}
	}
	return current, nil
}
