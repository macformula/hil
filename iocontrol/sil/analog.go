package sil

import pb "github.com/macformula/hil/iocontrol/sil/proto"

// AnalogPin defines an analog pin for SIL
type AnalogPin struct {
	info pb.SignalInfo
}

func NewAnalogInputPin(ecu, signal string) *AnalogPin {
	return &AnalogPin{
		info: pb.SignalInfo{
			EcuName:         ecu,
			SignalName:      signal,
			SignalType:      pb.SignalType_SIGNAL_TYPE_ANALOG,
			SignalDirection: pb.SignalDirection_SIGNAL_DIRECTION_INPUT,
		},
	}
}

func NewAnalogOutputPin(ecu, signal string) *AnalogPin {
	return &AnalogPin{
		info: pb.SignalInfo{
			EcuName:         ecu,
			SignalName:      signal,
			SignalType:      pb.SignalType_SIGNAL_TYPE_ANALOG,
			SignalDirection: pb.SignalDirection_SIGNAL_DIRECTION_OUTPUT,
		},
	}
}

// String returns the pin type
func (d *AnalogPin) String() string {
	return "sil_analog_pin"
}

// IsAnalogPin ensures the AnalogPin is inherited
func (d *AnalogPin) IsAnalogPin() {}
