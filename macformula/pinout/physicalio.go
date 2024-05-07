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
	// AccumulatorCurrent determines the current that is passing through the battery pack (accumulator).
	AccumulatorCurrent
	// LvController3v3RefVoltage is connected to the 3.3V touch point on the lv controller.
	LvController3v3RefVoltage
	FrontController3v3RefVoltage
	// IndicatorLed is used in the firmware DemoProject.
	IndicatorLed
	// IndicatorButton is used in the firmware DemoProject.
	IndicatorButton
	// HvCurrentSense does not go into a microcontroller, it is an electrical only signal.
	HvCurrentSense
	MotorControllerPrechargeEn
	InverterEn
	ShutdownCircuitEn
	AccumulatorEn

	// IoFcCheckoutProject:
	GlvmsDisable
	DebugLedEn
	DashboardEn
	HvilLedEn
	BrakeLightEn
	StatusLedEn
	RtdsEn
	AccelPedalPosition1
	AccelPedalPosition2
	BrakePedalPosition1
	BrakePedalPosition2
	SuspensionTravel1
	SuspensionTravel2
	HvilFeedback
	SteeringAngle
	StartButtonN
	WheelSpeedLeftA
	WheelSpeedLeftB
	WheelSpeedRightA
	WheelSpeedRightB
	WaitForStart
	HvilDisable

	SgTestOutput
	SgTestInput
)
