package speedgoat

// DigitalPin defines a digital pin for the Raspberry Pi
type DigitalPin struct {
}

// String returns the pin type
func (d *DigitalPin) String() string {
	return "speedgoat_digital_pin"
}

// IsDigitalPin ensures the DigitalPin is inherited
func (d *DigitalPin) IsDigitalPin() {}
