package test

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

// SetupErrorState returns an error when Run is called
type SetupErrorState struct{}

// Timeout returns the state setup and run timeout.
func (s *SetupErrorState) Timeout() time.Duration {
	return time.Minute
}

// Setup executes any necessary setup logic before run.
func (s *SetupErrorState) Setup(_ context.Context) error {
	return errors.New("there has been an error")
}

// Name is the name of the state.
func (s *SetupErrorState) Name() string {
	return "setup_error_state"
}

// Run is the logic that gets executed after setup.
func (s *SetupErrorState) Run(_ context.Context) error {
	return nil
}

// FatalError indicates if any non-recoverable errors have occured.
func (s *SetupErrorState) FatalError() error {
	return nil
}
