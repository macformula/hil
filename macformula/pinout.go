package macformula

import (
	"github.com/macformula/hil/iocontrol"
	"github.com/macformula/hil/iocontrol/raspi"
	"github.com/macformula/hil/iocontrol/speedgoat"
)

type (
	// DigitalPinout maps physical io to digital pins.
	DigitalPinout map[PhysicalIo]iocontrol.DigitalPin
	// AnalogPinout maps physical io to analog pins.
	AnalogPinout map[PhysicalIo]iocontrol.AnalogPin
)

var _revisionDigitalInputPinout = map[Revision]DigitalPinout{
	Ev5: {
		HvilOk: speedgoat.NewDigitalPin(15),
	},
	MockTest: {
		HvilOk: raspi.NewDigitalPin(),
	},
}

var _revisionDigitalOutputPinout = map[Revision]DigitalPinout{
	Ev5: {
		LvEnableButton:     speedgoat.NewDigitalPin(0),
		ReadyToDriveButton: speedgoat.NewDigitalPin(1),
	},
	MockTest: {
		LvEnableButton:     raspi.NewDigitalPin(),
		ReadyToDriveButton: raspi.NewDigitalPin(),
	},
}

var _revisionAnalogInputPinout = map[Revision]AnalogPinout{
	Ev5: {
		LvController3v3RefVoltage: speedgoat.NewAnalogPin(5),
	},
	MockTest: {
		LvController3v3RefVoltage: raspi.NewAnalogPin(),
	},
}

var _revisionAnalogOutputPinout = map[Revision]AnalogPinout{
	Ev5: {
		AcceleratorPedalPosition1: speedgoat.NewAnalogPin(0),
		AcceleratorPedalPosition2: speedgoat.NewAnalogPin(1),
		AccumulatorCurrent:        speedgoat.NewAnalogPin(2),
	},
	MockTest: {
		AcceleratorPedalPosition1: raspi.NewAnalogPin(),
		AcceleratorPedalPosition2: raspi.NewAnalogPin(),
		AccumulatorCurrent:        raspi.NewAnalogPin(),
	},
}
