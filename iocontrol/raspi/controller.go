package raspi

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

const (
	_loggerName = "raspi_controller"
)

// Pins reference by board number (1-40)
// CAN HAT uses 15, 19, 21, 23, 24, 26, 33
// HAT EEPROM uses 27, 28
// UART 8, 10 and I2C 3, 5 are optional
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

// Closes the controller instance
func (c *Controller) Close() error {
	c.l.Info("closing raspberry pi controller")
	c.mu.Lock()
	defer c.mu.Unlock()
	c.opened = false
	return nil
}

// SetDigital sets an output digital pin for a Raspberry Pi digital pin
func (c *Controller) SetDigital(output *DigitalPin, level bool) error {
	pin, err := resolvePin(output)
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
}

// ReadDigital returns the level of a Raspberry Pi digital pin
func (c *Controller) ReadDigital(input *DigitalPin) (bool, error) {
	pin, err := resolvePin(input)
	if err != nil {
		return false, err
	}

	if err := pin.In(gpio.Float, gpio.NoEdge); err != nil {
		return false, errors.Wrap(err, "set input mode")
	}
	return pin.Read() == gpio.High, nil
}

// WriteVoltage sets the voltage of a Raspberry Pi analog pin
func (c *Controller) WriteVoltage(output *AnalogPin, voltage float64) error {
	return errors.New("currently unsupported on raspi")
}

// ReadVoltage returns the voltage of a Raspberry Pi analog pin
func (c *Controller) ReadVoltage(output *AnalogPin) (float64, error) {
	return 0.00, errors.New("currently unsupported on raspi")
}

// WriteCurrent sets the current of a Raspberry Pi analog pin
func (c *Controller) WriteCurrent(output *AnalogPin, current float64) error {
	return errors.New("currently unsupported on raspi")
}

// ReadCurrent returns the current of a Raspberry Pi analog pin
func (c *Controller) ReadCurrent(output *AnalogPin) (float64, error) {
	return 0.00, errors.New("currently unsupported on raspi")
}

// helpers
func resolvePin(pin *DigitalPin) (gpio.PinIO, error) {
	if pin.id < 1 || pin.id > 40 {
		return nil, errors.Errorf("invalid board pin %d", pin.id)
	}

	if p := gpioreg.ByName(fmt.Sprintf("P1-%d", pin.id)); p != nil {
		return p, nil
	}

	bcm, ok := boardToBCM[pin.id]
	if !ok {
		return nil, errors.Errorf("no BCM mapping for board pin %d", pin.id)
	}

	if p := gpioreg.ByName(fmt.Sprintf("GPIO%d", bcm)); p != nil {
		return p, nil
	}

	return nil, errors.Errorf("GPIO for board pin %d unavailable", pin.id)
}
