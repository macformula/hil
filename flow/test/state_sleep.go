package test

import (
	"context"
	"time"
)

type SleepState struct {
	sleepTime time.Duration
}

func (s *SleepState) Timeout() time.Duration {
	return s.sleepTime + 5*time.Second
}

func (s *SleepState) Setup(_ context.Context) error {
	return nil
}

func (s *SleepState) Name() string {
	return "sleep_state"
}

func (s *SleepState) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(s.sleepTime):
		return nil
	}
}

func (s *SleepState) FatalError() error {
	return nil
}
