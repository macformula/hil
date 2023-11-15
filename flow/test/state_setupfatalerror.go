package test

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

// SetupFatalErrorState returns an error when Setup is called and when FatalError is called.
type SetupFatalErrorState struct{}

// Timeout returns the state setup and run timeout.
func (s *SetupFatalErrorState) Timeout() time.Duration {
	return time.Minute
}

// Setup executes any necessary setup logic before run.
func (s *SetupFatalErrorState) Setup(_ context.Context) error {
	return errors.New("there has been a setup error")
}

// Name is the name of the state.
func (s *SetupFatalErrorState) Name() string {
	return "setup_fatal_error_state"
}

// Run is the logic that gets executed after setup.
func (s *SetupFatalErrorState) Run(_ context.Context) error {
	return nil
}

// FatalError indicates if any non-recoverable errors have occured.
func (s *SetupFatalErrorState) FatalError() error {
	return errors.New("there has been a fatal setup error")
}
