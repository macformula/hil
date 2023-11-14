package test

import (
	"context"
	"time"
)

type SetupForeverState struct{}

func (s *SetupForeverState) Timeout() time.Duration {
	return time.Second
}

func (s *SetupForeverState) Setup(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *SetupForeverState) Name() string {
	return "setup_forever_state"
}

func (s *SetupForeverState) Run(_ context.Context) error {
	return nil
}

func (s *SetupForeverState) FatalError() error {
	return nil
}
