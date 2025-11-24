package raspi

// ### NOTE ###
// This driver currently requires permissions (sudo) to access the GPIO pins. Running HIL stack on root is not recommended.
// Need to add udev rule to allow gpio group to access /dev/gpiomem and related device nodes without sudo.
// Once the Ansible script is updated to install the udev rule, permissions will no longer be required.

import (
	"context"
	"fmt"
	"sync"

	"github.com/macformula/hil/iocontrol"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

const (
	_loggerName = "raspi_controller"
)

// Mapping from board pins (1-40) to their corresponding BCM GPIO numbers
// CAN HAT uses 8, 10, 3, 5, 15, 19, 21, 23, 24, 26, 27, 28, 33 (https://www.waveshare.com/wiki/2-CH_CAN_HAT)
// Remaining pins are listed below (see https://pinout.xyz/pinout/pin7_gpio4/)
var boardToBCM = map[uint8]int{
	7: 4, 11: 17, 12: 18, 13: 27, 16: 23, 18: 24, 22: 25,
	29: 5, 31: 6, 32: 12, 35: 19, 36: 16, 37: 26, 38: 20, 40: 21,
}

// Controller provides control for various Raspberry Pi pins
type Controller struct {
	l      *zap.Logger
	mu     sync.Mutex
	opened bool
}

// NewController returns a new Raspberry Pi controller
func NewController(l *zap.Logger) *Controller {
	return &Controller{
		l: l.Named(_loggerName),
	}
}

// Open configures the controller
func (c *Controller) Open(_ context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.opened {
		return errors.New("raspberry pi controller already opened")
	}

	// initialize periph io
	if _, err := host.Init(); err != nil {
		return errors.Wrap(err, "periph host init")
	}

	c.opened = true
	return nil
}

// Close tears down the controller instance
func (c *Controller) Close() error {
	c.l.Info("closing raspberry pi controller")
	c.mu.Lock()
	defer c.mu.Unlock()
	c.opened = false
	return nil
}

// SetDigital sets an output digital pin for a Raspberry Pi digital pin
func (c *Controller) WriteDigital(output iocontrol.DigitalPin, level bool) error {
	switch p := output.(type) {
	case *DigitalPin:
		pin, err := resolvePin(p)
		if err != nil {
			return err
		}

		// configure as output, then write
		if err := pin.Out(gpio.Low); err != nil {
			return errors.Wrap(err, "set output mode")
		}
		if level {
			return errors.Wrap(pin.Out(gpio.High), "write high")
		}
		return errors.Wrap(pin.Out(gpio.Low), "write low")

	default:
		return errors.Errorf("Invalid pin type")
	}
}

// ReadDigital returns the level of a Raspberry Pi digital pin
func (c *Controller) ReadDigital(input iocontrol.DigitalPin) (bool, error) {
	switch p := input.(type) {
	case *DigitalPin:
		pin, err := resolvePin(p)
		if err != nil {
			return false, err
		}

		if err := pin.In(gpio.Float, gpio.NoEdge); err != nil {
			return false, errors.Wrap(err, "set input mode")
		}
		return pin.Read() == gpio.High, nil
	default:
		return false, errors.Errorf("Invalid pin type")
	}

}

// WriteVoltage sets the voltage of a Raspberry Pi analog pin
func (c *Controller) WriteVoltage(output iocontrol.AnalogPin, voltage float64) error {
	return errors.New("currently unsupported on raspi")
}

// ReadVoltage returns the voltage of a Raspberry Pi analog pin
func (c *Controller) ReadVoltage(output iocontrol.AnalogPin) (float64, error) {
	return 0.00, errors.New("currently unsupported on raspi")
}

// helpers
func resolvePin(pin *DigitalPin) (gpio.PinIO, error) {
	bcm, ok := boardToBCM[pin.id]
	if !ok {
		return nil, errors.Errorf("no BCM mapping for board pin %d", pin.id)
	}

	if p := gpioreg.ByName(fmt.Sprintf("GPIO%d", bcm)); p != nil {
		return p, nil
	}

	return nil, errors.Errorf("GPIO for board pin %d unavailable", pin.id)
}
