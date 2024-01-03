package orchestrator

//go:generate enumer -type=State state.go
type State int

const (
	Unknown State = iota
	Idle
	Running
	FatalError
)
