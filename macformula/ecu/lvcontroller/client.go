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

func (c *Client) ReadDigital(pin pinout.PhysicalIo) (bool, error) {
	lvl, err := c.pinController.ReadDigitalLevel(pinout.TsalEn)
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
