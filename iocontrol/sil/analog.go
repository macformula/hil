package sil

// AnalogPin defines an analog pin for SIL
type AnalogPin struct {
	EcuName string
	SigName string
}

func (p *AnalogPin) GetEcuName() string {
	return p.EcuName
}

func (p *AnalogPin) GetSigName() string {
	return p.EcuName
}

func NewAnalogInputPin(ecu, signal string) *AnalogPin {
	return &AnalogPin{
		EcuName: ecu,
		SigName: signal,
	}
}

func NewAnalogOutputPin(ecu, signal string) *AnalogPin {
	return &AnalogPin{
		EcuName: ecu,
		SigName: signal,
	}
}

// String returns the pin type
func (d *AnalogPin) String() string {
	return "sil_analog_pin"
}

// IsAnalogPin ensures the AnalogPin is inherited
func (d *AnalogPin) IsAnalogPin() {}
