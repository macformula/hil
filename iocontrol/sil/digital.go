package sil

// DigitalPin defines a digital pin for SIL
type DigitalPin struct {
	EcuName string
	SigName string
}

func (p *DigitalPin) GetEcuName() string {
	return p.EcuName
}

func (p *DigitalPin) GetSigName() string {
	return p.SigName
}

// NewDigitalInputPin returns a new instance of a digital pin with input direction.
func NewDigitalInputPin(ecu string, signal string) *DigitalPin {
	return &DigitalPin{
		EcuName: ecu,
		SigName: signal,
	}
}

// NewDigitalOutputPin returns a new instance of a digital pin with output direction.
func NewDigitalOutputPin(ecu string, signal string) *DigitalPin {
	return &DigitalPin{
		EcuName: ecu,
		SigName: signal,
	}
}

// String returns the pin type
func (d *DigitalPin) String() string {
	return "sil_digital_pin"
}

// IsDigitalPin ensures the DigitalPin is inherited
func (d *DigitalPin) IsDigitalPin() {}
