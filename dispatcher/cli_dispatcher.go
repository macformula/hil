package dispatcher

import (
	"context"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
)

type cliDispatcher struct {
	start             chan flow.Sequence
	quit              chan struct{}
	recoverFromFatal  chan struct{}
	progress          chan flow.Progress
	orchestratorState orchestrator.State
	stop              chan struct{}
}

func newCliDispatcher() *cliDispatcher {
	return &cliDispatcher{
		start:    make(chan flow.Sequence),
		quit:     make(chan struct{}),
		progress: make(chan flow.Progress),
	}
}

func (c *cliDispatcher) Close() error {
	//TODO implement me
	panic("implement me")
}

func (c *cliDispatcher) Open(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (c *cliDispatcher) Start() <-chan orchestrator.StartSignal {
	//TODO implement me
	panic("implement me")
}

func (c *cliDispatcher) CancelTest() <-chan orchestrator.TestId {
	//TODO implement me
	panic("implement me")
}

func (c *cliDispatcher) RecoverFromFatal() <-chan struct{} {
	//TODO implement me
	panic("implement me")
}

func (c *cliDispatcher) Status() chan<- orchestrator.StatusSignal {
	//TODO implement me
	panic("implement me")
}

func (c *cliDispatcher) Results() chan<- orchestrator.ResultSignal {
	//TODO implement me
	panic("implement me")
}
