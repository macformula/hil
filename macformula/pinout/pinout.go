package pinout

import (
	"github.com/macformula/hil/iocontrol"
	"github.com/macformula/hil/iocontrol/raspi"
	"github.com/macformula/hil/iocontrol/speedgoat"
	"github.com/pkg/errors"
)

type (
	// DigitalPinout maps physical io to digital pins.
	DigitalPinout map[PhysicalIo]iocontrol.DigitalPin
	// AnalogPinout maps physical io to analog pins.
	AnalogPinout map[PhysicalIo]iocontrol.AnalogPin
)

var _revisionDigitalInputPinout = map[Revision]DigitalPinout{
	Ev5: {
		HvilOk: speedgoat.NewDigitalPin(0),
	},
	MockTest: {
		HvilOk: raspi.NewDigitalPin(),
	},
}

var _revisionDigitalOutputPinout = map[Revision]DigitalPinout{
	Ev5: {
		LvEnableButton:     speedgoat.NewDigitalPin(8),
		ReadyToDriveButton: speedgoat.NewDigitalPin(9),
	},
	MockTest: {
		LvEnableButton:     raspi.NewDigitalPin(),
		ReadyToDriveButton: raspi.NewDigitalPin(),
	},
}

var _revisionAnalogInputPinout = map[Revision]AnalogPinout{
	Ev5: {
		LvController3v3RefVoltage: speedgoat.NewAnalogPin(0),
	},
	MockTest: {
		LvController3v3RefVoltage: raspi.NewAnalogPin(),
	},
}

var _revisionAnalogOutputPinout = map[Revision]AnalogPinout{
	Ev5: {
		AcceleratorPedalPosition1: speedgoat.NewAnalogPin(8),
		AcceleratorPedalPosition2: speedgoat.NewAnalogPin(1),
		AccumulatorCurrent:        speedgoat.NewAnalogPin(2),
	},
	MockTest: {
		AcceleratorPedalPosition1: raspi.NewAnalogPin(),
		AcceleratorPedalPosition2: raspi.NewAnalogPin(),
		AccumulatorCurrent:        raspi.NewAnalogPin(),
	},
}

// GetDigitalInputs returns a digital input pinout for the given revision.
func GetDigitalInputs(rev Revision) (DigitalPinout, error) {
	ret, ok := _revisionDigitalInputPinout[rev]
	if !ok {
		return nil, errors.Errorf("no digital input pinout for revision (%s)", rev.String())
	}

	return ret, nil
}

// GetDigitalOutputs returns a digital output pinout for the given revision.
func GetDigitalOutputs(rev Revision) (DigitalPinout, error) {
	ret, ok := _revisionDigitalOutputPinout[rev]
	if !ok {
		return nil, errors.Errorf("no digital output pinout for revision (%s)", rev.String())
	}

	return ret, nil
}

// GetAnalogInputs returns a analog input pinout for the given revision.
func GetAnalogInputs(rev Revision) (AnalogPinout, error) {
	ret, ok := _revisionAnalogInputPinout[rev]
	if !ok {
		return nil, errors.Errorf("no analog input pinout for revision (%s)", rev.String())
	}

	return ret, nil
}

// GetAnalogOutputs returns a analog output pinout for the given revision.
func GetAnalogOutputs(rev Revision) (AnalogPinout, error) {
	ret, ok := _revisionAnalogOutputPinout[rev]
	if !ok {
		return nil, errors.Errorf("no analog output pinout for revision (%s)", rev.String())
	}

	return ret, nil
}
