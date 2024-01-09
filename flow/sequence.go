package flow

// Sequence is a list of states that will be run in the order provided.
type Sequence struct {
	Name   string
	Desc   string
	States []State
}
