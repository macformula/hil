package macformula

import (
	"context"
	"strconv"

	"github.com/macformula/hil/macformula/pinout"

	"github.com/fatih/color"
	"github.com/macformula/hil/iocontrol"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	// Labels (must be unique)
	_digitalVsAnalogLabel     = "Digital IO vs Analog IO"
	_digitalOutputSelectLabel = "Digital Output"
	_digitalInputSelectLabel  = "Digital Input"
	_analogOutputSelectLabel  = "Analog Output"
	_analogInputSelectLabel   = "Analog Input"
	_digitalLevelSelectLabel  = "Digital Level"

	// Select Options
	_exit          = "--- EXIT ---"
	_return        = "--- RETURN ---"
	_digitalOutput = "Digital Output"
	_digitalInput  = "Digital Input"
	_analogOutput  = "Analog Output"
	_analogInput   = "Analog Input"

	// Other
	_loggerName        = "iocheckout"
	_digitalHigh       = "High"
	_digitalLow        = "Low"
	_voltagePrompt     = "Voltage"
	_lowerVoltageBound = -3.3
	_upperVoltageBound = 3.3
)

var (
	ioValueColor = color.New(color.FgHiBlue, color.Bold)
)

// IoCheckout is a cli tool to control io set out in pinout.go.
type IoCheckout struct {
	l         *zap.Logger
	rev       pinout.Revision
	ioControl iocontrol.IOController

	diPins pinout.DigitalPinout
	doPins pinout.DigitalPinout
	aiPins pinout.AnalogPinout
	aoPins pinout.AnalogPinout

	diPinStrs []string
	doPinStrs []string
	aiPinStrs []string
	aoPinStrs []string

	currentLabel string
}

// NewIoCheckout returns a pointer to an IoCheckout object.
func NewIoCheckout(rev pinout.Revision, ioControl iocontrol.IOController, l *zap.Logger) *IoCheckout {
	return &IoCheckout{
		l:            l.Named(_loggerName),
		rev:          rev,
		ioControl:    ioControl,
		currentLabel: _digitalVsAnalogLabel,
	}
}

func (io *IoCheckout) Open(ctx context.Context) error {
	err := io.ioControl.Open(ctx)
	if err != nil {
		return errors.Wrap(err, "iocontrol open")
	}

	po, ok := pinout.Pinouts[io.rev]
	if !ok {
		return errors.Errorf("Invalid revision %s", io.rev.String())
	}
	io.diPins = po.DigitalInputs
	io.doPins = po.DigitalOutputs
	io.aiPins = po.AnalogInputs
	io.aoPins = po.AnalogOutputs

	return nil
}

// Start starts the cli tool.
func (io *IoCheckout) Start() error {
	var err error

	for {
		switch io.currentLabel {
		case _digitalVsAnalogLabel:
			err = io.handleDigitalVsAnalog()
			if err != nil {
				return errors.Wrap(err, "handle digital vs analog")
			}
		case _digitalInputSelectLabel:
			err = io.handleDigitalInputSelect()
			if err != nil {
				return errors.Wrap(err, "handle digital input select")
			}
		case _digitalOutputSelectLabel:
			err = io.handleDigitalOutputSelect()
			if err != nil {
				return errors.Wrap(err, "handle digital output select")
			}
		case _analogInputSelectLabel:
			err = io.handleAnalogInputSelect()
			if err != nil {
				return errors.Wrap(err, "handle analog input select")
			}
		case _analogOutputSelectLabel:
			err = io.handleAnalogOutputSelect()
			if err != nil {
				return errors.Wrap(err, "handle analog output select")
			}
		case _exit:
			io.l.Info("received exit prompt")

			return nil
		default:
			return errors.New("reached unknown label")
		}
	}
}

func (io *IoCheckout) Close() error {
	err := io.ioControl.Close()
	if err != nil {
		return errors.Wrap(err, "close iocontrol")
	}

	return nil
}

func (io *IoCheckout) promptSelect(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return "", errors.Wrap(err, "run")
	}

	return result, nil
}

func (io *IoCheckout) promptVoltage() (float64, error) {
	prompt := promptui.Prompt{
		Label:    _voltagePrompt,
		Validate: validateVoltage,
	}

	resultStr, err := prompt.Run()
	if err != nil {
		return 0.0, errors.Wrap(err, "run")
	}

	ret, err := strconv.ParseFloat(resultStr, 64)
	if err != nil {
		return 0.0, errors.New("invalid voltage input")
	}

	return ret, nil
}

func (io *IoCheckout) handleDigitalVsAnalog() error {
	digVsAnalogStr, err := io.promptSelect(_digitalVsAnalogLabel, []string{
		_exit,
		_digitalInput,
		_digitalOutput,
		_analogInput,
		_analogOutput,
	})
	if err != nil {
		return errors.Wrap(err, "prompt select")
	}

	switch digVsAnalogStr {
	case _digitalInput:
		io.currentLabel = _digitalInputSelectLabel
	case _digitalOutput:
		io.currentLabel = _digitalOutputSelectLabel
	case _analogInput:
		io.currentLabel = _analogInputSelectLabel
	case _analogOutput:
		io.currentLabel = _analogOutputSelectLabel
	case _exit:
		io.currentLabel = _exit
	default:
		return errors.Errorf("unexpected choice for digital vs analog (%s)", digVsAnalogStr)
	}

	return nil
}

func (io *IoCheckout) digitalInputStrings() []string {
	if io.diPinStrs != nil {
		return io.diPinStrs
	}

	var ret []string

	for digitalIn := range io.diPins {
		ret = append(ret, digitalIn.String())
	}

	io.diPinStrs = ret

	return ret
}

func (io *IoCheckout) digitalOutputStrings() []string {
	if io.doPinStrs != nil {
		return io.doPinStrs
	}

	var ret []string

	for digitalOut := range io.doPins {
		ret = append(ret, digitalOut.String())
	}

	io.doPinStrs = ret

	return ret
}

func (io *IoCheckout) analogInputStrings() []string {
	if io.aiPinStrs != nil {
		return io.aiPinStrs
	}

	var ret []string

	for analogIn := range io.aiPins {
		ret = append(ret, analogIn.String())
	}

	io.aiPinStrs = ret

	return ret
}

func (io *IoCheckout) analogOutputStrings() []string {
	if io.aoPinStrs != nil {
		return io.aoPinStrs
	}

	var ret []string

	for analogOut := range io.aoPins {
		ret = append(ret, analogOut.String())
	}

	io.aoPinStrs = ret

	return ret
}

func (io *IoCheckout) handleDigitalInputSelect() error {
	digitalInputStr, err := io.promptSelect(
		_digitalInputSelectLabel,
		append([]string{_return}, io.digitalInputStrings()...),
	)
	if err != nil {
		return errors.Wrap(err, "prompt select")
	}

	if digitalInputStr == _return {
		io.currentLabel = _digitalVsAnalogLabel

		return nil
	}

	physicalIn, err := pinout.PhysicalIoString(digitalInputStr)
	if err != nil {
		return errors.Wrap(err, "physical io string")
	}

	digitalIn, ok := io.diPins[physicalIn]
	if !ok {
		return errors.Wrap(err, "physical io does not map to digital input")
	}

	lvl, err := io.ioControl.ReadDigital(digitalIn)
	if err != nil {
		return errors.Wrap(err, "read digital")
	}

	_, _ = ioValueColor.Printf("\n%s: %s\n", physicalIn.String(), boolToHighLow(lvl))

	return nil
}

func (io *IoCheckout) handleDigitalOutputSelect() error {
	digitalOutputStr, err := io.promptSelect(
		_digitalOutputSelectLabel,
		append([]string{_return}, io.digitalOutputStrings()...),
	)
	if err != nil {
		return errors.Wrap(err, "prompt select")
	}

	if digitalOutputStr == _return {
		io.currentLabel = _digitalVsAnalogLabel

		return nil
	}

	physicalOut, err := pinout.PhysicalIoString(digitalOutputStr)
	if err != nil {
		return errors.Wrap(err, "physical io string")
	}

	digitalOut, ok := io.doPins[physicalOut]
	if !ok {
		return errors.Wrap(err, "physical io does not map to digital output")
	}

	digitalLevelStr, err := io.promptSelect(
		_digitalLevelSelectLabel,
		[]string{_digitalLow, _digitalHigh},
	)
	if err != nil {
		return errors.Wrap(err, "prompt select")
	}

	var lvl bool

	lvl, err = highLowToBool(digitalLevelStr)
	if err != nil {
		return errors.Wrap(err, "high to low bool")
	}

	err = io.ioControl.WriteDigital(digitalOut, lvl)
	if err != nil {
		return errors.Wrap(err, "set digital")
	}

	_, _ = ioValueColor.Printf("\n%s: %s\n", physicalOut.String(), boolToHighLow(lvl))

	return nil
}

func (io *IoCheckout) handleAnalogInputSelect() error {
	analogInputStr, err := io.promptSelect(_analogInputSelectLabel, append([]string{_return}, io.analogInputStrings()...))
	if err != nil {
		return errors.Wrap(err, "prompt select")
	}

	if analogInputStr == _return {
		io.currentLabel = _digitalVsAnalogLabel

		return nil
	}

	physicalIn, err := pinout.PhysicalIoString(analogInputStr)
	if err != nil {
		return errors.Wrap(err, "physical io string")
	}

	analogIn, ok := io.aiPins[physicalIn]
	if !ok {
		return errors.Wrap(err, "physical io does not map to analog input")
	}

	voltage, err := io.ioControl.ReadVoltage(analogIn)
	if err != nil {
		return errors.Wrap(err, "read voltage")
	}

	_, _ = ioValueColor.Printf("\n%s: %3f V\n", physicalIn.String(), voltage)

	return nil
}

func (io *IoCheckout) handleAnalogOutputSelect() error {
	analogOutputStr, err := io.promptSelect(_analogInputSelectLabel, append([]string{_return}, io.analogOutputStrings()...))
	if err != nil {
		return errors.Wrap(err, "prompt select")
	}

	if analogOutputStr == _return {
		io.currentLabel = _digitalVsAnalogLabel

		return nil
	}

	physicalOut, err := pinout.PhysicalIoString(analogOutputStr)
	if err != nil {
		return errors.Wrap(err, "physical io string")
	}

	analogOut, ok := io.aoPins[physicalOut]
	if !ok {
		return errors.Wrap(err, "physical io does not map to digital input")
	}

	voltage, err := io.promptVoltage()
	if err != nil {
		return errors.Wrap(err, "prompt voltage")
	}

	err = io.ioControl.WriteVoltage(analogOut, voltage)
	if err != nil {
		return errors.Wrap(err, "write voltage")
	}

	_, _ = ioValueColor.Printf("\n%s: %3f V\n", physicalOut.String(), voltage)

	return nil
}

func validateVoltage(input string) error {
	val, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return errors.New("invalid voltage input")
	}

	if val < _lowerVoltageBound || val > _upperVoltageBound {
		return errors.Errorf("invalid voltage provided (%v) must be within bounds [%v, %v]",
			val, _lowerVoltageBound, _upperVoltageBound)
	}

	return nil
}

func highLowToBool(lvl string) (bool, error) {
	if lvl == _digitalHigh {
		return true, nil
	} else if lvl == _digitalLow {
		return false, nil
	}

	return false, errors.Errorf("unknown digital level string (%s)", lvl)
}

func boolToHighLow(b bool) string {
	if b {
		return _digitalHigh
	}

	return _digitalLow
}
