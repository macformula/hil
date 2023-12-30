package orchestrator

import (
	"context"
	"github.com/macformula/hil/flow"

	"github.com/ethereum/go-ethereum/event"
)

// Sequencer is responsible for managining execution of a sequence of test states
type Sequencer interface {
	// SubscribeToProgress subscribes to the progress of the Sequencer across its Sequence runs.
	// The Progress channel gets updated whenever there is new information available.
	SubscribeToProgress(progCh chan flow.Progress) event.Subscription
	// Run will run the sequence provided. FatalError must be called after Run to check for any non-recoverable errors.
	Run(ctx context.Context, seq flow.Sequence) error
	// FatalError indicates that there is an error that requires intervention.
	// The Sequencer will stop executing all remaining states in the Sequence if it encounters a fatal error.
	FatalError() error
	// ResetFatalError sets the fatal error to nil.
	ResetFatalError()
}
