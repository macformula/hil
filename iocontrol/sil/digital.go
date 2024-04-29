package sil

import proto "github.com/macformula/hil/iocontrol/sil/generated"

// DigitalPin defines a digital pin for SIL
type DigitalPin struct {
	Info *proto.SignalInfo
}

// NewDigitalPin returns a new instance of a digital pin
func NewDigitalPin(ecu, signal string, access proto.SignalAccess, signalType proto.SignalType) *DigitalPin {
	return &DigitalPin{
		Info: &proto.SignalInfo{
			EcuName:      ecu,
			SignalName:   signal,
			SignalAccess: access,
			SignalType:   signalType,
		},
	}
}

// String returns the pin type
func (d *DigitalPin) String() string {
	return "sil_digital_pin"
}

// IsDigitalPin ensures the DigitalPin is inherited
func (d *DigitalPin) IsDigitalPin() {}
