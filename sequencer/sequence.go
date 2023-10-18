package sequencer

// Sequence is a list of states that will be run in the order provided
type Sequence []State

type Progress struct {
	CurrentState State
	StateIndex   int
	TotalStates  int
}
