package speedgoat

// DigitalPin defines a digital pin for the Speedgoat
type DigitalPin struct {
	index uint8
}

// NewDigitalPin returns a new instance of a digital pin
func NewDigitalPin(idx uint8) *DigitalPin {
	pin := DigitalPin{
		index: idx,
	}

	return &pin
}

// String returns the pin type
func (d *DigitalPin) String() string {
	return "speedgoat_digital_pin"
}

// IsDigitalPin ensures the DigitalPin is inherited
func (d *DigitalPin) IsDigitalPin() {}
