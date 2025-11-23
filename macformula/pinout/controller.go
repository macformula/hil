package pinout

import (
	"context"

	"go.uber.org/zap"

	"github.com/macformula/hil/iocontrol"
	"github.com/pkg/errors"
)

const _controllerName = "pinout_controller"

// Controller allows for easy control of the I/O's given the current pinout Revision.
type Controller struct {
	l            *zap.Logger
	ioController iocontrol.IOController
	rev          Revision

	digitalOutputs DigitalPinout
	digitalInputs  DigitalPinout
	analogOutputs  AnalogPinout
	analogInputs   AnalogPinout
}

// NewController creates a new pinout controller.
func NewController(rev Revision, ioController iocontrol.IOController, l *zap.Logger) *Controller {
	return &Controller{
		l:            l.Named(_controllerName),
		ioController: ioController,
		rev:          rev,
	}
}

// Open opens the controller and initializes the digital and analog I/O's.
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

// SetDigitalLevel sets the digital level of the given output.
func (c *Controller) SetDigitalLevel(out PhysicalIo, level bool) error {
	digitalOutput, ok := c.digitalOutputs[out]
	if !ok {
		return errors.Errorf("no digital output for physical io (%s) in revision (%s)",
			out.String(), c.rev.String())
	}

	err := c.ioController.WriteDigital(digitalOutput, level)
	if err != nil {
		return errors.Wrap(err, "set digital")
	}

	return nil
}

// ReadDigitalLevel reads the digital level of the given input.
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

// SetVoltage sets the voltage of the given output.
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

// ReadVoltage reads the voltage of the given input.
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
