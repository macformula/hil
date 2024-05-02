package speedgoat

// AnalogPin defines an analog pin for the Speedgoat
type AnalogPin struct {
	index uint8
}

func NewAnalogPin(idx uint8) *AnalogPin {
	pin := AnalogPin{
		index: idx,
	}

	return &pin
}

// String returns the pin type
func (d *AnalogPin) String() string {
	return "speedgoat_analog_pin"
}

// IsAnalogPin ensures the AnalogPin is inherited
func (d *AnalogPin) IsAnalogPin() {}
