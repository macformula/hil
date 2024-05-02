package pinout

// PhysicalIo represents an input output in its physical meaning.
type PhysicalIo int

//go:generate enumer -type=PhysicalIo "physicalio.go"

const (
	UnknownPhysicalIo PhysicalIo = iota
	// LvEnableButton enables the low voltage system.
	LvEnableButton
	// ReadyToDriveButton is the final button to be pressed before the car can be driven.
	ReadyToDriveButton
	// HvilOk indicates that the high voltage interlock loop is satisfied.
	HvilOk
	// AcceleratorPedalPosition1 determines how fast the car should go. It should be offset from AcceleratorPedalPosition2.
	AcceleratorPedalPosition1
	// AcceleratorPedalPosition2 determines how fast the car should go. It should be offset from AcceleratorPedalPosition1.
	AcceleratorPedalPosition2
	// AccumulatorCurrent determines the current that is passing through the battery pack (accumulator).
	AccumulatorCurrent
	// LvController3v3RefVoltage is connected to the 3.3V touch point on the lv controller.
	LvController3v3RefVoltage
	// IndicatorLed is used in the firmware DemoProject.
	IndicatorLed
	// IndicatorButton is used in the firmware DemoProject.
	IndicatorButton
)
