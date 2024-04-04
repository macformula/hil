package raspi

// DigitalPin defines a digital pin for the Raspberry Pi
type DigitalPin struct {
}

func NewDigitalPin() *DigitalPin {
	return &DigitalPin{}
}

// String returns the pin type
func (d *DigitalPin) String() string {
	return "raspi_digital_pin"
}

// IsDigitalPin ensures the DigitalPin is inherited
func (d *DigitalPin) IsDigitalPin() {}
