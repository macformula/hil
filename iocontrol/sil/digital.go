package sil

import pb "github.com/macformula/hil/iocontrol/sil/proto"

// DigitalPin defines a digital pin for SIL
type DigitalPin struct {
	info pb.SignalInfo
}

// NewDigitalInputPin returns a new instance of a digital pin with input direction.
func NewDigitalInputPin(ecu, signal string) *DigitalPin {
	return &DigitalPin{
		info: pb.SignalInfo{
			EcuName:         ecu,
			SignalName:      signal,
			SignalType:      pb.SignalType_SIGNAL_TYPE_DIGITAL,
			SignalDirection: pb.SignalDirection_SIGNAL_DIRECTION_INPUT,
		},
	}
}

// NewDigitalOutputPin returns a new instance of a digital pin with output direction.
func NewDigitalOutputPin(ecu, signal string) *DigitalPin {
	return &DigitalPin{
		info: pb.SignalInfo{
			EcuName:         ecu,
			SignalName:      signal,
			SignalType:      pb.SignalType_SIGNAL_TYPE_DIGITAL,
			SignalDirection: pb.SignalDirection_SIGNAL_DIRECTION_OUTPUT,
		},
	}
}

// String returns the pin type
func (d *DigitalPin) String() string {
	return "sil_digital_pin"
}

// IsDigitalPin ensures the DigitalPin is inherited
func (d *DigitalPin) IsDigitalPin() {}
