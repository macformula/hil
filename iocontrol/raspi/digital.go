package raspi

// DigitalPin defines a digital pin for the Raspberry Pi
type DigitalPin struct {
	id string
}

func NewDigitalPin(idx string) *DigitalPin {
	return &DigitalPin{id: idx}
}

// String returns the pin type
func (d *DigitalPin) String() string {
	return "raspi_digital_pin"
}

// IsDigitalPin ensures the DigitalPin is inherited
func (d *DigitalPin) IsDigitalPin() {}
