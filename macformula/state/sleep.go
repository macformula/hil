package state

import (
	"context"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/macformula/config"
	"time"
)

const _sleepStateName = "sleep_state"

// sleep sleeps in the Run function for a specified amount of time.
type sleep struct {
	sleepTime time.Duration
}

func newSleep(sleepTime time.Duration) *sleep {
	return &sleep{sleepTime: sleepTime}
}

func (s *sleep) GetResults() map[flow.Tag]any {
	return map[flow.Tag]any{
		config.TestTags.TestTag1: "no value",
	}
}

func (s *sleep) ContinueOnFail() bool {
	return true
}

// Timeout returns the state setup and run timeout.
func (s *sleep) Timeout() time.Duration {
	return s.sleepTime + 1*time.Second
}

// Setup executes any necessary setup logic before run.
func (s *sleep) Setup(_ context.Context) error {
	return nil
}

// Name is the name of the state.
func (s *sleep) Name() string {
	return _sleepStateName
}

// Run is the logic that gets executed after setup.
func (s *sleep) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(s.sleepTime):
		return nil
	}
}

// FatalError indicates if any non-recoverable errors have occurred.
func (s *sleep) FatalError() error {
	return nil
}
