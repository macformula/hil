package pinout

import (
	"github.com/macformula/hil/iocontrol"
	// "github.com/macformula/hil/iocontrol/raspi"
	"github.com/macformula/hil/iocontrol/sil"
	"github.com/macformula/hil/iocontrol/speedgoat"
)

// DigitalPinout maps physical IO to digital pins.
type DigitalPinout map[PhysicalIo]iocontrol.DigitalPin

// AnalogPinout maps physical IO to analog pins.
type AnalogPinout map[PhysicalIo]iocontrol.AnalogPin

// These inputs and outputs are relative to the SIL/HIL (not the firmware!)
// Ex. A button is a DigitalInput to firmware but is a DigitalOutput of the HIL
type Pinout struct {
	DigitalInputs  DigitalPinout
	DigitalOutputs DigitalPinout
	AnalogInputs   AnalogPinout
	AnalogOutputs  AnalogPinout
}

var _ev5Pinout = Pinout{
	DigitalInputs: DigitalPinout{
		InverterSwitchEn:           speedgoat.NewDigitalPin(0),
		MotorControllerPrechargeEn: speedgoat.NewDigitalPin(1),
		ShutdownCircuitEn:          speedgoat.NewDigitalPin(2),
		AccumulatorEn:              speedgoat.NewDigitalPin(3),
	},
	DigitalOutputs: DigitalPinout{
		GlvmsDisable: speedgoat.NewDigitalPin(8),
		StartButtonN: speedgoat.NewDigitalPin(9),
		HvilDisable:  speedgoat.NewDigitalPin(10),
	},
	AnalogInputs: AnalogPinout{
		HvilFeedback:                 speedgoat.NewAnalogPin(0),
		LvController3v3RefVoltage:    speedgoat.NewAnalogPin(1),
		FrontController3v3RefVoltage: speedgoat.NewAnalogPin(2),
	},
	AnalogOutputs: AnalogPinout{
		SteeringAngle:       speedgoat.NewAnalogPin(8),
		HvCurrentSense:      speedgoat.NewAnalogPin(9),
		AccelPedalPosition1: speedgoat.NewAnalogPin(10),
		AccelPedalPosition2: speedgoat.NewAnalogPin(11),
	},
}

var _mockPinout = Pinout{
	DigitalInputs:  DigitalPinout{},
	DigitalOutputs: DigitalPinout{
		// StartButtonN: raspi.NewDigitalPin(),
	},
	AnalogInputs: AnalogPinout{
		// LvController3v3RefVoltage: raspi.NewAnalogPin(),
	},
	AnalogOutputs: AnalogPinout{
		// AccelPedalPosition1: raspi.NewAnalogPin(),
		// AccelPedalPosition2: raspi.NewAnalogPin(),
		// HvCurrentSense:      raspi.NewAnalogPin(),
	},
}

var _sgTestPinout = Pinout{
	DigitalInputs:  DigitalPinout{},
	DigitalOutputs: DigitalPinout{},
	AnalogInputs:   AnalogPinout{},
	AnalogOutputs:  AnalogPinout{},
}

var _silPinout = Pinout{
	DigitalInputs: DigitalPinout{
		IndicatorLed:               sil.NewDigitalInputPin("DemoProject", IndicatorLed.String()),
		MotorControllerPrechargeEn: sil.NewDigitalInputPin("LvController", MotorControllerPrechargeEn.String()),
		InverterSwitchEn:           sil.NewDigitalInputPin("LvController", InverterSwitchEn.String()),
		AccumulatorEn:              sil.NewDigitalInputPin("LvController", AccumulatorEn.String()),
		ShutdownCircuitEn:          sil.NewDigitalInputPin("LvController", ShutdownCircuitEn.String()),
		TsalEn:                     sil.NewDigitalInputPin("LvController", TsalEn.String()),
		RaspiEn:                    sil.NewDigitalInputPin("LvController", RaspiEn.String()),
		FrontControllerEn:          sil.NewDigitalInputPin("LvController", FrontControllerEn.String()),
		SpeedgoatEn:                sil.NewDigitalInputPin("LvController", SpeedgoatEn.String()),
		MotorControllerEn:          sil.NewDigitalInputPin("LvController", MotorControllerEn.String()),
		ImuGpsEn:                   sil.NewDigitalInputPin("LvController", ImuGpsEn.String()),
		DcdcEn:                     sil.NewDigitalInputPin("LvController", DcdcEn.String()),
		DcdcValid:                  sil.NewDigitalInputPin("LvController", DcdcValid.String()),
		DebugLedEn:                 sil.NewDigitalInputPin("FrontController", DebugLedEn.String()),
		DashboardEn:                sil.NewDigitalInputPin("FrontController", DashboardEn.String()),
		HvilLedEn:                  sil.NewDigitalInputPin("FrontController", HvilLedEn.String()),
		BrakeLightEn:               sil.NewDigitalInputPin("FrontController", BrakeLightEn.String()),
		StatusLedEn:                sil.NewDigitalInputPin("FrontController", StatusLedEn.String()),
		RtdsEn:                     sil.NewDigitalInputPin("FrontController", RtdsEn.String()),
	},
	DigitalOutputs: DigitalPinout{
		IndicatorButton:  sil.NewDigitalOutputPin("DemoProject", IndicatorButton.String()),
		StartButtonN:     sil.NewDigitalOutputPin("FrontController", StartButtonN.String()),
		WheelSpeedLeftA:  sil.NewDigitalOutputPin("FrontController", WheelSpeedLeftA.String()),
		WheelSpeedLeftB:  sil.NewDigitalOutputPin("FrontController", WheelSpeedLeftB.String()),
		WheelSpeedRightA: sil.NewDigitalOutputPin("FrontController", WheelSpeedRightA.String()),
		WheelSpeedRightB: sil.NewDigitalOutputPin("FrontController", WheelSpeedRightB.String()),
		HvilDisable:      sil.NewDigitalOutputPin("FrontController", HvilDisable.String()),
		GlvmsDisable:     sil.NewDigitalOutputPin("LvController", GlvmsDisable.String()),
	},
	AnalogInputs: AnalogPinout{
		HvilFeedback: sil.NewAnalogInputPin("FrontController", HvilFeedback.String()),
	},
	AnalogOutputs: AnalogPinout{
		AccelPedalPosition1: sil.NewAnalogOutputPin("FrontController", AccelPedalPosition1.String()),
		AccelPedalPosition2: sil.NewAnalogOutputPin("FrontController", AccelPedalPosition2.String()),
		SteeringAngle:       sil.NewAnalogOutputPin("FrontController", SteeringAngle.String()),
	},
}

var Pinouts = map[Revision]Pinout{
	Ev5:      _ev5Pinout,
	MockTest: _mockPinout,
	Sil:      _silPinout,
	SgTest:   _sgTestPinout,
}
