package frontcontroller

import (
	"context"

	"go.uber.org/zap"

	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/macformula/cangen/vehcan"
	"github.com/macformula/hil/macformula/pinout"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"
)

const (
	_clientName = "front_controller_client"
)

// Client allows for easy control/query of the front controller.
type Client struct {
	l             *zap.Logger
	pinController *pinout.Controller
	vehBusManager *canlink.BusManager
}

// NewClient creates a new front controller client.
func NewClient(pc *pinout.Controller, veh *canlink.BusManager, l *zap.Logger) *Client {
	return &Client{
		l:             l.Named(_clientName),
		pinController: pc,
		vehBusManager: veh,
	}
}

// CommandContactors sends a command to the BMS/LV controller (pretending to be the front controller) to
// control the contactors.
func (c *Client) CommandContactors(ctx context.Context, hvPositive, hvNegative, precharge bool) error {
	var contactorCommand = vehcan.NewContactorStates()

	contactorCommand.SetPackPositive(utils.BoolToNumeric(hvPositive))
	contactorCommand.SetPackNegative(utils.BoolToNumeric(hvNegative))
	contactorCommand.SetPackPrecharge(utils.BoolToNumeric(precharge))

	err := c.vehBusManager.Send(ctx, contactorCommand)
	if err != nil {
		return errors.Wrap(err, "veh can client send")
	}

	return nil
}

// CommandInverter send a command to the inverters to enable or disable them.
func (c *Client) CommandInverter(ctx context.Context, enable bool) error {
	var inverterCommand = vehcan.NewInverterCommand()

	inverterCommand.SetEnableInverter(enable)

	err := c.vehBusManager.Send(ctx, inverterCommand)
	if err != nil {
		return errors.Wrap(err, "veh can client send")
	}

	return nil
}
