package speedgoat

import (
	"context"
	"encoding/binary"
	"math"
	"net"
	"os/exec"
	"sync"
	"time"

	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_digitalInputCount       = 8
	_digitalOutputCount      = 8
	_digitalArraySize        = _digitalInputCount + _digitalOutputCount
	_digitalOutputStartIndex = 8
	_digitalInputStartIndex  = 0
	_analogInputCount        = 8
	_analogOutputCount       = 4
	_analogArraySize         = _analogInputCount + _analogOutputCount
	_analogOutputStartIndex  = 8
	_analogInputStartIndex   = 0
	_loggerName              = "speedgoat_controller"
	_tickTime                = time.Millisecond * 10
	_readDeadline            = time.Second * 5
	_startScript             = "iocontrol/speedgoat/scripts/startModel.sh"
	_stopScript              = "iocontrol/speedgoat/scripts/stopModel.sh"
)

// Controller provides control for various Speedgoat pins
type Controller struct {
	addr string
	conn net.Conn
	err  *utils.ResettableError
	l    *zap.Logger

	opened bool

	digital   [_digitalArraySize]bool
	analog    [_analogInputCount + _analogOutputCount]float64
	muDigital sync.Mutex
	muAnalog  sync.Mutex

	autoloadModel bool
	sgSSH         string
	sgPassword    string
	modelName     string
}

type SpeedgoatOption func(*Controller)

// WithModelAutoload will automatically load the given model into Simulink realtime remotely from the Pi. sgPassword is
// Speedgoat password, sgSSH is the "user@ip" format for SSH, modelName is the name for the simulink model.
func WithModelAutoload(sgPassword, sgSSH, modelName string) SpeedgoatOption {
	return func(c *Controller) {
		c.autoloadModel = true
		c.sgSSH = sgSSH
		c.sgPassword = sgPassword
		c.modelName = modelName
	}
}

// NewController returns a new Speedgoat controller.
func NewController(l *zap.Logger, address string, opts ...SpeedgoatOption) *Controller {
	sg := &Controller{
		addr: address,
		l:    l.Named(_loggerName),
	}

	for _, o := range opts {
		o(sg)
	}

	return sg
}

// Open configures the controller.
func (c *Controller) Open(_ context.Context) error {
	c.l.Info("opening speedgoat controller")

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return errors.Wrap(err, "dial speedgoat")
	}

	c.conn = conn
	c.opened = true

	go c.tickOutputs()
	go c.tickInputs()

	if c.autoloadModel {
		err = c.runSpeedgoatScript(_startScript, c.sgPassword, c.sgSSH, c.modelName)
		if err != nil {
			return errors.Wrap(err, "run speedgoat script")
		}
	}

	return nil
}

// Close closes any resources related to the controller.
func (c *Controller) Close() error {
	c.l.Info("closing speedgoat controller")

	c.opened = false

	err := c.conn.Close()
	if err != nil {
		return errors.Wrap(err, "close speedgoat connection")
	}

	if c.autoloadModel {
		err = c.runSpeedgoatScript(_stopScript, c.sgPassword, c.sgSSH)
		if err != nil {
			return errors.Wrap(err, "run speedgoat script")
		}
	}

	return nil
}

// SetDigital sets an output digital pin for a Speedgoat digital pin.
func (c *Controller) SetDigital(output *DigitalPin, b bool) error {
	c.muDigital.Lock()
	defer c.muDigital.Unlock()

	if output.index >= _digitalArraySize || output.index < _digitalOutputStartIndex {
		return errors.Errorf("invalid output index (%d)", output.index)
	}

	c.digital[output.index] = b

	return nil
}

// ReadDigital returns the level of a Speedgoat digital pin.
func (c *Controller) ReadDigital(input *DigitalPin) (bool, error) {
	c.muDigital.Lock()
	defer c.muDigital.Unlock()

	if input.index >= _digitalArraySize || input.index < _digitalInputStartIndex {
		return false, errors.Errorf("invalid input index (%d)", input.index)
	}

	return c.digital[input.index], nil
}

// WriteVoltage sets the voltage of a Speedgoat analog pin.
func (c *Controller) WriteVoltage(output *AnalogPin, voltage float64) error {
	c.muAnalog.Lock()
	defer c.muAnalog.Unlock()

	if output.index >= _analogArraySize || output.index < _analogOutputStartIndex {
		return errors.Errorf("invalid output index (%d)", output.index)
	}

	c.analog[output.index] = voltage

	return nil
}

// ReadVoltage returns the voltage of a Speedgoat analog pin.
func (c *Controller) ReadVoltage(input *AnalogPin) (float64, error) {
	c.muAnalog.Lock()
	defer c.muAnalog.Unlock()

	if input.index >= _analogArraySize || input.index < _analogInputStartIndex {
		return 0.0, errors.Errorf("invalid input index (%d)", input.index)
	}

	return c.analog[input.index], nil
}

// WriteCurrent sets the current of a Speedgoat analog pin (unimplemented for Speedgoat).
func (c *Controller) WriteCurrent(_ *AnalogPin, _ float64) error {
	return errors.New("unimplemented function on speedgoat controller")
}

// ReadCurrent returns the current of a Speedgoat analog pin (unimplemented for Speedgoat).
func (c *Controller) ReadCurrent(_ *AnalogPin) (float64, error) {
	return 0.0, errors.New("unimplemented function on speedgoat controller")
}

// tickOutputs transmits the packed data for the digital and analog outputs to the Speedgoat at a set time interval.
func (c *Controller) tickOutputs() {
	ticker := time.NewTicker(_tickTime)
	defer ticker.Stop()

	for c.opened {
		for range ticker.C {
			_, err := c.conn.Write(c.packOutputs())
			if err != nil {
				c.l.Error("speedgoat controller", zap.Error(errors.Wrap(err, "connection write")))
			}
		}
	}
}

// tickInputs reads the pin data from the connection and unpacks it into its respective pin arrays.
func (c *Controller) tickInputs() {
	ticker := time.NewTicker(_tickTime)
	defer ticker.Stop()

	for c.opened {
		for range ticker.C {
			data := make([]byte, _digitalInputCount+_analogInputCount*8)

			err := c.conn.SetReadDeadline(time.Now().Add(_readDeadline))
			if err != nil {
				c.l.Error("set read deadline", zap.Error(err))
			}

			_, err = c.conn.Read(data)
			if err != nil {
				c.l.Error("connection read", zap.Error(err))
				panic(err) // temporary until we better handle this
			}

			c.unpackInputs(data)
		}
	}
}

// packOutputs packs the data in the output arrays so that it can be sent over TCP.
func (c *Controller) packOutputs() []byte {
	// Digital IO will be ordered in the array first, followed by analog outputs
	data := make([]byte, _digitalOutputCount+_analogOutputCount*8)

	c.muDigital.Lock()
	for i, digitalOut := range c.digital[_digitalOutputStartIndex:] {
		if digitalOut {
			data[i] = byte(1)
		} else {
			data[i] = byte(0)
		}
	}
	c.muDigital.Unlock()

	c.muAnalog.Lock()
	for i, analogOutput := range c.analog[_analogOutputStartIndex:] {
		// Convert the float64 to uint64 and append it as a byte array
		binary.LittleEndian.PutUint64(data[_digitalOutputStartIndex+i*8:], math.Float64bits(analogOutput))
	}
	c.muAnalog.Unlock()

	return data
}

// unpackInputs takes the received data over TCP and unpacks it into the respective input arrays.
func (c *Controller) unpackInputs(data []byte) {
	c.muDigital.Lock()
	for i := 0; i < _digitalInputCount; i++ {
		c.digital[i] = data[i] != 0
	}
	c.muDigital.Unlock()

	c.muAnalog.Lock()
	for i := 0; i < _analogInputCount; i++ {
		offset := _digitalInputCount + i*8
		analogInput := data[offset : offset+8]
		c.analog[i] = math.Float64frombits(binary.NativeEndian.Uint64(analogInput))
	}
	c.muAnalog.Unlock()
}

// runSpeedgoatScript runs the given script for the Speedgoat. Meant for scripts that SSH into the Speedgoat and run.
func (c *Controller) runSpeedgoatScript(script string, args ...string) error {
	cmdArgs := append([]string{script}, args...)
	cmd := exec.Command("/bin/sh", cmdArgs...)

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "cmd run")
	}
	return nil
}
