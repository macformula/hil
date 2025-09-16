package sil

// DigitalPin defines a digital pin for SIL
type DigitalPin struct {
	EcuName string
	SigName string
}

// String returns the pin type
func (d *DigitalPin) String() string {
	return "sil_digital_pin"
}

// IsDigitalPin ensures the DigitalPin is inherited
func (d *DigitalPin) IsDigitalPin() {}
