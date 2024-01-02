package test

import (
	"context"
	"errors"
	"github.com/macformula/hil"
	"github.com/macformula/hil/flow"
	"time"
)

const (
	_name = "do_nothing_state"
)

// DoNothingState does nothing
type DoNothingState struct {
	setupCalled bool
}

func (d *DoNothingState) GetResults() map[flow.Tag]any {
	return map[flow.Tag]any{
		hil.FwTags.FrontControllerFlashed: true,
		hil.FwTags.TmsFlashed:             true,
	}
}

func (d *DoNothingState) ContinueOnFail() bool {
	return true
}

// Timeout returns the state setup and run timeout.
func (d *DoNothingState) Timeout() time.Duration {
	return time.Minute
}

// Setup executes any necessary setup logic before run.
func (d *DoNothingState) Setup(ctx context.Context) error {
	d.setupCalled = true

	return nil
}

// Name is the name of the state.
func (d *DoNothingState) Name() string {
	return "do_nothing_state"
}

// Run is the logic that gets executed after setup.
func (d *DoNothingState) Run(_ context.Context) error {
	if !d.setupCalled {
		return errors.New("setup never called")
	}

	return nil
}

// FatalError indicates if any non-recoverable errors have occured.
func (d *DoNothingState) FatalError() error {
	return nil
}
