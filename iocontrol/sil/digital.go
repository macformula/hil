package sil

// DigitalPin defines a digital pin for SIL
type DigitalPin struct {
	Ecu_name string
	Sig_name string
}

// NewDigitalInputPin returns a new instance of a digital pin with input direction.
func NewDigitalInputPin(ecu, signal string) *DigitalPin {
	return &DigitalPin{
		Ecu_name: ecu,
		Sig_name: signal,
	}
}

// NewDigitalOutputPin returns a new instance of a digital pin with output direction.
func NewDigitalOutputPin(ecu, signal string) *DigitalPin {
	return &DigitalPin{
		Ecu_name: ecu,
		Sig_name: signal,
	}
}

// String returns the pin type
func (d *DigitalPin) String() string {
	return "sil_digital_pin"
}

// IsDigitalPin ensures the DigitalPin is inherited
func (d *DigitalPin) IsDigitalPin() {}
