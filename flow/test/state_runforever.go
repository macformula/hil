package test

import (
	"context"
	"time"
)

type RunForeverState struct{}

func (r *RunForeverState) Timeout() time.Duration {
	return time.Second
}

func (r *RunForeverState) Setup(_ context.Context) error {
	return nil
}

func (r *RunForeverState) Name() string {
	return "run_forever_state"
}

func (r *RunForeverState) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *RunForeverState) FatalError() error {
	return nil
}
