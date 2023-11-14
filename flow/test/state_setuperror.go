package test

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

type SetupErrorState struct{}

func (s *SetupErrorState) Timeout() time.Duration {
	return time.Minute
}

func (s *SetupErrorState) Setup(_ context.Context) error {
	return errors.New("there has been an error")
}

func (s *SetupErrorState) Name() string {
	return "setup_error_state"
}

func (s *SetupErrorState) Run(_ context.Context) error {
	return nil
}

func (s *SetupErrorState) FatalError() error {
	return nil
}
