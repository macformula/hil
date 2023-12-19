package orchestrator

import (
	"context"
	"github.com/ethereum/go-ethereum/event"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"

	"github.com/macformula/hil/flow"
	"go.uber.org/zap"
)

type Orchestrator struct {
	l     *zap.Logger
	state State

	sequencer   *flow.Sequencer
	dispatchers []Dispatcher
	progSubs    []event.Subscription

	startSequence chan flow.Sequence
	recoverFatal  chan struct{}
	quit          chan struct{}

	fatalErr *utils.ResettableError
}

type Option = func(*Orchestrator)

func WithDispatcher(d Dispatcher) Option {
	return func(o *Orchestrator) {
		o.dispatchers = append(o.dispatchers, d)
	}
}

func NewOrchestrator(l *zap.Logger, opts ...Option) *Orchestrator {
	ret := &Orchestrator{
		l:        l,
		state:    Idle,
		fatalErr: utils.NewResettaleError(),
	}

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}

func (o *Orchestrator) Open(ctx context.Context) error {
	if len(o.dispatchers) == 0 {
		return errors.Errorf("orchestrator requires at least one dispatcher")
	}

	// Setup dispatchers
	for _, d := range o.dispatchers {
		err := d.Open(ctx)
		if err != nil {
			return errors.Wrap(err, "dispatcher open")
		}

		// Subscribe all dispatchers to sequencer progress
		o.progSubs = append(o.progSubs, o.sequencer.SubscribeToProgress(d.Progress()))

		// Monitor Dispatcher signals
		go o.monitorDispatcher(ctx, d)
	}

	return nil
}

func (o *Orchestrator) Run(ctx context.Context) error {
	for {
		select {
		case seq := <-o.startSequence:
			o.state = Running

			err := o.sequencer.Run(ctx, seq)
			if err != nil {
				// TODO: figure out how to manage errors
			}

			err = o.sequencer.FatalError()
			if err != nil {
				o.fatalErr.Set(err)
			}
		case <-o.recoverFatal:
			o.state = Idle

			o.sequencer.ResetFatalError()
		case <-o.quit:
			return nil
		}
	}
}

func (o *Orchestrator) Close() error {
	var resettableErr = utils.NewResettaleError()
	// Setup dispatchers
	for _, d := range o.dispatchers {
		err := d.Close()
		if err != nil {
			resettableErr.Set(errors.Wrap(err, "dispatcher close"))
		}
	}

	for _, sub := range o.progSubs {
		sub.Unsubscribe()
	}

	return nil
}

func (o *Orchestrator) monitorDispatcher(ctx context.Context, d Dispatcher) {
	for {
		select {
		case <-d.RecoverFromFatal():
			o.l.Info("recover from fatal signal received")

			switch o.state {
			case Idle, Running:
				o.l.Info("test is not in fatal error state")
			case FatalError:
				o.recoverFatal <- struct{}{}
			}
		case seq := <-d.Start():
			o.l.Info("start signal received")

			switch o.state {
			case Idle:
				o.startSequence <- seq
			case Running:
				o.l.Info("test is already running")
			case FatalError:
				o.l.Info("orchestrator is in fatal error state, must recover from fatal error")
			}
		case <-d.Quit():
			o.l.Info("quit signal received")

			o.quit <- struct{}{}
			return
		case <-ctx.Done():
			o.l.Info("context done signal received")

			o.quit <- struct{}{}
			return
		}
	}
}
