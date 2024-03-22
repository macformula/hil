package speedgoat

// AnalogPin defines an analog pin for the Raspberry Pi
type AnalogPin struct{}

// String returns the pin type
func (d *AnalogPin) String() string {
	return "speedgoat_analog_pin"
}

// IsAnalogPin ensures the AnalogPin is inherited
func (d *AnalogPin) IsAnalogPin() {}
