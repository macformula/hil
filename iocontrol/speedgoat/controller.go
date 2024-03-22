package speedgoat

// Controller provides control for various Speedgoat pins
type Controller struct {
}

// NewController returns a new Speedgoat controller
func NewController() *Controller {
	return nil
}

// Open configures the controller
func (c *Controller) Open() error {
	return nil
}

// SetDigital sets an output digital pin for a Speedgoat digital pin
func (c *Controller) SetDigital(output *DigitalPin, b bool) error {
	return nil
}

// ReadDigital returns the level of a Speedgoat digital pin
func (c *Controller) ReadDigital(output *DigitalPin) (bool, error) {
	return false, nil
}

// WriteVoltage sets the voltage of a Speedgoat analog pin
func (c *Controller) WriteVoltage(output *AnalogPin, voltage float64) error {
	return nil
}

// ReadVoltage returns the voltage of a Speedgoat analog pin
func (c *Controller) ReadVoltage(output *AnalogPin) (float64, error) {
	return 0.00, nil
}

// WriteCurrent sets the current of a Speedgoat analog pin
func (c *Controller) WriteCurrent(output *AnalogPin, current float64) error {
	return nil
}

// ReadCurrent returns the current of a Speedgoat analog pin
func (c *Controller) ReadCurrent(output *AnalogPin) (float64, error) {
	return 0.00, nil
}
