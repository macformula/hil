package httpdispatcher

import (
	"context"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
)

type Server interface {
	Open(context.Context, []flow.Sequence) error
	Close() error
	Start() <-chan orchestrator.StartSignal
	CancelTest() <-chan orchestrator.CancelTestSignal
	Shutdown() <-chan orchestrator.ShutdownSignal
	RecoverFromFatal() <-chan orchestrator.RecoverFromFatalSignal
	Status() chan<- orchestrator.StatusSignal
	Results() chan<- orchestrator.ResultsSignal
}
