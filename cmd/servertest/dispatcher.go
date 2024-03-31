package main

import (
	"context"
	"github.com/macformula/hil/flow"

	"github.com/macformula/hil/orchestrator"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_dispatcherLoggerName = "serverDispatcher"
	_httpLoggerName       = "httpDispatcher"
)

// ServerDispatcher is the server implementation of the orchestrator.DispatcherIface.
type ServerDispatcher struct {
	l                *zap.Logger
	start            chan orchestrator.StartSignal
	results          chan orchestrator.ResultsSignal
	status           chan orchestrator.StatusSignal
	cancelTest       chan orchestrator.CancelTestSignal
	recoverFromFatal chan orchestrator.RecoverFromFatalSignal
	shutdown         chan orchestrator.ShutdownSignal
	sequences        []flow.Sequence
	server           Server
}

// NewServerDispatcher creates a cli dispatcher object.
func NewServerDispatcher(sequences []flow.Sequence, server Server, l *zap.Logger) *ServerDispatcher {
	return &ServerDispatcher{
		l:                l.Named(_dispatcherLoggerName),
		start:            make(chan orchestrator.StartSignal, 5),
		results:          make(chan orchestrator.ResultsSignal),
		status:           make(chan orchestrator.StatusSignal),
		cancelTest:       make(chan orchestrator.CancelTestSignal),
		recoverFromFatal: make(chan orchestrator.RecoverFromFatalSignal),
		shutdown:         make(chan orchestrator.ShutdownSignal),
		sequences:        sequences,
		server:           server,
	}
}

// Shutdown will shut down the hil app.
func (s *ServerDispatcher) Shutdown() <-chan orchestrator.ShutdownSignal {
	return s.shutdown
}

// Close should close all objects held by the dispatcher.
func (s *ServerDispatcher) Close() error {
	err := s.server.Close()
	return err
}

// Open should set up all initial calls for the dispatcher.
func (s *ServerDispatcher) Open(ctx context.Context) error {
	err := s.server.Open(ctx, s.sequences)

	if err != nil {
		return errors.Wrap(err, "cli open")
	}

	go s.monitorServer(ctx)
	go s.monitorOrchestrator(ctx)

	return nil
}

// Start signal is sent by the dispatcher to the orchestrator to start a test sequence.
func (s *ServerDispatcher) Start() <-chan orchestrator.StartSignal {
	return s.start
}

// CancelTest will cancel execution of the test with the given ID.
func (s *ServerDispatcher) CancelTest() <-chan orchestrator.CancelTestSignal {
	return s.cancelTest
}

// RecoverFromFatal will tell the orchestrator to leave the fatal error state and go back to idle.
func (s *ServerDispatcher) RecoverFromFatal() <-chan orchestrator.RecoverFromFatalSignal {
	return s.recoverFromFatal
}

// Status signal is sent on updates from the orchestrator.
func (s *ServerDispatcher) Status() chan<- orchestrator.StatusSignal {
	return s.status
}

// Results signal is sent at the end of a test execution or on test cancel.
func (s *ServerDispatcher) Results() chan<- orchestrator.ResultsSignal {
	return s.results
}

// Quit signal will shut down the app.
func (s *ServerDispatcher) Quit() chan orchestrator.ShutdownSignal {
	return s.shutdown
}

// Name returns the dispatcher name.
func (s *ServerDispatcher) Name() string {
	return "http_dispatcher"
}

func (s *ServerDispatcher) monitorServer(ctx context.Context) {
	for {
		select {
		case recoverSignal := <-s.server.RecoverFromFatal():
			s.l.Info("recover from fatal signal received")

			s.recoverFromFatal <- recoverSignal
		case startSignal := <-s.server.Start():
			s.l.Info("start signal received")

			s.start <- startSignal
		case cancelSignal := <-s.server.CancelTest():
			s.l.Info("cancel test signal received")

			s.cancelTest <- cancelSignal
		case fatalSignal := <-s.server.RecoverFromFatal():
			s.l.Info("fatal recovery signal received")

			s.recoverFromFatal <- fatalSignal
		case <-ctx.Done():
			s.l.Info("context done signal received")

			return
		}
	}
}

func (s *ServerDispatcher) monitorOrchestrator(ctx context.Context) {
	for {
		select {
		case status := <-s.status:
			s.l.Info("status signal received")

			s.server.Status() <- status
		case results := <-s.results:
			s.l.Info("results signal received")

			s.server.Results() <- results
		case <-ctx.Done():
			s.l.Info("context done signal received")

			return
		}
	}
}
