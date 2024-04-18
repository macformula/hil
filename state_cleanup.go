package hil

import (
	"context"
	"github.com/macformula/hil/flow"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

const (
	_cleanupStateName    = "cleanup_state"
	_cleanupStateTimeout = 10 * time.Second
)

type CleanupState struct {
	l   *zap.Logger
	app *AppState
}

func NewCleanupState(a *AppState, l *zap.Logger) *CleanupState {
	return &CleanupState{
		l:   l.Named(_cleanupStateName),
		app: a,
	}
}

func (c *CleanupState) Name() string {
	return _cleanupStateName
}

func (c *CleanupState) Setup(_ context.Context) error {
	return nil
}

func (c *CleanupState) Run(ctx context.Context) error {
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

func (c *CleanupState) GetResults() map[flow.Tag]any {
	// No results to return
	return nil
}

func (c *CleanupState) ContinueOnFail() bool {
	return false
}

func (c *CleanupState) Timeout() time.Duration {
	return _cleanupStateTimeout
}

func (c *CleanupState) FatalError() error {
	return nil
}
