package cli

import (
	"context"
	"github.com/macformula/hil/flow"

	"github.com/macformula/hil/orchestrator"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_loggerName = "cliDispatcher"
)

// CliDispatcher is the cli implementation of the orchestrator.DispatcherIface.
type CliDispatcher struct {
	l                *zap.Logger
	start            chan orchestrator.StartSignal
	results          chan orchestrator.ResultsSignal
	status           chan orchestrator.StatusSignal
	cancelTest       chan orchestrator.CancelTestSignal
	recoverFromFatal chan orchestrator.RecoverFromFatalSignal
	cli              cliIface
}

// NewCliDispatcher creates a cli dispatcher object.
func NewCliDispatcher(sequences []flow.Sequence, l *zap.Logger) *CliDispatcher {
	return &CliDispatcher{
		l:                l.Named(_loggerName),
		start:            make(chan orchestrator.StartSignal, 5),
		results:          make(chan orchestrator.ResultsSignal),
		status:           make(chan orchestrator.StatusSignal),
		cancelTest:       make(chan orchestrator.CancelTestSignal),
		recoverFromFatal: make(chan orchestrator.RecoverFromFatalSignal),
		cli:              newCliModel(sequences, l),
	}
}

// Shutdown will shut down the hil app.
func (c *CliDispatcher) Shutdown() <-chan orchestrator.ShutdownSignal {
	return c.cli.Quit()
}

// Close should close all objects held by the dispatcher.
func (c *CliDispatcher) Close() error {
	err := c.cli.Close()
	return err
}

// Open should set up all initial calls for the dispatcher.
func (c *CliDispatcher) Open(ctx context.Context) error {
	err := c.cli.Open(ctx)

	if err != nil {
		return errors.Wrap(err, "cli open")
	}

	go c.monitorCli(ctx, c.cli)
	go c.monitorOrchestrator(ctx)

	return nil
}

// Start signal is sent by the dispatcher to the orchestrator to start a test sequence.
func (c *CliDispatcher) Start() <-chan orchestrator.StartSignal {
	return c.start
}

// CancelTest will cancel execution of the test with the given ID.
func (c *CliDispatcher) CancelTest() <-chan orchestrator.CancelTestSignal {
	return c.cancelTest
}

// RecoverFromFatal will tell the orchestrator to leave the fatal error state and go back to idle.
func (c *CliDispatcher) RecoverFromFatal() <-chan orchestrator.RecoverFromFatalSignal {
	return c.recoverFromFatal
}

// Status signal is sent on updates from the orchestrator.
func (c *CliDispatcher) Status() chan<- orchestrator.StatusSignal {
	return c.status
}

// Results signal is sent at the end of a test execution or on test cancel.
func (c *CliDispatcher) Results() chan<- orchestrator.ResultsSignal {
	return c.results
}

// Quit signal will shut down the app.
func (c *CliDispatcher) Quit() chan orchestrator.ShutdownSignal {
	return c.cli.Quit()
}

// Name returns the dispatcher name.
func (c *CliDispatcher) Name() string {
	return "cli_dispatcher"
}

func (c *CliDispatcher) monitorCli(ctx context.Context, cli cliIface) {
	for {
		select {
		case recoverSignal := <-cli.RecoverFromFatal():
			c.l.Info("recover from fatal signal received")

			c.recoverFromFatal <- recoverSignal
		case startSignal := <-cli.Start():
			c.l.Info("start signal received")

			c.start <- startSignal
		case cancelSignal := <-cli.CancelTest():
			c.l.Info("cancel test signal received")

			c.cancelTest <- cancelSignal
		case fatalSignal := <-cli.RecoverFromFatal():
			c.l.Info("fatal recovery signal received")

			c.recoverFromFatal <- fatalSignal
		case <-ctx.Done():
			c.l.Info("context done signal received")

			return
		}
	}
}

func (c *CliDispatcher) monitorOrchestrator(ctx context.Context) {
	for {
		select {
		case status := <-c.status:
			c.l.Info("status signal received")

			c.cli.Status() <- status
			c.l.Info("after status sent to cli")
		case results := <-c.results:
			c.l.Info("results signal received")

			c.cli.Results() <- results
		case <-ctx.Done():
			c.l.Info("context done signal received")

			return
		}
	}
}
