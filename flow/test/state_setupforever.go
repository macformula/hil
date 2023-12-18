package test

import (
	"context"
	"time"
)

// SetupForeverState hangs forever until the context deadline expires in the Setup function.
type SetupForeverState struct{}

// Timeout returns the state setup and run timeout.
func (s *SetupForeverState) Timeout() time.Duration {
	return time.Second
}

// Setup executes any necessary setup logic before run.
func (s *SetupForeverState) Setup(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Name is the name of the state.
func (s *SetupForeverState) Name() string {
	return "setup_forever_state"
}

// Run is the logic that gets executed after setup.
func (s *SetupForeverState) Run(_ context.Context) error {
	return nil
}

// FatalError indicates if any non-recoverable errors have occured.
func (s *SetupForeverState) FatalError() error {
	return nil
}
