package test

import (
	"context"
	"github.com/macformula/hil/orchestrator"
)

type simpleDispatcher struct {
	startSig        chan orchestrator.StartSignal
	shutdownSig     chan orchestrator.ShutdownSignal
	cancelSig       chan orchestrator.CancelTestSignal
	recoverFatalSig chan orchestrator.RecoverFromFatalSignal
	status          chan orchestrator.StatusSignal
	resultsSig      chan orchestrator.ResultsSignal
}

func (s simpleDispatcher) Close() error {
	return nil
}

func (s simpleDispatcher) Name() string {
	return "simple_dispatcher"
}

func (s simpleDispatcher) Open(ctx context.Context) error {
	return nil
}

func (s simpleDispatcher) Start() <-chan orchestrator.StartSignal {
	return s.startSig
}

func (s simpleDispatcher) CancelTest() <-chan orchestrator.CancelTestSignal {
	return s.cancelSig
}

func (s simpleDispatcher) Shutdown() <-chan orchestrator.ShutdownSignal {
	return s.shutdownSig
}

func (s simpleDispatcher) RecoverFromFatal() <-chan orchestrator.RecoverFromFatalSignal {
	return s.recoverFatalSig
}

func (s simpleDispatcher) Status() chan<- orchestrator.StatusSignal {
	return s.status
}

func (s simpleDispatcher) Results() chan<- orchestrator.ResultsSignal {
	return s.resultsSig
}
