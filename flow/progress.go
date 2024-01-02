package flow

import "time"

// Sequence is a list of states that will be run in the order provided.
type Sequence []State

// Progress represents the progress of a currently running Sequence.
type Progress struct {
	CurrentState  State
	StateIndex    int
	Sequence      Sequence
	StatePassed   []bool
	StateDuration []time.Duration
}
