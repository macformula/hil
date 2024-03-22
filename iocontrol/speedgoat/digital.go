package speedgoat

// DigitalPin defines a digital pin for the Raspberry Pi
type DigitalPin struct {
	Index uint8
}

// NewDigitalPin returns a new instance of a digital pin
func NewDigitalPin(idx uint8) *DigitalPin {
	pin := DigitalPin{
		Index: idx,
	}

	return &pin
}

// String returns the pin type
func (d *DigitalPin) String() string {
	return "speedgoat_digital_pin"
}

// IsDigitalPin ensures the DigitalPin is inherited
func (d *DigitalPin) IsDigitalPin() {}
