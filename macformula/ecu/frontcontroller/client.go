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

type Client struct {
	l             *zap.Logger
	pinController *pinout.Controller
	vehCanClient  *canlink.CanClient
}

func (c *Client) CommandContactors(ctx context.Context, hvPositive, hvNegative, precharge bool) error {
	var contactorCommand = vehcan.NewContactor_States()

	contactorCommand.SetPack_Positive(utils.BoolToNumeric(hvPositive))
	contactorCommand.SetPack_Negative(utils.BoolToNumeric(hvNegative))
	contactorCommand.SetPack_Precharge(utils.BoolToNumeric(precharge))

	err := c.vehCanClient.Send(ctx, contactorCommand)
	if err != nil {
		return errors.Wrap(err, "veh can client send")
	}

	return nil
}

func (c *Client) CommandInverter(ctx context.Context, enable bool) error {
	panic("implement me")
}
