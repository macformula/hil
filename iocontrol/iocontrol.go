package iocontrol

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/macformula/hil/iocontrol/raspi"
	"github.com/macformula/hil/iocontrol/sil"
	"github.com/macformula/hil/iocontrol/speedgoat"
)

const (
	_loggerName = "iocontrol"
)

// IOControlOption is a type for functions operating on IOControl
type IOControlOption func(*IOControl)

// IOControl contains I/O controllers
type IOControl struct {
	sg  *speedgoat.Controller
	rp  *raspi.Controller
	sil *sil.Controller

	l *zap.Logger
}

// NewIOControl returns a new IOControl
func NewIOControl(
	l *zap.Logger,
	opts ...IOControlOption) *IOControl {
	io := &IOControl{
		l:  l.Named(_loggerName),
		sg: nil,
		rp: nil,
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

// WithRaspi initializes the iocontroller with a raspi device
func WithRaspi(rp *raspi.Controller) IOControlOption {
	return func(i *IOControl) {
		i.rp = rp
	}
}

// WithSil initializes the iocontroller with a sil device
func WithSil(sil *sil.Controller) IOControlOption {
	return func(i *IOControl) {
		i.sil = sil
	}
}

func (io *IOControl) Open(ctx context.Context) error {
	if io.rp != nil {
		err := io.rp.Open(ctx)
		if err != nil {
			return errors.Wrap(err, "raspi controller open")
		}
	}

	if io.sg != nil {
		err := io.sg.Open(ctx)
		if err != nil {
			return errors.Wrap(err, "speedgoat controller open")
		}
	}

	if io.sil != nil {
		err := io.sil.Open(ctx)
		if err != nil {
			return errors.Wrap(err, "sil controller open")
		}
	}

	return nil
}

// SetDigital sets an output digital pin for a specified pin
func (io *IOControl) SetDigital(output DigitalPin, b bool) error {
	var err error

	switch pin := output.(type) {
	case *speedgoat.DigitalPin:
		if io.sg == nil {
			return errors.New("speedgoat target is nil")
		}

		err = io.sg.SetDigital(pin, b)
		if err != nil {
			return errors.Wrap(err, "set digital")
		}
	case *raspi.DigitalPin:
		if io.rp == nil {
			return errors.New("raspi target is nil")
		}

		err = io.rp.SetDigital(pin, b)
		if err != nil {
			return errors.Wrap(err, "set digital")
		}
	case *sil.DigitalPin:
		if io.sil == nil {
			return errors.New("sil target is nil")
		}

		err = io.sil.SetDigital(pin, b)
		if err != nil {
			return errors.Wrap(err, "set digital")
		}
	default:
		return errors.Errorf("unknown digital pin type (%s)", pin.String())
	}

	return nil
}

// ReadDigital reads an input digital pin for a specified pin
func (io *IOControl) ReadDigital(input DigitalPin) (bool, error) {
	var (
		lvl bool
		err error
	)

	switch pin := input.(type) {
	case *speedgoat.DigitalPin:
		if io.sg == nil {
			return lvl, errors.New("speedgoat target is nil")
		}

		lvl, err = io.sg.ReadDigital(pin)
		if err != nil {
			return false, errors.Wrap(err, "read digital")
		}
	case *raspi.DigitalPin:
		if io.rp == nil {
			return lvl, errors.New("raspi target is nil")
		}

		lvl, err = io.rp.ReadDigital(pin)
		if err != nil {
			return lvl, errors.Wrap(err, "read digital")
		}
	case *sil.DigitalPin:
		if io.sil == nil {
			return lvl, errors.New("sil target is nil")
		}

		lvl, err = io.sil.ReadDigital(pin)
		if err != nil {
			return false, errors.Wrap(err, "read digital")
		}
	default:
		return false, errors.Errorf("unknown digital pin type (%s)", pin.String())
	}

	return lvl, nil
}

// WriteVoltage sets a voltage for a specified output analog pin
func (io *IOControl) WriteVoltage(output AnalogPin, voltage float64) error {
	var err error

	switch pin := output.(type) {
	case *speedgoat.AnalogPin:
		if io.sg == nil {
			return errors.New("speedgoat target is nil")
		}

		err = io.sg.WriteVoltage(pin, voltage)
		if err != nil {
			return errors.Wrap(err, "write voltage")
		}
	case *raspi.AnalogPin:
		if io.rp == nil {
			return errors.New("raspi target is nil")
		}

		err = io.rp.WriteVoltage(pin, voltage)
		if err != nil {
			return errors.Wrap(err, "write voltage")
		}
	case *sil.AnalogPin:
		if io.sil == nil {
			return errors.New("sil target is nil")
		}

		err = io.sil.WriteVoltage(pin, voltage)
		if err != nil {
			return errors.Wrap(err, "write voltage")
		}
	default:
		return errors.Errorf("unknown analog pin type (%s)", pin.String())
	}
	return nil
}

// ReadVoltage returns the voltage of a specified input analog pin
func (io *IOControl) ReadVoltage(input AnalogPin) (float64, error) {
	var (
		voltage float64
		err     error
	)

	switch pin := input.(type) {
	case *speedgoat.AnalogPin:
		if io.sg == nil {
			return voltage, errors.New("speedgoat target is nil")
		}

		voltage, err = io.sg.ReadVoltage(pin)
		if err != nil {
			return 0.0, errors.Wrap(err, "read voltage")
		}
	case *raspi.AnalogPin:
		if io.rp != nil {
			return voltage, errors.New("raspi target is nil")
		}

		voltage, err = io.rp.ReadVoltage(pin)
		if err != nil {
			return voltage, errors.Wrap(err, "read voltage")
		}
	case *sil.AnalogPin:
		if io.sil == nil {
			return voltage, errors.New("sil target is nil")
		}

		voltage, err = io.sil.ReadVoltage(pin)
		if err != nil {
			return 0.0, errors.Wrap(err, "read voltage")
		}
	default:
		return 0.0, errors.Errorf("unknown analog pin type (%s)", pin.String())
	}
	return voltage, nil
}

// WriteCurrent sets the current of a specified output analog pin
func (io *IOControl) WriteCurrent(output AnalogPin, current float64) error {
	var err error

	switch pin := output.(type) {
	case *speedgoat.AnalogPin:
		if io.sg == nil {
			return errors.New("speedgoat target is nil")
		}

		err = io.sg.WriteCurrent(pin, current)
		if err != nil {
			return errors.Wrap(err, "write current")
		}
	case *raspi.AnalogPin:
		if io.rp == nil {
			return errors.New("raspi target is nil")
		}

		err = io.rp.WriteCurrent(pin, current)
		if err != nil {
			return errors.Wrap(err, "write current")
		}
	case *sil.AnalogPin:
		if io.sil == nil {
			return errors.New("sil target is nil")
		}

		err = io.sil.WriteCurrent(pin, current)
		if err != nil {
			return errors.Wrap(err, "write current")
		}
	default:
		return errors.Errorf("unknown analog pin type (%s)", pin.String())
	}

	return nil
}

// ReadCurrent returns the current of a specified input analog pin
func (io *IOControl) ReadCurrent(input AnalogPin) (float64, error) {
	var (
		current float64
		err     error
	)

	switch pin := input.(type) {
	case *speedgoat.AnalogPin:
		if io.sg == nil {
			return current, errors.New("speedgoat target is nil")
		}

		current, err = io.sg.ReadCurrent(pin)
		if err != nil {
			return current, errors.Wrap(err, "read current")
		}
	case *raspi.AnalogPin:
		if io.rp == nil {
			return current, errors.New("raspi target is nil")
		}

		current, err = io.rp.ReadCurrent(pin)
		if err != nil {
			return current, errors.Wrap(err, "read current")
		}
	case *sil.AnalogPin:
		if io.sil == nil {
			return current, errors.New("sil target is nil")
		}

		current, err = io.sil.ReadCurrent(pin)
		if err != nil {
			return current, errors.Wrap(err, "read current")
		}
	default:
		return 0.0, errors.Errorf("unknown analog pin type (%s)", pin.String())
	}

	return current, nil
}

func (io *IOControl) Close() error {
	if io.rp != nil {
		err := io.rp.Close()
		if err != nil {
			return errors.Wrap(err, "raspi controller close")
		}
	}

	if io.sg != nil {
		err := io.sg.Close()
		if err != nil {
			return errors.Wrap(err, "speedgoat controller close")
		}
	}

	if io.sil != nil {
		err := io.sil.Close()
		if err != nil {
			return errors.Wrap(err, "sil controller close")
		}
	}

	return nil
}
