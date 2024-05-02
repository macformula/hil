package raspi

import "context"

// Controller provides control for various Raspberry Pi pins
type Controller struct {
}

// NewController returns a new Raspberry Pi controller
func NewController() *Controller {
	return &Controller{}
}

// Open configures the controller
func (c *Controller) Open(_ context.Context) error {
	return nil
}

// SetDigital sets an output digital pin for a Raspberry Pi digital pin
func (c *Controller) SetDigital(output *DigitalPin, b bool) error {
	return nil
}

// ReadDigital returns the level of a Raspberry Pi digital pin
func (c *Controller) ReadDigital(output *DigitalPin) (bool, error) {
	return false, nil
}

// WriteVoltage sets the voltage of a Raspberry Pi analog pin
func (c *Controller) WriteVoltage(output *AnalogPin, voltage float64) error {
	return nil
}

// ReadVoltage returns the voltage of a Raspberry Pi analog pin
func (c *Controller) ReadVoltage(output *AnalogPin) (float64, error) {
	return 0.00, nil
}

// WriteCurrent sets the current of a Raspberry Pi analog pin
func (c *Controller) WriteCurrent(output *AnalogPin, current float64) error {
	return nil
}

// ReadCurrent returns the current of a Raspberry Pi analog pin
func (c *Controller) ReadCurrent(output *AnalogPin) (float64, error) {
	return 0.00, nil
}

func (c *Controller) Close() error {
	return nil
}
