package frontcontroller

import (
	"context"
	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/macformula/cangen/vehcan"
	"github.com/macformula/hil/macformula/pinout"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_clientName = "front_controller_client"
)

type Client struct {
	l             *zap.Logger
	pinController *pinout.Controller
	vehCanClient  *canlink.CanClient
}

func NewClient(pc *pinout.Controller, veh *canlink.CanClient, l *zap.Logger) *Client {
	return &Client{
		l:             l.Named(_clientName),
		pinController: pc,
		vehCanClient:  veh,
	}
}

func (c *Client) CommandContactors(ctx context.Context, hvPositive, hvNegative, precharge bool) error {
	var contactorCommand = vehcan.NewContactorStates()

	contactorCommand.SetPackPositive(utils.BoolToNumeric(hvPositive))
	contactorCommand.SetPackNegative(utils.BoolToNumeric(hvNegative))
	contactorCommand.SetPackPrecharge(utils.BoolToNumeric(precharge))

	err := c.vehCanClient.Send(ctx, contactorCommand)
	if err != nil {
		return errors.Wrap(err, "veh can client send")
	}

	return nil
}

func (c *Client) CommandInverter(ctx context.Context, enable bool) error {
	var inverterCommand = vehcan.NewInverterCommand()

	inverterCommand.SetEnableInverter(enable)

	err := c.vehCanClient.Send(ctx, inverterCommand)
	if err != nil {
		return errors.Wrap(err, "veh can client send")
	}

	return nil
}
