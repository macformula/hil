package test

import (
	"context"
	"github.com/macformula/hil"
	"github.com/macformula/hil/flow"
	"time"
)

// SleepState sleeps in the Run function for a specified amount of time.
type SleepState struct {
	SleepTime time.Duration
}

func (s *SleepState) GetResults() map[flow.Tag]any {
	return map[flow.Tag]any{
		hil.FwTags.FrontControllerFlashed: true,
		hil.FwTags.TmsFlashed:             true,
	}
}

func (s *SleepState) ContinueOnFail() bool {
	return true
}

// Timeout returns the state setup and run timeout.
func (s *SleepState) Timeout() time.Duration {
	return s.SleepTime + 1*time.Second
}

// Setup executes any necessary setup logic before run.
func (s *SleepState) Setup(_ context.Context) error {
	return nil
}

// Name is the name of the state.
func (s *SleepState) Name() string {
	return "sleep_state"
}

// Run is the logic that gets executed after setup.
func (s *SleepState) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(s.SleepTime):
		return nil
	}
}

// FatalError indicates if any non-recoverable errors have occurred.
func (s *SleepState) FatalError() error {
	return nil
}
