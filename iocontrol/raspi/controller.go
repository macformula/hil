package raspi

import (
	"context"
	"os/exec"
	"strings"
	"sync"
	"time"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/gpio/gpioreg"
)

const (
	_loggerName   = "raspi_controller"
	_startCommand = "start"
	_stopCommand  = "stop"
)

// Controller provides control for various Raspberry Pi pins
type Controller struct {
	addr string
	l *zap.Logger
	opened bool

	digital [_digitalArraySize]bool
	analog [_analogInputCount + _analogOutputCount]float64
	mu sync.Mutex

	autoload bool
	scriptPath string
	piSSH string
	piPassword string
}

type RaspiOption func(*Controller)

// automatically loads given script onto raspi
func withScriptAutoload(scriptPath, piSSH, piPassword string) RaspiOption {
	return func(c *Controller) {
		c.autoload = true
		c.scriptPath = scriptPath
		c.piSSH = piSSH
		c.piPassword = piPassword
	}
}

// NewController returns a new Raspberry Pi controller
func NewController(l *zap.Logger, address string, opts ...RaspiOption) *Controller {
	pi := &Controller {
		addr: address,
		l: l.Named(_loggerName),
	}

	for _, o := range opts {
		o(pi)
	}

	return pi
}

func (c*Controller) runPiScript(script string, args ...string) error {
	cmdArgs := append([]string{script}, args...)
	cmd := exec.Command("/bin/sh", cmdArgs...)

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "cmd run")
	}

	return nil
}

// Open configures the controller
func (c *Controller) Open(_ context.Context) error {
	c.l.Info("opening raspi controller")
	if c.opened {return nil}

	// initialize periph io
    if _, err := host.Init(); err != nil {
        return errors.Wrap(err, "periph host init")
    }

	if c.autoload {
		err := c.runPiScript(c.scriptPath, c.piSSH, c.piPassword, _startCommand)
		if err != nil {
			return errors.Wrap(err, "run raspi script")
		}
	}

	c.opened = true
	return nil
}

func (c *Controller) Close() error {
	c.l.Info("closing raspberry pi controller")

	if c.autoload {
		err := c.runPiScript(c.scriptPath, c.piSSH, c.piPassword, _stopCommand)
		if err != nil {
			return errors.Wrap(err, "run raspi script")
		}
	}

	c.opened = false
	return nil
}

// SetDigital sets an output digital pin for a Raspberry Pi digital pin
func (c *Controller) SetDigital(output *DigitalPin, b bool) error {
    pin := gpioreg.ByName(strings.ToUpper(id))
    if pin == nil {
        return errors.Errorf("unknown pin %q", output.id)
    }
    // configure output mode
    if err := pin.Out(gpio.Low); err != nil {
        return errors.Wrap(err, "set output mode")
    }
	// set pin level
    if b {
        return errors.Wrap(pin.Out(gpio.High), "write high")
    }
    return errors.Wrap(pin.Out(gpio.Low), "write low")
}

// ReadDigital returns the level of a Raspberry Pi digital pin
func (c *Controller) ReadDigital(output *DigitalPin) (bool, error) {
    pin := gpioreg.ByName(strings.ToUpper(id))
    if pin == nil {
        return false, errors.Errorf("unknown pin %q", input.id)
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
