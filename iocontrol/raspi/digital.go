package raspi

// DigitalPin defines a digital pin for the Raspberry Pi
// Works with physical board pin numbering (1-40)
type DigitalPin struct {
	id uint8
}

func NewDigitalPin(idx uint8) *DigitalPin {
	return &DigitalPin{id: idx}
}

// String returns the pin type
func (d *DigitalPin) String() string {
	return "raspi_digital_pin"
}

// IsDigitalPin ensures the DigitalPin is inherited
func (d *DigitalPin) IsDigitalPin() {}
