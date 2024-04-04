package raspi

// AnalogPin defines an analog pin for the Raspberry Pi
type AnalogPin struct{}

func NewAnalogPin() *AnalogPin {
	return &AnalogPin{}
}

// String returns the pin type
func (d *AnalogPin) String() string {
	return "raspi_analog_pin"
}

// IsAnalogPin ensures the AnalogPin is inherited
func (d *AnalogPin) IsAnalogPin() {}
