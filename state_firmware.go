package hil

import (
	"context"
	"errors"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/utils"
	"time"
)

const (
	_stateName = "firmware"
)

type Firmware struct {
	err       utils.ResettableError
	FcFlashed bool
}

func (f Firmware) Name() string {
	//TODO implement me
	panic("implement me")
}

func (f Firmware) Setup(ctx context.Context) error {
	//TODO implement me
	panic("implement me")

	return errors.New("divisor cant be 0")
}

func (f Firmware) Run(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (f Firmware) GetResults() map[flow.Tag]any {
	panic("implement me")
}

func (f Firmware) ContinueOnFail() bool {
	//TODO implement me
	panic("implement me")
}

func (f Firmware) Timeout() time.Duration {
	//TODO implement me
	panic("implement me")
}

func (f Firmware) FatalError() error {
	//TODO implement me
	panic("implement me")
}
