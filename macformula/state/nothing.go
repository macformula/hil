package state

import (
	"context"
	"errors"
	"time"

	"github.com/macformula/hil/flow"
)

const (
	_name = "do_nothing_state"
)

// DoNothingState does nothing
type nothing struct {
	setupCalled bool
}

func newNothing() *nothing {
	return &nothing{setupCalled: false}
}

func (n *nothing) GetResults() map[flow.Tag]any {
	return map[flow.Tag]any{}
}

func (n *nothing) ContinueOnFail() bool {
	return true
}

// Timeout returns the state setup and run timeout.
func (n *nothing) Timeout() time.Duration {
	return time.Minute
}

// Setup executes any necessary setup logic before run.
func (n *nothing) Setup(ctx context.Context) error {
	n.setupCalled = true

	return nil
}

// Name is the name of the state.
func (n *nothing) Name() string {
	return _name
}

// Run is the logic that gets executed after setup.
func (n *nothing) Run(_ context.Context) error {
	if !n.setupCalled {
		return errors.New("setup never called")
	}

	return nil
}

// FatalError indicates if any non-recoverable errors have occured.
func (n *nothing) FatalError() error {
	return nil
}
