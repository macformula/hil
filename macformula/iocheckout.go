package macformula

import (
	"github.com/fatih/color"
	"github.com/macformula/hil/iocontrol"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	// Labels: MUST BE UNIQUE!!
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
	_digitalHigh = "High"
	_digitalLow  = "Low"
)

type IoCheckout struct {
	l         *zap.Logger
	ioControl *iocontrol.IOControl

	diPins DigitalPinout
	doPins DigitalPinout
	aiPins AnalogPinout
	aoPins AnalogPinout

	diPinStrs []string
	doPinStrs []string
	aiPinStrs []string
	aoPinStrs []string

	currentLabel string
}

func NewIoCheckout(rev Revision, ioControl *iocontrol.IOControl, l *zap.Logger) *IoCheckout {
	return &IoCheckout{
		l:            l.Named("iocheckout"),
		ioControl:    ioControl,
		diPins:       _revisionDigitalInputPinout[rev],
		doPins:       _revisionDigitalOutputPinout[rev],
		aiPins:       _revisionAnalogInputPinout[rev],
		aoPins:       _revisionAnalogOutputPinout[rev],
		currentLabel: _digitalVsAnalogLabel,
	}
}

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
		default:
			return errors.New("reached unknown label")
		}
	}

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
	digtialInputStr, err := io.promptSelect(_digitalInputSelectLabel, append([]string{_return}, io.digitalInputStrings()...))
	if err != nil {
		return errors.Wrap(err, "prompt select")
	}

	if digtialInputStr == _return {
		io.currentLabel = _digitalVsAnalogLabel

		return nil
	}

	physicalIn, err := PhysicalIoString(digtialInputStr)
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

	if lvl {
		_, _ = color.New(color.FgHiBlue, color.Bold).Printf("\n%s: %s\n", physicalIn.String(), _digitalHigh)
	} else {
		_, _ = color.New(color.FgHiBlue, color.Bold).Printf("\n%s: %s\n", physicalIn.String(), _digitalLow)
	}

	return nil
}

func (io *IoCheckout) handleAnalogOutputSelect() error {
	digtialInputStr, err := io.promptSelect(_digitalInputSelectLabel, append([]string{_return}, io.digitalInputStrings()...))
	if err != nil {
		return errors.Wrap(err, "prompt select")
	}

	if digtialInputStr == _return {
		io.currentLabel = _digitalVsAnalogLabel

		return nil
	}

	physicalIn, err := PhysicalIoString(digtialInputStr)
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

	if lvl {
		_, _ = color.New(color.FgHiBlue, color.Bold).Printf("\n%s: %s\n", physicalIn.String(), _digitalHigh)
	} else {
		_, _ = color.New(color.FgHiBlue, color.Bold).Printf("\n%s: %s\n", physicalIn.String(), _digitalLow)
	}

	return nil
}

func (io *IoCheckout) handleAnalogInputSelect() error {
	return errors.New("unimplemented")
}

func (io *IoCheckout) handleDigitalOutputSelect() error {
	digtialOutputStr, err := io.promptSelect(
		_digitalOutputSelectLabel,
		append([]string{_return}, io.digitalOutputStrings()...),
	)
	if err != nil {
		return errors.Wrap(err, "prompt select")
	}

	if digtialOutputStr == _return {
		io.currentLabel = _digitalVsAnalogLabel

		return nil
	}

	physicalOut, err := PhysicalIoString(digtialOutputStr)
	if err != nil {
		return errors.Wrap(err, "physical io string")
	}

	digitalOut, ok := io.doPins[physicalOut]
	if !ok {
		return errors.Wrap(err, "physical io does not map to digital output")
	}

	digtialLevelStr, err := io.promptSelect(
		_digitalLevelSelectLabel,
		[]string{_digitalLow, _digitalHigh},
	)
	if err != nil {
		return errors.Wrap(err, "prompt select")
	}

	var lvl bool

	lvl, err = highLowToBool(digtialLevelStr)
	if err != nil {
		return errors.Wrap(err, "high to low bool")
	}

	err = io.ioControl.SetDigital(digitalOut, lvl)
	if err != nil {
		return errors.Wrap(err, "read digital")
	}

	if lvl {
		_, _ = color.New(color.FgHiBlue, color.Bold).Printf("\n%s: %s\n", physicalOut.String(), _digitalHigh)
	} else {
		_, _ = color.New(color.FgHiBlue, color.Bold).Printf("\n%s: %s\n", physicalOut.String(), _digitalLow)
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
