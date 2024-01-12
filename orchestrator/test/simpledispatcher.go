package test

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/macformula/hil/dispatcher/test"
	"github.com/macformula/hil/orchestrator"
	"go.uber.org/zap"
)

const (
	_loggerName = "simpledispatcher"
)

type SimpleDispatcher struct {
	startSig        chan orchestrator.StartSignal
	shutdownSig     chan orchestrator.ShutdownSignal
	cancelSig       chan orchestrator.CancelTestSignal
	recoverFatalSig chan orchestrator.RecoverFromFatalSignal
	status          chan orchestrator.StatusSignal
	resultsSig      chan orchestrator.ResultsSignal
	durations       []time.Duration
	l               *zap.Logger
}

func NewSimpleDispatcher(l *zap.Logger, durations ...time.Duration) *SimpleDispatcher {
	return &SimpleDispatcher{
		l:               l.Named(_loggerName),
		startSig:        make(chan orchestrator.StartSignal),
		shutdownSig:     make(chan orchestrator.ShutdownSignal),
		cancelSig:       make(chan orchestrator.CancelTestSignal),
		recoverFatalSig: make(chan orchestrator.RecoverFromFatalSignal),
		status:          make(chan orchestrator.StatusSignal),
		resultsSig:      make(chan orchestrator.ResultsSignal),
		durations:       durations,
	}
}

func (s *SimpleDispatcher) Close() error {
	return nil
}

func (s *SimpleDispatcher) Name() string {
	return "simple_dispatcher"
}

func (s *SimpleDispatcher) Open(ctx context.Context) error {
	go s.simulate(s.durations)
	go s.monitorOrchestrator(ctx)
	return nil
}

func (s *SimpleDispatcher) simulate(durations []time.Duration) {
	for _, duration := range durations {
		time.Sleep(duration)
		testId := uuid.New()
		s.startSig <- orchestrator.StartSignal{
			TestId:   testId,
			Seq:      test.SleepSequence,
			Metadata: nil,
		}
		for {
			results := <-s.resultsSig
			if results.TestId == testId {
				break
			}
		}
	}
}

func (s *SimpleDispatcher) Start() <-chan orchestrator.StartSignal {
	return s.startSig
}

func (s *SimpleDispatcher) CancelTest() <-chan orchestrator.CancelTestSignal {
	return s.cancelSig
}

func (s *SimpleDispatcher) Shutdown() <-chan orchestrator.ShutdownSignal {
	return s.shutdownSig
}

func (s *SimpleDispatcher) RecoverFromFatal() <-chan orchestrator.RecoverFromFatalSignal {
	return s.recoverFatalSig
}

func (s *SimpleDispatcher) Status() chan<- orchestrator.StatusSignal {
	return s.status
}

func (s *SimpleDispatcher) Results() chan<- orchestrator.ResultsSignal {
	return s.resultsSig
}

func (s *SimpleDispatcher) monitorOrchestrator(ctx context.Context) {
	for {
		select {
		case status := <-s.status:
			s.l.Info("status signal received",
				zap.String("testid", status.TestId.String()),
				zap.String("orchestrator state", status.OrchestratorState.String()))

		case results := <-s.resultsSig:
			s.l.Info("results signal received",
				zap.String("testid", results.TestId.String()))

		case <-ctx.Done():
			s.l.Info("context done signal received")

			s.Close()
			return
		}
	}
}
