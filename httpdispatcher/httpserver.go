package httpdispatcher

import (
	"context"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
	"go.uber.org/zap"
)

type HttpServer struct {
	l                *zap.Logger
	start            chan orchestrator.StartSignal
	results          chan orchestrator.ResultsSignal
	status           chan orchestrator.StatusSignal
	cancelTest       chan orchestrator.CancelTestSignal
	recoverFromFatal chan orchestrator.RecoverFromFatalSignal
	shutdown         chan orchestrator.ShutdownSignal
	sequences        []flow.Sequence
}

func NewHttpServer(l *zap.Logger) *HttpServer {
	return &HttpServer{
		l:                l.Named(_httpLoggerName),
		start:            make(chan orchestrator.StartSignal),
		results:          make(chan orchestrator.ResultsSignal),
		status:           make(chan orchestrator.StatusSignal),
		cancelTest:       make(chan orchestrator.CancelTestSignal),
		recoverFromFatal: make(chan orchestrator.RecoverFromFatalSignal),
		shutdown:         make(chan orchestrator.ShutdownSignal),
	}
}

func (h *HttpServer) Open(ctx context.Context, sequences []flow.Sequence) error {
	h.sequences = sequences

	err := h.setupServer()
	if err != nil {
		return err
	}

	go h.startServer()

	return nil
}

func (h *HttpServer) Close() error {
	err := h.closeServer()
	if err != nil {
		return err
	}

	return nil
}

func (h *HttpServer) Start() <-chan orchestrator.StartSignal {
	return h.start
}

func (h *HttpServer) CancelTest() <-chan orchestrator.CancelTestSignal {
	return h.cancelTest
}

func (h *HttpServer) Shutdown() <-chan orchestrator.ShutdownSignal {
	return h.shutdown
}

func (h *HttpServer) RecoverFromFatal() <-chan orchestrator.RecoverFromFatalSignal {
	return h.recoverFromFatal
}

func (h *HttpServer) Status() chan<- orchestrator.StatusSignal {
	return h.status
}

func (h *HttpServer) Results() chan<- orchestrator.ResultsSignal {
	return h.results
}

func (h *HttpServer) setupServer() error {
	// TODO: Implement this
	panic("NOT IMPLEMENTED")
	return nil
}

func (h *HttpServer) startServer() {
	// TODO: Implement this
	panic("NOT IMPLEMENTED")
}

func (h *HttpServer) closeServer() error {
	// TODO: Implement this
	panic("NOT IMPLEMENTED")
	return nil
}
