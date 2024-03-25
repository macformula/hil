package speedgoat

import (
	"encoding/binary"
	"math"
	"net"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_digitalPinCount   = 16
	_analogOutputCount = 4
	_analogInputCount  = 8
	_analogOutputIndex = 8
	_analogPinCount    = 12
	_loggerName        = "speedgoat_controller"
	_tickTime          = time.Millisecond * 10
	_readDeadline      = time.Second
)

// Controller provides control for various Speedgoat pins
type Controller struct {
	addr string
	conn net.Conn
	l    *zap.Logger

	opened bool

	digital [_digitalPinCount]bool
	analog  [_analogPinCount]float64
}

// NewController returns a new Speedgoat controller
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
	c.digital[output.Index] = b
}

// ReadDigital returns the level of a Speedgoat digital pin
func (c *Controller) ReadDigital(output *DigitalPin) bool {
	return c.digital[output.Index]
}

// WriteVoltage sets the voltage of a Speedgoat analog pin
func (c *Controller) WriteVoltage(output *AnalogPin, voltage float64) {
	c.analog[output.Index] = voltage
}

// ReadVoltage returns the voltage of a Speedgoat analog pin
func (c *Controller) ReadVoltage(output *AnalogPin) float64 {
	return c.analog[output.Index]
}

// WriteCurrent sets the current of a Speedgoat analog pin
func (c *Controller) WriteCurrent(output *AnalogPin, current float64) error {
	return nil
}

// ReadCurrent returns the current of a Speedgoat analog pin
func (c *Controller) ReadCurrent(output *AnalogPin) (float64, error) {
	return 0.00, nil
}

// tickOutputs transmits the packed data for the digital and analog outputs to the speedgoat at a set time interval.
func (c *Controller) tickOutputs() {
	// call a pack function for the digital and analog arrays here, transmit every 10 milliseconds

	ticker := time.NewTicker(_tickTime)
	for c.opened {
		for range ticker.C {
			_, err := c.conn.Write(c.packOutputs())
			if err != nil {
				c.l.Error("speedgoat controller", zap.Error(errors.Wrap(err, "connection write")))
			}
		}
	}
	ticker.Stop()
}

// tickInputs reads the pin data from the connection and unpacks it into its respective pin arrays.
func (c *Controller) tickInputs() {
	// call unpack here on digital and analog arrays, receive every 10 milliseconds
	// if we have not received a tcp packet in over a second, error out

	ticker := time.NewTicker(_tickTime)
	for c.opened {
		for range ticker.C {
			data := make([]byte, _digitalPinCount+_analogInputCount*8)

			err := c.conn.SetReadDeadline(time.Now().Add(_readDeadline))
			if err != nil {
				c.l.Error("set read deadline", zap.Error(err))
			}

			_, err = c.conn.Read(data)
			if err != nil {
				c.l.Error("connection read", zap.Error(err))
			}

			c.unpackInputs(data)
		}
	}
	ticker.Stop()
}

// packOutputs packs the data in the output arrays so that it can be sent over TCP.
func (c *Controller) packOutputs() []byte {
	data := make([]byte, _digitalPinCount+_analogOutputCount*8)

	// Digital IO will be ordered in the array first, followed by analog outputs
	for i, digitalPin := range c.digital {
		if digitalPin {
			data[i] = byte(1)
		} else {
			data[i] = byte(0)
		}
	}

	for i, analogOutput := range c.analog[_analogOutputIndex:] {
		// Convert the float64 to uint64 and append it as a byte array
		binary.LittleEndian.PutUint64(data[_digitalPinCount+i:], math.Float64bits(analogOutput))
	}

	return data
}

// unpackInputs takes the received data over TCP and unpacks it into the respective input arrays.
func (c *Controller) unpackInputs(data []byte) {
	for i, digitalPin := range data[:_digitalPinCount] {
		c.digital[i] = digitalPin != 0
	}

	for i := 0; i < _analogInputCount; i++ {
		offset := _digitalPinCount + i*8
		analogInput := data[offset : offset+8]
		c.analog[i] = math.Float64frombits(binary.NativeEndian.Uint64(analogInput))
	}
}
