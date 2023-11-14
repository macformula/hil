package test

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

type SetupFatalErrorState struct{}

func (s *SetupFatalErrorState) Timeout() time.Duration {
	return time.Minute
}

func (s *SetupFatalErrorState) Setup(_ context.Context) error {
	return errors.New("there has been a setup error")
}

func (s *SetupFatalErrorState) Name() string {
	return "setup_fatal_error_state"
}

func (s *SetupFatalErrorState) Run(_ context.Context) error {
	return nil
}

func (s *SetupFatalErrorState) FatalError() error {
	return errors.New("there has been a fatal setup error")
}
