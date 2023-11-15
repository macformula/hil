package test

import (
	"context"
	"time"
)

// RunForeverState hangs until the context deadline expires in the Run function.
type RunForeverState struct{}

// Timeout returns the state setup and run timeout.
func (r *RunForeverState) Timeout() time.Duration {
	return time.Second
}

// Setup executes any necessary setup logic before run.
func (r *RunForeverState) Setup(_ context.Context) error {
	return nil
}

// Name is the name of the state.
func (r *RunForeverState) Name() string {
	return "run_forever_state"
}

// Run is the logic that gets executed after setup.
func (r *RunForeverState) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// FatalError indicates if any non-recoverable errors have occured.
func (r *RunForeverState) FatalError() error {
	return nil
}
