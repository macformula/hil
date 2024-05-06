package test

import (
	"context"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/macformula/config"
	"time"
)

type PanicState struct {
}

func (p PanicState) Name() string {
	return "panic_state"
}

func (p PanicState) Setup(ctx context.Context) error {
	return nil
}

func (p PanicState) Run(ctx context.Context) error {
	panic("this is the panic state")
}

func (p PanicState) GetResults() map[flow.Tag]any {
	return map[flow.Tag]any{
		config.FirmwareTags.FrontControllerFlashed: true,
	}
}

func (p PanicState) ContinueOnFail() bool {
	return false
}

func (p PanicState) Timeout() time.Duration {
	return 1 * time.Second
}

func (p PanicState) FatalError() error {
	return nil
}
