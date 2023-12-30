package orchestrator

import (
	"context"
	"github.com/ethereum/go-ethereum/event"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"

	"github.com/macformula/hil/flow"
	"go.uber.org/zap"
)

const (
	_loggerName = "orchestrator"
)

type Orchestrator struct {
	l     *zap.Logger
	state State

	sequencer   Sequencer
	dispatchers []Dispatcher
	progSub     event.Subscription

	startSequence chan flow.Sequence
	recoverFatal  chan struct{}
	quit          chan struct{}

	fatalErr *utils.ResettableError
}

func NewOrchestrator(l *zap.Logger, sequencer Sequencer, dispatchers ...Dispatcher) *Orchestrator {
	ret := &Orchestrator{
		l:             l.Named(_loggerName),
		state:         Idle,
		fatalErr:      utils.NewResettaleError(),
		sequencer:     sequencer,
		startSequence: make(chan flow.Sequence),
		recoverFatal:  make(chan struct{}),
		quit:          make(chan struct{}),
	}

	for _, d := range dispatchers {
		ret.dispatchers = append(ret.dispatchers, d)
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
				// This would be a fatal error unrelated to state logic
				o.state = FatalError
				o.fatalErr.Set(err)
			}

			err = o.sequencer.FatalError()
			if err != nil {
				o.fatalErr.Set(err)
				o.state = FatalError
			} else {
				o.state = Idle
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

	o.l.Info("closing orchestrator")

	for _, d := range o.dispatchers {
		err := d.Close()
		if err != nil {
			resettableErr.Set(errors.Wrap(err, "dispatcher close"))
		}
	}

	o.progSub.Unsubscribe()

	return nil
}

func (o *Orchestrator) monitorDispatcher(ctx context.Context, d Dispatcher) {
	for {
		select {
		case <-d.RecoverFromFatal():
			o.l.Info("recover from fatal signal received")

			switch o.state {
			case Idle, Running:
				o.l.Info("orchestrator is not in fatal error state")
			case FatalError:
				o.recoverFatal <- struct{}{}
			}
		case startSignal := <-d.Start():
			o.l.Info("start signal received")

			switch o.state {
			case Idle:
				o.startSequence <- startSignal.Seq
			case Running:
				o.l.Info("test is already running")
			case FatalError:
				o.l.Info("orchestrator is in fatal error state, must recover from fatal error")
			}
		case cancelTestSignal := <-d.CancelTest():
			o.l.Info("cancel test signal received")

			o.quit <- struct{}{}
			return
		case <-ctx.Done():
			o.l.Info("context done signal received")

			o.quit <- struct{}{}
			return
		}
	}
}
