package sil

// AnalogPin defines an analog pin for SIL
type AnalogPin struct {
	Ecu_name string
	Sig_name string
}

func NewAnalogInputPin(ecu, signal string) *AnalogPin {
	return &AnalogPin{
		Ecu_name: ecu,
		Sig_name: signal,
	}
}

func NewAnalogOutputPin(ecu, signal string) *AnalogPin {
	return &AnalogPin{
		Ecu_name: ecu,
		Sig_name: signal,
	}
}

// String returns the pin type
func (d *AnalogPin) String() string {
	return "sil_analog_pin"
}

// IsAnalogPin ensures the AnalogPin is inherited
func (d *AnalogPin) IsAnalogPin() {}
