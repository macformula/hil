package orchestrator

import (
	"github.com/google/uuid"
	"github.com/macformula/hil/flow"
)

type TestId = uuid.UUID

type StatusSignal struct {
	// State is the current state of the Orchestrator
	OrchestratorState State
	// TestId is the current running test. It is only valid if State is Running
	TestId TestId
	// Progress is the current Progress of the flow.Sequencer.
	Progress flow.Progress
	// QueueLength is the current length of the test queue.
	QueueLength int
	// FatalError is the current fatal error. It is only valid if state is FatalError
	FatalError error
}

// StartSignal is the signal that is sent to the orchestrator from the dispatcher to start a test
type StartSignal struct {
	TestId   TestId
	Seq      flow.Sequence
	Metadata map[string]string
}

type ResultSignal struct {
	TestId    TestId
	IsPassing bool
	// Should be of type Tag, will replace this later
	FailedTags []string
}
