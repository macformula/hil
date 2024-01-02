package dispatcher

import (
	"context"
	"github.com/google/uuid"
	"github.com/macformula/hil/orchestrator"
	"go.uber.org/zap"
)

const (
	_loggerName = "cliDispatcher"
)

type CliDispatcher struct {
	l                *zap.Logger
	start            chan orchestrator.StartSignal
	results          chan orchestrator.ResultsSignal
	status           chan orchestrator.StatusSignal
	cancelTest       chan orchestrator.CancelTestSignal
	recoverFromFatal chan orchestrator.RecoverFromFatalSignal
	testToRun        uuid.UUID
	stop             chan struct{}
	cli              cliInterface
}

func NewCliDispatcher(l *zap.Logger) *CliDispatcher {
	return &CliDispatcher{
		l:                l.Named(_loggerName),
		start:            make(chan orchestrator.StartSignal, 5),
		results:          make(chan orchestrator.ResultsSignal),
		status:           make(chan orchestrator.StatusSignal),
		cancelTest:       make(chan orchestrator.CancelTestSignal),
		recoverFromFatal: make(chan orchestrator.RecoverFromFatalSignal),
		cli:              newCli(zap.L()),
	}
}

func (c *CliDispatcher) Shutdown() <-chan orchestrator.ShutdownSignal {
	return c.cli.Quit()
}

func (c *CliDispatcher) Close() error {
	err := c.cli.Close()
	return err
}

func (c *CliDispatcher) Open(ctx context.Context) error {
	err := c.cli.Open(ctx)
	if err != nil {
		return err
	}

	c.cli.Start()

	go c.monitorCli(ctx, c.cli)
	go c.monitorOrchestrator(ctx)

	return nil
}

func (c *CliDispatcher) Start() <-chan orchestrator.StartSignal {
	return c.start
}

func (c *CliDispatcher) CancelTest() <-chan orchestrator.CancelTestSignal {
	return c.cancelTest
}

func (c *CliDispatcher) RecoverFromFatal() <-chan orchestrator.RecoverFromFatalSignal {
	return c.recoverFromFatal
}

func (c *CliDispatcher) Status() chan<- orchestrator.StatusSignal {
	return c.status
}

func (c *CliDispatcher) Results() chan<- orchestrator.ResultsSignal {
	return c.results
}

func (c *CliDispatcher) Quit() chan orchestrator.ShutdownSignal {
	return c.cli.Quit()
}

func (c *CliDispatcher) Name() string {
	return "CliDispatcher"
}

func (c *CliDispatcher) monitorCli(ctx context.Context, cli cliInterface) {
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

			c.Close()
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

			c.Close()
			return
		}
	}
}
