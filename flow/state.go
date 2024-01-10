package flow

import (
	"context"
	"time"
)

// State is a set of logic that gets executed as a part of a Sequence.
type State interface {
	// Name of the state, should be in lower_snake_case
	Name() string
	// Setup will be called before run, an error here will result in a fatal error in the sequencer
	Setup(ctx context.Context) error
	// Run should be responsible for executing the main logic of the state
	Run(ctx context.Context) error
	// Timeout is the max duration of time the state will be allowed to run for
	Timeout() time.Duration
	// FatalError indicates that there is likely some hardware failure, or some other type of non-recoverable error
	FatalError() error
}
