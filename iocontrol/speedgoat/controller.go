package speedgoat

import (
	"encoding/binary"
	"math"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_digitalInputCount  = 8
	_digitalOutputCount = 8
	_digitalOutputIndex = 8
	_analogInputCount   = 8
	_analogOutputCount  = 4
	_analogOutputIndex  = 8
	_loggerName         = "speedgoat_controller"
	_tickTime           = time.Millisecond * 10
	_readDeadline       = time.Second * 5
)

// Controller provides control for various Speedgoat pins
type Controller struct {
	addr string
	conn net.Conn
	l    *zap.Logger

	opened bool

	digital   [_digitalInputCount + _digitalOutputCount]bool
	analog    [_analogInputCount + _analogOutputCount]float64
	muDigital sync.Mutex
	muAnalog  sync.Mutex
}

// NewController returns a new Speedgoat controller.
func NewController(l *zap.Logger, address string) *Controller {
	sg := Controller{
		addr: address,
		l:    l.Named(_loggerName),
	}
	return &sg
}

// Open configures the controller.
func (c *Controller) Open() error {
	c.l.Info("opening speedgoat controller")

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return errors.Wrap(err, "dial speedgoat")
	}

	c.conn = conn
	c.opened = true

	go c.tickOutputs()
	go c.tickInputs()

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
	return nil
}

// SetDigital sets an output digital pin for a Speedgoat digital pin.
func (c *Controller) SetDigital(output *DigitalPin, b bool) {
	c.muDigital.Lock()
	defer c.muDigital.Unlock()
	c.digital[output.Index] = b
}

// ReadDigital returns the level of a Speedgoat digital pin.
func (c *Controller) ReadDigital(output *DigitalPin) bool {
	c.muDigital.Lock()
	defer c.muDigital.Unlock()
	return c.digital[output.Index]
}

// WriteVoltage sets the voltage of a Speedgoat analog pin.
func (c *Controller) WriteVoltage(output *AnalogPin, voltage float64) {
	c.muAnalog.Lock()
	defer c.muDigital.Unlock()
	c.analog[output.Index] = voltage
}

// ReadVoltage returns the voltage of a Speedgoat analog pin.
func (c *Controller) ReadVoltage(output *AnalogPin) float64 {
	c.muAnalog.Lock()
	defer c.muDigital.Unlock()
	return c.analog[output.Index]
}

// WriteCurrent sets the current of a Speedgoat analog pin (unimplemented for Speedgoat).
func (c *Controller) WriteCurrent(output *AnalogPin, current float64) error {
	return nil
}

// ReadCurrent returns the current of a Speedgoat analog pin (unimplemented for Speedgoat).
func (c *Controller) ReadCurrent(output *AnalogPin) (float64, error) {
	return 0.00, nil
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
				panic(err)
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
	for i, digitalOut := range c.digital[_digitalOutputIndex:] {
		if digitalOut {
			data[i] = byte(1)
		} else {
			data[i] = byte(0)
		}
	}
	c.muDigital.Unlock()

	c.muAnalog.Lock()
	for i, analogOutput := range c.analog[_analogOutputIndex:] {
		// Convert the float64 to uint64 and append it as a byte array
		binary.LittleEndian.PutUint64(data[_digitalOutputIndex+i*8:], math.Float64bits(analogOutput))
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
