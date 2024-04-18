package hil

import (
	"context"
	"github.com/macformula/hil/flow"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

const (
	_initStateName    = "init_state"
	_initStateTimeout = 10 * time.Second
)

type InitState struct {
	l   *zap.Logger
	app *AppState
}

func NewInitState(a *AppState, l *zap.Logger) *InitState {
	return &InitState{
		l:   l.Named(_initStateName),
		app: a,
	}
}

func (s *InitState) Name() string {
	return _initStateName
}

func (s *InitState) Setup(ctx context.Context) error {
	return nil
}

func (s *InitState) Run(ctx context.Context) error {
	s.app.currProcess = NewProcessInfo()

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

func (s *InitState) GetResults() map[flow.Tag]any {
	// No results for init state.
	return nil
}

func (s *InitState) ContinueOnFail() bool {
	return false
}

func (s *InitState) Timeout() time.Duration {
	return _initStateTimeout
}

func (s *InitState) FatalError() error {
	return nil
}
