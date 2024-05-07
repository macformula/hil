package pinout

import (
	"github.com/macformula/hil/fwutils"
	"github.com/macformula/hil/iocontrol"
	"github.com/macformula/hil/iocontrol/raspi"
	"github.com/macformula/hil/iocontrol/sil"
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
		InverterEn:                 speedgoat.NewDigitalPin(0),
		MotorControllerPrechargeEn: speedgoat.NewDigitalPin(1),
		ShutdownCircuitEn:          speedgoat.NewDigitalPin(2),
		AccumulatorEn:              speedgoat.NewDigitalPin(3),
	},
	MockTest: {
		HvilOk: raspi.NewDigitalPin(),
	},
	Sil: {
		IndicatorLed: sil.NewDigitalInputPin("DemoProject", IndicatorLed.String()),
		DebugLedEn:   sil.NewDigitalInputPin(fwutils.FrontController.String(), DebugLedEn.String()),
		DashboardEn:  sil.NewDigitalInputPin(fwutils.FrontController.String(), DashboardEn.String()),
		HvilLedEn:    sil.NewDigitalInputPin(fwutils.FrontController.String(), HvilLedEn.String()),
		BrakeLightEn: sil.NewDigitalInputPin(fwutils.FrontController.String(), BrakeLightEn.String()),
		StatusLedEn:  sil.NewDigitalInputPin(fwutils.FrontController.String(), StatusLedEn.String()),
		RtdsEn:       sil.NewDigitalInputPin(fwutils.FrontController.String(), RtdsEn.String()),
	},
	SgTest: {
		SgTestInput: speedgoat.NewDigitalPin(6),
	},
}

var _revisionDigitalOutputPinout = map[Revision]DigitalPinout{
	Ev5: {
		GlvmsDisable:       speedgoat.NewDigitalPin(8),
		ReadyToDriveButton: speedgoat.NewDigitalPin(9),
		HvilDisable:        speedgoat.NewDigitalPin(10),
	},
	MockTest: {
		LvEnableButton:     raspi.NewDigitalPin(),
		ReadyToDriveButton: raspi.NewDigitalPin(),
	},
	Sil: {
		IndicatorButton:  sil.NewDigitalOutputPin("DemoProject", IndicatorButton.String()),
		StartButtonN:     sil.NewDigitalOutputPin(fwutils.FrontController.String(), StartButtonN.String()),
		WheelSpeedLeftA:  sil.NewDigitalOutputPin(fwutils.FrontController.String(), WheelSpeedLeftA.String()),
		WheelSpeedLeftB:  sil.NewDigitalOutputPin(fwutils.FrontController.String(), WheelSpeedLeftB.String()),
		WheelSpeedRightA: sil.NewDigitalOutputPin(fwutils.FrontController.String(), WheelSpeedRightA.String()),
		WheelSpeedRightB: sil.NewDigitalOutputPin(fwutils.FrontController.String(), WheelSpeedRightB.String()),
		WaitForStart:     sil.NewDigitalOutputPin(fwutils.FrontController.String(), WaitForStart.String()),
		HvilDisable:      sil.NewDigitalOutputPin(fwutils.FrontController.String(), HvilDisable.String()),
	},
	SgTest: {
		SgTestOutput: speedgoat.NewDigitalPin(12),
	},
}

var _revisionAnalogInputPinout = map[Revision]AnalogPinout{
	Ev5: {
		HvilFeedback:                 speedgoat.NewAnalogPin(0),
		LvController3v3RefVoltage:    speedgoat.NewAnalogPin(1),
		FrontController3v3RefVoltage: speedgoat.NewAnalogPin(2),
	},
	MockTest: {
		LvController3v3RefVoltage: raspi.NewAnalogPin(),
	},
	Sil: {
		HvilFeedback: sil.NewAnalogInputPin(fwutils.FrontController.String(), HvilFeedback.String()),
	},
	SgTest: {},
}

var _revisionAnalogOutputPinout = map[Revision]AnalogPinout{
	Ev5: {
		SteeringAngle:       speedgoat.NewAnalogPin(8),
		HvCurrentSense:      speedgoat.NewAnalogPin(9),
		AccelPedalPosition1: speedgoat.NewAnalogPin(10),
		AccelPedalPosition2: speedgoat.NewAnalogPin(11),
	},
	MockTest: {
		AccelPedalPosition1: raspi.NewAnalogPin(),
		AccelPedalPosition2: raspi.NewAnalogPin(),
		AccumulatorCurrent:  raspi.NewAnalogPin(),
	},
	Sil: {
		AccelPedalPosition1: sil.NewAnalogOutputPin(fwutils.FrontController.String(), AccelPedalPosition1.String()),
		AccelPedalPosition2: sil.NewAnalogOutputPin(fwutils.FrontController.String(), AccelPedalPosition2.String()),
		SteeringAngle:       sil.NewAnalogOutputPin(fwutils.FrontController.String(), SteeringAngle.String()),
	},
	SgTest: {},
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
