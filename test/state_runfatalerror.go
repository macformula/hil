package test

import (
	"context"
	"github.com/macformula/hil"
	"github.com/macformula/hil/flow"
	"time"

	"github.com/pkg/errors"
)

// RunFatalErrorState returns an error when Run is called and when FatalError is called.
type RunFatalErrorState struct{}

func (r *RunFatalErrorState) GetResults() map[flow.Tag]any {
	return map[flow.Tag]any{
		hil.FwTags.FrontControllerFlashed: true,
		hil.FwTags.TmsFlashed:             true,
	}
}

func (r *RunFatalErrorState) ContinueOnFail() bool {
	return true
}

// Timeout returns the state setup and run timeout.
func (r *RunFatalErrorState) Timeout() time.Duration {
	return time.Minute
}

// Setup executes any necessary setup logic before run.
func (r *RunFatalErrorState) Setup(_ context.Context) error {
	return nil
}

// Name is the name of the state.
func (r *RunFatalErrorState) Name() string {
	return "run_fatal_error_state"
}

// Run is the logic that gets executed after setup.
func (r *RunFatalErrorState) Run(_ context.Context) error {
	return errors.New("there has been a run error")
}

// FatalError indicates if any non-recoverable errors have occured.
func (r *RunFatalErrorState) FatalError() error {
	return errors.New("there has been a fatal run error")
}
