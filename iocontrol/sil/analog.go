package sil

import proto "github.com/macformula/hil/iocontrol/sil/generated"

// AnalogPin defines an analog pin for SIL
type AnalogPin struct {
	Info *proto.SignalInfo
}

func NewAnalogPin(ecu, signal string, access proto.SignalAccess, signalType proto.SignalType) *AnalogPin {
	return &AnalogPin{
		Info: &proto.SignalInfo{
			EcuName:      ecu,
			SignalName:   signal,
			SignalAccess: access,
			SignalType:   signalType,
		},
	}
}

// String returns the pin type
func (d *AnalogPin) String() string {
	return "sil_analog_pin"
}

// IsAnalogPin ensures the AnalogPin is inherited
func (d *AnalogPin) IsAnalogPin() {}
