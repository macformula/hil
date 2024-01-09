package flow

import "time"

// Progress represents the progress of a currently running Sequence.
type Progress struct {
	CurrentState  State
	StateIndex    int
	Sequence      Sequence
	StatePassed   []bool
	StateDuration []time.Duration
}
