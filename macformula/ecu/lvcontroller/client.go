package lvcontroller

import (
	"github.com/macformula/hil/macformula/pinout"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_clientName = "lv_controller_client"
)

type Client struct {
	l             *zap.Logger
	pinController *pinout.Controller
}

func NewClient(pc *pinout.Controller, l *zap.Logger) *Client {
	return &Client{
		l:             l.Named(_clientName),
		pinController: pc,
	}
}

func (c *Client) ReadTsalEn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.TsalEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

func (c *Client) ReadRaspiOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.RaspiEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}
func (c *Client) ReadFrontControllerOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.FrontControllerEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

func (c *Client) ReadSpeedgoatOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.SpeedgoatEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

func (c *Client) ReadAccumulatorOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.AccumulatorEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

func (c *Client) ReadMotorControllerPrechargeOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.MotorControllerPrechargeEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

func (c *Client) ReadMotorControllerPrechargeOff() (bool, error) {
	lvl, err := c.ReadMotorControllerPrechargeOn()
	return !lvl, err
}

func (c *Client) ReadMotorControllerOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.MotorControllerEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

func (c *Client) ReadImuGpsOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.ImuGpsEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

func (c *Client) ReadShutdownCircuitOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.ShutdownCircuitEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}

func (c *Client) ReadInverterSwitchOn() (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.InverterEn)
	if err != nil {
		return false, errors.Wrap(err, "read digital level")
	}

	return lvl, nil
}
