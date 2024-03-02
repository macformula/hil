package flow

import (
	"context"
	"github.com/google/uuid"
	"io"
	"time"
)

// ResultProcessorIface will be used to get pass/fail statuses on tags.
type ResultProcessorIface interface {
	io.Closer
	// Open will be called at the start of the app.
	Open(context.Context) error
	// SubmitTag will return the passing status of a given tag.
	SubmitTag(ctx context.Context, tagId string, value any) (bool, error)
	// CompleteTest will signal the result processor that a test has been completed. The overall pass/fail is returned.
	CompleteTest(ctx context.Context, testId uuid.UUID, sequenceName string) (bool, error)
	// SubmitError will be stored by the result processor and should make the sequence an overall fail.
	SubmitError(ctx context.Context, err error) error
}

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
