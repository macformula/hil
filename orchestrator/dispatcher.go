package orchestrator

import (
	"context"
	"io"
)

type Dispatcher interface {
	io.Closer
	Open(context.Context) error
	Start() <-chan StartSignal
	CancelTest() <-chan CancelTestSignal
	RecoverFromFatal() <-chan RecoverFromFatalSignal
	Status() chan<- StatusSignal
	Results() chan<- ResultsSignal
}
