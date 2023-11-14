package test

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

type RunErrorState struct{}

func (r *RunErrorState) Timeout() time.Duration {
	return time.Minute
}

func (r *RunErrorState) Setup(_ context.Context) error {
	return nil
}

func (r *RunErrorState) Name() string {
	return "run_error_state"
}

func (r *RunErrorState) Run(_ context.Context) error {

	return errors.New("there has been an error")
}

func (r *RunErrorState) FatalError() error {
	return nil
}
