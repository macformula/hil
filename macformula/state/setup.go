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
	_initStateName    = "init_state"
	_initStateTimeout = 10 * time.Second
)

type setup struct {
	l   *zap.Logger
	app *macformula.App
}

func newSetup(a *macformula.App, l *zap.Logger) *setup {
	return &setup{
		l:   l.Named(_initStateName),
		app: a,
	}
}

func (s *setup) Name() string {
	return _initStateName
}

func (s *setup) Setup(ctx context.Context) error {
	return nil
}

func (s *setup) Run(ctx context.Context) error {
	s.app.CurrProcess = macformula.NewProcessInfo()

	err := s.app.VehCanTracer.StartTrace(ctx)
	if err != nil {
		return errors.Wrap(err, "start trace (veh)")
	}

	err = s.app.VehCanTracer.StartTrace(ctx)
	if err != nil {
		return errors.Wrap(err, "start trace (pt)")
	}

	return nil
}

func (s *setup) GetResults() map[flow.Tag]any {
	// No results for init state.
	return nil
}

func (s *setup) ContinueOnFail() bool {
	return false
}

func (s *setup) Timeout() time.Duration {
	return _initStateTimeout
}

func (s *setup) FatalError() error {
	return nil
}
