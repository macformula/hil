package raspi

import (
	"strings"

	io "github.com/macformula/hil/iocontrol"
	"github.com/pkg/errors"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

// DigitalPin is a digital input/output pin for the Raspberry Pi
type DigitalPin struct {
	pin Pin
	dir io.Direction
	lvl io.Level

	gpio gpio.PinIO
}

// NewDigitalPin returns a new DigitalPin
func NewDigitalPin(
	pin Pin,
	dir io.Direction,
) *DigitalPin {
	digitalPin := &DigitalPin{
		pin: pin,
		dir: dir,
	}

	return digitalPin
}

// Open initializes the host device drivers and acquires a gpio pin
func (d *DigitalPin) Open() error {
	if !d.pin.IsAPin() {
		return errors.New("provided pin is reserved or not available")
	}

	if _, err := host.Init(); err != nil {
		return err
	}

	d.gpio = gpioreg.ByName(strings.ToUpper(d.pin.String()))
	if d.gpio == nil {
		return errors.New("requested gpio is not present")
	}

	return nil
}

// SetDirection modifies the function of the pin
func (d *DigitalPin) SetDirection(dir io.Direction) error {
	if !dir.IsADirection() {
		return errors.New("invalid direction provided")
	}
	d.dir = dir

	return nil
}

// Read returns the current state of the DigitalPin if it is set to Input mode
func (d *DigitalPin) Read() (io.Level, error) {
	if d.dir != io.Input {
		return io.Unknown, errors.Errorf("direction must be set to input to read")
	}

	d.lvl = levelReadConvert(d.gpio.Read())
	return d.lvl, nil
}

// Write sets the state of the DigitalPin if it is set to Output mode
func (d *DigitalPin) Write(lvl io.Level) error {
	if d.dir != io.Output {
		return errors.Errorf("direction must be set to input to write")
	}

	if err := d.gpio.Out(levelWriteConvert(lvl)); err != nil {
		return err
	}

	d.lvl = lvl

	return nil
}

func levelReadConvert(lvl gpio.Level) io.Level {
	var l io.Level
	if lvl == gpio.High {
		l = io.High
	} else {
		l = io.Low
	}
	return l
}

func levelWriteConvert(lvl io.Level) gpio.Level {
	var l gpio.Level
	if lvl == io.High {
		l = gpio.High
	} else {
		l = gpio.Low
	}
	return l
}
