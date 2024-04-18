package state

import (
	"context"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/macformula"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

const (
	_cleanupStateName    = "cleanup_state"
	_cleanupStateTimeout = 10 * time.Second
)

type cleanup struct {
	l   *zap.Logger
	app *macformula.App
}

func newCleanup(a *macformula.App, l *zap.Logger) *cleanup {
	return &cleanup{
		l:   l.Named(_cleanupStateName),
		app: a,
	}
}

func (c *cleanup) Name() string {
	return _cleanupStateName
}

func (c *cleanup) Setup(_ context.Context) error {
	return nil
}

func (c *cleanup) Run(ctx context.Context) error {
	err := c.app.VehCanTracer.StopTrace()
	if err != nil {
		return errors.Wrap(err, "stop trace (veh)")
	}

	err = c.app.PtCanTracer.StopTrace()
	if err != nil {
		return errors.Wrap(err, "stop trace (pt)")
	}

	return nil
}

func (c *cleanup) GetResults() map[flow.Tag]any {
	// No results to return
	return nil
}

func (c *cleanup) ContinueOnFail() bool {
	return false
}

func (c *cleanup) Timeout() time.Duration {
	return _cleanupStateTimeout
}

func (c *cleanup) FatalError() error {
	return nil
}