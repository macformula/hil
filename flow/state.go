package flow

import (
	"context"
	"time"
)

// Sequence is a list of states that will be run in the order provided.
type Sequence []State

// State is a set of logic that gets executed as a part of a Sequence.
type State interface {
	// Name of the state, should be in lower_snake_case.
	Name() string
	// Setup will be called before Run.
	Setup(ctx context.Context) error
	// Run should be responsible for executing the main logic of the State.
	Run(ctx context.Context) error
	// GetResults will be called after Run. It returns a map of tag-value pairs.
	GetResults() map[Tag]any
	// ContinueOnFail indicates whether the Sequencer should continue running next states if the State fails.
	ContinueOnFail() bool
	// Timeout is the max duration of time the state will be allowed to run for.
	Timeout() time.Duration
	// FatalError indicates that there is likely some hardware failure, or some other type of non-recoverable error.
	FatalError() error
}
