package pinout

//go:generate enumer -type=PhysicalIo "physicalio.go"

// PhysicalIo represents an input/output in its physical meaning.
type PhysicalIo int

const (
	UnknownPhysicalIo PhysicalIo = iota

	// Front Controller IOs
	DebugLedEn
	DashboardEn
	HvilLedEn
	BrakeLightEn
	StatusLedEn
	RtdsEn
	AccelPedalPosition1
	AccelPedalPosition2
	SteeringAngle
	StartButtonN
	WheelSpeedLeftA
	WheelSpeedLeftB
	WheelSpeedRightA
	WheelSpeedRightB

	// LV Controller IOs
	MotorControllerPrechargeEn // Digital output from LV controller, read as digital input
	InverterSwitchEn           // Digital output from LV controller, read as digital input
	AccumulatorEn              // Digital output from LV controller, read as digital input
	ShutdownCircuitEn          // Digital output from LV controller, read as digital input
	TsalEn
	RaspiEn
	FrontControllerEn
	SpeedgoatEn
	MotorControllerEn
	ImuGpsEn
	DcdcEn
	DcdcValid

	// Test bench control IOs
	HvCurrentSense               // Current through the accumulator (electrical signal only)
	LvController3v3RefVoltage    // Connected to 3.3V touch point on LV controller
	FrontController3v3RefVoltage // Connected to 3.3V touch point on front controller
	HvilDisable                  // Disables HVIL on the test bench
	HvilFeedback                 // Reads back HVIL circuit voltage on the testbench (max 5V)
	GlvmsDisable                 // Disables power to the test bench

	// Demo Project IOs
	IndicatorLed
	IndicatorButton

	// New constants
	TsalEn
	RaspiEn
	FrontControllerEn
	SpeedgoatEn
	MotorControllerEn
	ImuGpsEn
	DcdcEn
	DcdcValid
)
