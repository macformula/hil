package test

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

// RunErrorState returns an error when Run is called
type RunErrorState struct{}

// Timeout returns the state setup and run timeout.
func (r *RunErrorState) Timeout() time.Duration {
	return time.Minute
}

// Setup executes any necessary setup logic before run.
func (r *RunErrorState) Setup(_ context.Context) error {
	return nil
}

// Name is the name of the state.
func (r *RunErrorState) Name() string {
	return "run_error_state"
}

// Run is the logic that gets executed after setup.
func (r *RunErrorState) Run(_ context.Context) error {

	return errors.New("there has been an error")
}

// FatalError indicates if any non-recoverable errors have occured.
func (r *RunErrorState) FatalError() error {
	return nil
}
