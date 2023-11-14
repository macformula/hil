package orchestrator

import (
	"context"
	"github.com/macformula/hil/flow"
	"go.uber.org/zap"
)

type Orchestrator struct {
	l *zap.Logger
	state State

	s *flow.Sequencer
	d []Dispatcher

	testStart chan struct{}
	testComplete chan struct{}
}

type Option = func(*Orchestrator)

func WithDispatcher(d Dispatcher) Option {
	return func(o *Orchestrator) {
		o.d = append(o.d, d)
	}
}

func NewOrchestrator(l *zap.Logger, opts ...Option) *Orchestrator {
	ret := &Orchestrator{
		l: l,
	}

	for _, o := range opts {
		o(ret)
	}

	return ret
}

func (o *Orchestrator) Open(ctx context.Context) error {

	for
	go o.waitForStart(ctx)

	return nil
}

func (o *Orchestrator) Close(ctx context.Context) error {
	return nil
}

func (o *Orchestrator) monitorDispatcher(ctx context.Context, d Dispatcher) {

	for {
		switch o.state {
		case Idle:
			select {
			case <-ctx.Done():
			case <-d.Start():
				o.testStart <- struct {}{}
			case <-o.testStart:
			}
		}

	}

}
