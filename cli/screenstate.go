package cli

//go:generate enumer -type=screenState screenstate.go
type screenState int

const (
	Unknown screenState = iota
	Idle
	Running
	FatalError
	Results
)
