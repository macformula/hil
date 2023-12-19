package orchestrator

type State int

const (
	Unknown State = iota
	Idle
	Running
	FatalError
)
