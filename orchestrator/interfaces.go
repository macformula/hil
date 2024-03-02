package orchestrator

import (
	"context"
	"io"

	"github.com/ethereum/go-ethereum/event"
	"github.com/macformula/hil/flow"
)

// SequencerIface is responsible for managing execution of a sequence of test states.
type SequencerIface interface {
	io.Closer
	// Open sets up the Sequencer.
	Open(ctx context.Context) error
	// SubscribeToProgress subscribes to the progress of the Sequencer across its Sequence runs.
	SubscribeToProgress(progCh chan flow.Progress) event.Subscription
	// Run will run the sequence provided. FatalError must be called after Run to check for any non-recoverable errors.
	Run(context.Context, flow.Sequence, chan struct{}, TestId) (bool, []flow.Tag, []error, error)
	// FatalError indicates that there is an error that requires intervention.
	FatalError() error
	// ResetFatalError sets the fatal error to nil.
	ResetFatalError()
}

// DispatcherIface is responsible for commanding start of execution.
type DispatcherIface interface {
	io.Closer
	// Name of the Dispatcher
	Name() string
	// Open will be called on orchestrator open
	Open(context.Context) error
	// Start signal is sent by the dispatcher to the orchestrator to start a test sequence.
	Start() <-chan StartSignal
	// CancelTest will cancel execution of the test with the given ID.
	CancelTest() <-chan CancelTestSignal
	// Shutdown will shut down the hil app.
	Shutdown() <-chan ShutdownSignal
	// RecoverFromFatal will tell the orchestrator to leave the fatal error state and go back to idle.
	RecoverFromFatal() <-chan RecoverFromFatalSignal
	// Status signal is sent on updates to the dispatchers.
	Status() chan<- StatusSignal
	// Results signal is sent at the end of a test execution or on test cancel.
	Results() chan<- ResultsSignal
}
