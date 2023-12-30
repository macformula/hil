package dispatcher

import (
	"context"
	"github.com/google/uuid"
	"github.com/macformula/hil/orchestrator"
)

type CliDispatcher struct {
	start            chan orchestrator.StartSignal
	results          chan orchestrator.ResultSignal
	status           chan orchestrator.StatusSignal
	cancelTest       chan orchestrator.TestId
	recoverFromFatal chan struct{}
	testsToRun       uuid.UUID
	stop             chan struct{}
	cli              cliInterface
}

func NewCliDispatcher() *CliDispatcher {
	return &CliDispatcher{
		start:      make(chan orchestrator.StartSignal, 5),
		results:    make(chan orchestrator.ResultSignal),
		status:     make(chan orchestrator.StatusSignal),
		cancelTest: make(chan orchestrator.TestId),
		//testsToRun: make([]uuid.UUID, 1), // keep it to 1
		cli: newCli(),
	}
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

	return nil
}

func (c *CliDispatcher) Start() <-chan orchestrator.StartSignal {
	return c.start
}

func (c *CliDispatcher) CancelTest() <-chan orchestrator.TestId {
	//TODO implement me
	panic("implement me")
}

func (c *CliDispatcher) RecoverFromFatal() <-chan struct{} {
	//TODO implement me
	panic("implement me")
}

func (c *CliDispatcher) Status() chan<- orchestrator.StatusSignal {
	return c.status
}

func (c *CliDispatcher) Results() chan<- orchestrator.ResultSignal {
	return c.results
}
