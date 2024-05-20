package pinout

// PhysicalIo represents an input output in its physical meaning.
type PhysicalIo int

//go:generate enumer -type=PhysicalIo "physicalio.go"

const (
	UnknownPhysicalIo PhysicalIo = iota

	// Front Controller

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
	SteeringAngle
	StartButtonN
	WheelSpeedLeftA
	WheelSpeedLeftB
	WheelSpeedRightA
	WheelSpeedRightB

	// Lv Controller

	// MotorControllerPrechargeEn is a digital output coming from the lv controller. It is read into our system as a
	// digital input. It enables precharge of the motor controller.
	MotorControllerPrechargeEn
	// InverterEn is a digital output coming from the lv controller. It is read into our system as a digital input. It
	// enables power to the motor inverters.
	InverterEn
	// AccumulatorEn is a digital output coming from the lv controller. It is read into our system as a digital input.
	// This line enables power to the accumulator.
	AccumulatorEn
	// ShutdownCircuitEn is a digital output coming from the lv controller. It is read into our system as a digital
	// input. This is the line that gives the HVIL circuit power.
	ShutdownCircuitEn

	// Note the following are not connected to testbench.
	TsalEn
	RaspiEn
	FrontControllerEn
	SpeedgoatEn
	MotorControllerEn
	ImuGpsEn
	DcdcValid
	DcdcEn
	DcdcEnLed
	PowerTrainPumpEn
	PowertrainFanEn

	// Test bench control

	// HvCurrentSense represents the current going through the accumulator. It does not go into a microcontroller; it is
	// an electrical only signal.
	HvCurrentSense
	// LvController3v3RefVoltage is connected to the 3.3V touch point on the lv controller.
	LvController3v3RefVoltage
	// FrontController3v3RefVoltage is connected to the 3.3V touch point on the front controller.
	FrontController3v3RefVoltage
	// HvilDisable will disable hvil on the test bench.
	HvilDisable
	// HvilFeedback will read back the voltage of the HVIL circuit on the testbench. Pinned to a maximum of 5V.
	HvilFeedback
	// GlvmsDisable will disable power to the test bench.
	GlvmsDisable

	// Demo Project IO

	// IndicatorLed is used in the firmware DemoProject.
	IndicatorLed
	// IndicatorButton is used in the firmware DemoProject.
	IndicatorButton
)
