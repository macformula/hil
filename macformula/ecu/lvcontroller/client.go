package lvcontroller

import (
	"go.uber.org/zap"

	"github.com/macformula/hil/macformula/pinout"
	"github.com/pkg/errors"
)

const (
	_clientName = "lv_controller_client"
)

// Client allows for easy control/query of the lv controller.
type Client struct {
	l             *zap.Logger
	pinController *pinout.Controller
}

// NewClient creates a new lv controller client.
func NewClient(pc *pinout.Controller, l *zap.Logger) *Client {
	return &Client{
		l:             l.Named(_clientName),
		pinController: pc,
	}
}

// ReadTsalEn reads the TSAL enable pin.
func (c *Client) ReadTsalEn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.TsalEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

// ReadRaspiOn reads the Raspi enable pin.
func (c *Client) ReadRaspiOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.RaspiEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

// ReadFrontControllerOn reads the front controller enable pin.
func (c *Client) ReadFrontControllerOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.FrontControllerEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

// ReadSpeedgoatOn reads the Speedgoat enable pin.
func (c *Client) ReadSpeedgoatOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.SpeedgoatEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

// ReadAccumulatorOn reads the accumulator enable pin.
func (c *Client) ReadAccumulatorOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.AccumulatorEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

// ReadMotorControllerPrechargeOn reads the motor controller precharge enable pin.
func (c *Client) ReadMotorControllerPrechargeOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.MotorControllerPrechargeEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

// ReadMotorControllerOn reads the motor controller enable pin.
func (c *Client) ReadMotorControllerOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.MotorControllerEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

// ReadImuGpsOn reads the IMU/GPS enable pin.
func (c *Client) ReadImuGpsOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.ImuGpsEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

// ReadShutdownCircuitOn reads the shutdown circuit enable pin.
func (c *Client) ReadShutdownCircuitOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.ShutdownCircuitEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

// ReadInverterSwitchOn reads the inverter switch enable pin.
func (c *Client) ReadInverterSwitchOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.InverterSwitchEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

// ReadDcdcOn reads the DCDC enable pin.
func (c *Client) ReadDcdcOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.DcdcEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

// SetDcdcValid sets the DCDC valid pin.
func (c *Client) SetDcdcValid(b bool) error {
	err := c.pinController.SetDigitalLevel(pinout.DcdcValid, b)
	if err != nil {
		return errors.Wrap(err, "set digital level")
	}

	return nil
}
