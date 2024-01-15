package flow

// Sequence is a list of states that will be run in the order provided.
type Sequence struct {
	// Name of the sequence.
	Name string
	// Desc is the description for the purpose of the sequence.
	Desc string
	// States are the states to be run in the order provided.
	States []State
}
