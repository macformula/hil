package pinout

import (
	"context"
	"github.com/macformula/hil/iocontrol"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const _controllerName = "pinout_controller"

type Controller struct {
	l            *zap.Logger
	ioController *iocontrol.IOControl
	rev          Revision

	digitalOutputs DigitalPinout
	digitalInputs  DigitalPinout
	analogOutputs  AnalogPinout
	analogInputs   AnalogPinout
}

func NewController(rev Revision, ioController *iocontrol.IOControl, l *zap.Logger) *Controller {
	return &Controller{
		l:            l.Named(_controllerName),
		ioController: ioController,
		rev:          rev,
	}
}

func (c *Controller) Open(_ context.Context) error {
	var err error

	c.digitalOutputs, err = GetDigitalOutputs(c.rev)
	if err != nil {
		return errors.Wrap(err, "get digital outputs")
	}

	c.digitalInputs, err = GetDigitalInputs(c.rev)
	if err != nil {
		return errors.Wrap(err, "get digital inputs")
	}

	c.analogOutputs, err = GetAnalogOutputs(c.rev)
	if err != nil {
		return errors.Wrap(err, "get analog outputs")
	}

	c.analogInputs, err = GetAnalogInputs(c.rev)
	if err != nil {
		return errors.Wrap(err, "get analog inputs")
	}

	return nil
}

func (c *Controller) SetDigitalLevel(out PhysicalIo, level bool) error {
	digitalOutput, ok := c.digitalOutputs[out]
	if !ok {
		return errors.Errorf("no digital output for physical io (%s) in revision (%s)",
			out.String(), c.rev.String())
	}

	err := c.ioController.SetDigital(digitalOutput, level)
	if err != nil {
		return errors.Wrap(err, "set digital")
	}

	return nil
}

func (c *Controller) ReadDigitalLevel(in PhysicalIo) (bool, error) {
	digitalInput, ok := c.digitalInputs[in]
	if !ok {
		return false, errors.Errorf("no digital input for physical io (%s) in revision (%s)",
			in.String(), c.rev.String())
	}

	level, err := c.ioController.ReadDigital(digitalInput)
	if err != nil {
		return false, errors.Wrap(err, "read digital")
	}

	return level, nil
}

func (c *Controller) SetVoltage(out PhysicalIo, voltage float64) error {
	analogOutput, ok := c.analogOutputs[out]
	if !ok {
		return errors.Errorf("no analog output for physical io (%s) in revision (%s)",
			out.String(), c.rev.String())
	}

	err := c.ioController.WriteVoltage(analogOutput, voltage)
	if err != nil {
		return errors.Wrap(err, "set analog")
	}

	return nil
}

func (c *Controller) ReadVoltage(in PhysicalIo) (float64, error) {
	analogInput, ok := c.analogInputs[in]
	if !ok {
		return 0.0, errors.Errorf("no analog inputs for physical io (%s) in revision (%s)",
			in.String(), c.rev.String())
	}

	voltage, err := c.ioController.ReadVoltage(analogInput)
	if err != nil {
		return 0, errors.Wrap(err, "read voltage")
	}

	return voltage, nil
}
