package test

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type RunFatalErrorState struct{}

func (r *RunFatalErrorState) Timeout() time.Duration {
	return time.Minute
}

func (r *RunFatalErrorState) Setup(_ context.Context) error {
	return nil
}

func (r *RunFatalErrorState) Name() string {
	return "run_fatal_error_state"
}

func (r *RunFatalErrorState) Run(_ context.Context) error {
	return errors.New("there has been a run error")
}

func (r *RunFatalErrorState) FatalError() error {
	return errors.New("there has been a fatal run error")
}
