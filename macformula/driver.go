package macformula

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/macformula/hil/macformula/pinout"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"
)

const (
	_minimumPercentage         = 0.0
	_maximumPercentage         = 100.0
	_accelPedalMaxVoltage      = 3.3
	_accelPedalMinVoltage      = 0.0
	_driverButtonPressDuration = 200 * time.Millisecond
)

// Driver represents the driver of the car.
type Driver struct {
	l             *zap.Logger
	pinController *pinout.Controller
}

// NewDriver creates a new driver instance.
func NewDriver(pc *pinout.Controller, l *zap.Logger) *Driver {
	return &Driver{
		l:             l,
		pinController: pc,
	}
}

// SetAcceleratorPosition takes a percentage [0.0, 100.0] and sets the pedal positions correctly.
func (d *Driver) SetAcceleratorPosition(pedalPercentage float64) error {
	if !checkInclusiveBounds(pedalPercentage, _minimumPercentage, _maximumPercentage) {
		return errors.Errorf("pedal position not in bounds (value: %v, bounds: [%v <= X <= %v])",
			pedalPercentage, _minimumPercentage, _maximumPercentage)
	}

	pedalVoltage1 := rescaler(pedalPercentage,
		_minimumPercentage, _maximumPercentage, _accelPedalMinVoltage, _accelPedalMaxVoltage)

	// TODO: fix pedalVoltage2 to do proper scaling
	pedalVoltage2 := pedalVoltage1

	err := d.pinController.SetVoltage(pinout.AccelPedalPosition1, pedalVoltage1)
	if err != nil {
		return errors.Wrap(err, "set voltage (accel pedal position 1)")
	}

	err = d.pinController.SetVoltage(pinout.AccelPedalPosition2, pedalVoltage2)
	if err != nil {
		return errors.Wrap(err, "set voltage (accel pedal position 2)")
	}

	return nil
}

// SetAcceleratorPositionsOverride takes two percentages [0.0, 100.0] and sets the pedal positions correctly.
func (d *Driver) SetAcceleratorPositionsOverride(pedalPercentage1, pedalPercentage2 float64) error {
	if !checkInclusiveBounds(pedalPercentage1, _minimumPercentage, _maximumPercentage) {
		return errors.Errorf("pedal position1 not in bounds (value: %v, bounds: [%v <= X <= %v])",
			pedalPercentage1, _minimumPercentage, _maximumPercentage)
	}

	if !checkInclusiveBounds(pedalPercentage2, _minimumPercentage, _maximumPercentage) {
		return errors.Errorf("pedal position2 not in bounds (value: %v, bounds: [%v <= X <= %v])",
			pedalPercentage2, _minimumPercentage, _maximumPercentage)
	}

	pedalVoltage1 := rescaler(pedalPercentage1,
		_minimumPercentage, _maximumPercentage, _accelPedalMinVoltage, _accelPedalMaxVoltage)

	pedalVoltage2 := rescaler(pedalPercentage2,
		_minimumPercentage, _maximumPercentage, _accelPedalMinVoltage, _accelPedalMaxVoltage)

	err := d.pinController.SetVoltage(pinout.AccelPedalPosition1, pedalVoltage1)
	if err != nil {
		return errors.Wrap(err, "set accel pedal position 1")
	}

	err = d.pinController.SetVoltage(pinout.AccelPedalPosition2, pedalVoltage2)
	if err != nil {
		return errors.Wrap(err, "set accel pedal position 2")
	}

	return nil
}

// PressStartButton presses the start button for a short duration.
func (d *Driver) PressStartButton(ctx context.Context) error {
	err := d.pinController.SetDigitalLevel(pinout.StartButtonN, false)
	if err != nil {
		return errors.Wrap(err, "set digital level (start button n)")
	}

	err = utils.Sleep(ctx, _driverButtonPressDuration)
	if err != nil {
		return errors.Wrapf(err, "sleep (%vms)", _driverButtonPressDuration.Milliseconds())
	}

	err = d.pinController.SetDigitalLevel(pinout.StartButtonN, true)
	if err != nil {
		return errors.Wrap(err, "set digital level (start button n)")
	}

	return nil
}

func checkInclusiveBounds(value, min, max float64) bool {
	if value > max || value < min {
		return false
	}

	return true
}

func rescaler(value, min1, max1, min2, max2 float64) float64 {
	diff1 := max1 - min1
	diff2 := max2 - min2
	scaleFactor := diff2 / diff1

	return (value-min1)*scaleFactor + min2
}
