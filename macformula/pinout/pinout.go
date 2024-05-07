package pinout

import (
	"github.com/macformula/hil/iocontrol"
	"github.com/macformula/hil/iocontrol/raspi"
	"github.com/macformula/hil/iocontrol/sil"
	"github.com/macformula/hil/iocontrol/speedgoat"
	"github.com/pkg/errors"
)

// DigitalPinout maps physical IO to digital pins.
type DigitalPinout map[PhysicalIo]iocontrol.DigitalPin

// AnalogPinout maps physical IO to analog pins.
type AnalogPinout map[PhysicalIo]iocontrol.AnalogPin

// Revision-specific digital input pinouts
var _revisionDigitalInputPinout = map[Revision]DigitalPinout{
	Ev5: {
		InverterSwitchEn:           speedgoat.NewDigitalPin(0),
		MotorControllerPrechargeEn: speedgoat.NewDigitalPin(1),
		ShutdownCircuitEn:          speedgoat.NewDigitalPin(2),
		AccumulatorEn:              speedgoat.NewDigitalPin(3),
	},
	MockTest: {},
	Sil: {
		IndicatorLed:               sil.NewDigitalInputPin("DemoProject", IndicatorLed.String()),
		MotorControllerPrechargeEn: sil.NewDigitalInputPin("LvController", MotorControllerPrechargeEn.String()),
		InverterSwitchEn:           sil.NewDigitalInputPin("LvController", InverterSwitchEn.String()),
		AccumulatorEn:              sil.NewDigitalInputPin("LvController", AccumulatorEn.String()),
		ShutdownCircuitEn:          sil.NewDigitalInputPin("LvController", ShutdownCircuitEn.String()),
		DebugLedEn:                 sil.NewDigitalInputPin("FrontController", DebugLedEn.String()),
		DashboardEn:                sil.NewDigitalInputPin("FrontController", DashboardEn.String()),
		HvilLedEn:                  sil.NewDigitalInputPin("FrontController", HvilLedEn.String()),
		BrakeLightEn:               sil.NewDigitalInputPin("FrontController", BrakeLightEn.String()),
		StatusLedEn:                sil.NewDigitalInputPin("FrontController", StatusLedEn.String()),
		RtdsEn:                     sil.NewDigitalInputPin("FrontController", RtdsEn.String()),
	},
}

// Revision-specific digital output pinouts
var _revisionDigitalOutputPinout = map[Revision]DigitalPinout{
	Ev5: {
		GlvmsDisable: speedgoat.NewDigitalPin(8),
		StartButtonN: speedgoat.NewDigitalPin(9),
		HvilDisable:  speedgoat.NewDigitalPin(10),
	},
	MockTest: {
		StartButtonN: raspi.NewDigitalPin(),
	},
	Sil: {
		IndicatorButton:  sil.NewDigitalOutputPin("DemoProject", IndicatorButton.String()),
		StartButtonN:     sil.NewDigitalOutputPin("FrontController", StartButtonN.String()),
		WheelSpeedLeftA:  sil.NewDigitalOutputPin("FrontController", WheelSpeedLeftA.String()),
		WheelSpeedLeftB:  sil.NewDigitalOutputPin("FrontController", WheelSpeedLeftB.String()),
		WheelSpeedRightA: sil.NewDigitalOutputPin("FrontController", WheelSpeedRightA.String()),
		WheelSpeedRightB: sil.NewDigitalOutputPin("FrontController", WheelSpeedRightB.String()),
		HvilDisable:      sil.NewDigitalOutputPin("FrontController", HvilDisable.String()),
		GlvmsDisable:     sil.NewDigitalOutputPin("LvController", GlvmsDisable.String()),
	},
}

// Revision-specific analog input pinouts
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
		HvilFeedback: sil.NewAnalogInputPin("FrontController", HvilFeedback.String()),
	},
}

// Revision-specific analog output pinouts
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
		HvCurrentSense:      raspi.NewAnalogPin(),
	},
	Sil: {
		AccelPedalPosition1: sil.NewAnalogOutputPin("FrontController", AccelPedalPosition1.String()),
		AccelPedalPosition2: sil.NewAnalogOutputPin("FrontController", AccelPedalPosition2.String()),
		SteeringAngle:       sil.NewAnalogOutputPin("FrontController", SteeringAngle.String()),
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

// GetAnalogInputs returns an analog input pinout for the given revision.
func GetAnalogInputs(rev Revision) (AnalogPinout, error) {
	ret, ok := _revisionAnalogInputPinout[rev]
	if !ok {
		return nil, errors.Errorf("no analog input pinout for revision (%s)", rev.String())
	}

	return ret, nil
}

// GetAnalogOutputs returns an analog output pinout for the given revision.
func GetAnalogOutputs(rev Revision) (AnalogPinout, error) {
	ret, ok := _revisionAnalogOutputPinout[rev]
	if !ok {
		return nil, errors.Errorf("no analog output pinout for revision (%s)", rev.String())
	}

	return ret, nil
}
