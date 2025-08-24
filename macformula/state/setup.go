package state

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/macformula"
	"go.uber.org/zap"
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

	sequenceResultsDir := filepath.Join(s.app.Config.ResultsDir, time.Now().Format("15:04:05.0000"))
	os.Mkdir(sequenceResultsDir, 0755)

	if s.app.WithVcan {
		s.app.VehCanTracer.SetTraceDir(sequenceResultsDir)
		s.app.PtCanTracer.SetTraceDir(sequenceResultsDir)

		s.app.VehBusManager.Register(s.app.VehCanTracer)
		s.app.PtBusManager.Register(s.app.PtCanTracer)

		s.app.VehBusManager.Start(ctx)
		s.app.PtBusManager.Start(ctx)
	}

	s.app.ResultsProcessor.SetReportsDir(sequenceResultsDir)

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
