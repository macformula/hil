package flow

import "time"

// Progress represents the progress of a currently running Sequence.
type Progress struct {
	// CurrentState is the currently running state.
	CurrentState State
	// StateIndex is the index of the currently running state in the sequence.
	StateIndex int
	// Sequence is the currently running sequence.
	Sequence Sequence
	// StatePassed indicates if the state at the given index has passed or failed.
	StatePassed []bool
	// StateDuration indicates the duration for which the state at the given index ran for.
	StateDuration []time.Duration
}
