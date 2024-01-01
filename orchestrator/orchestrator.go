package orchestrator

import (
	"context"
	"github.com/ethereum/go-ethereum/event"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"
	"sync"
	"time"

	"github.com/macformula/hil/flow"
	"go.uber.org/zap"
)

const (
	_loggerName = "orchestrator"
	_maxTestQueueLen
)

type Orchestrator struct {
	l           *zap.Logger
	state       State
	currentTest TestId

	sequencer   SequencerIface
	dispatchers []DispatcherIface
	progSub     event.Subscription
	progCh      chan flow.Progress

	testQueue []StartSignal

	startSignal     chan StartSignal
	resultsSig      chan ResultsSignal
	statusSig       chan StatusSignal
	recoverFatalSig chan RecoverFromFatalSignal

	cancelCurrentTest context.CancelFunc

	status  StatusSignal
	results ResultsSignal

	testQueueMtx sync.Mutex

	fatalErr *utils.ResettableError
}

func NewOrchestrator(s SequencerIface, l *zap.Logger, dispatchers ...DispatcherIface) *Orchestrator {
	ret := &Orchestrator{
		l:               l.Named(_loggerName),
		fatalErr:        utils.NewResettaleError(),
		sequencer:       s,
		testQueue:       make([]StartSignal, 0),
		startSignal:     make(chan StartSignal),
		resultsSig:      make(chan ResultsSignal),
		statusSig:       make(chan StatusSignal),
		recoverFatalSig: make(chan RecoverFromFatalSignal),
		progCh:          make(chan flow.Progress),
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

	// Subscribe to sequencer progress
	o.progSub = o.sequencer.SubscribeToProgress(o.progCh)

	// Setup dispatchers
	for _, d := range o.dispatchers {
		err := d.Open(ctx)
		if err != nil {
			return errors.Wrap(err, "dispatcher open")
		}

		// Monitor DispatcherIface signals
		go o.monitorDispatcher(ctx, d)
	}

	return nil
}

func (o *Orchestrator) Run(ctx context.Context) error {
	const _checkForStartSignalPeriod = 20 * time.Millisecond

	for {
		select {
		case <-time.After(_checkForStartSignalPeriod):
		case <-o.recoverFatalSig:
			o.sequencer.ResetFatalError()
			o.state = Idle
		}

		startSig, ok := o.dequeueNextTest()
		if !ok {
			continue
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

func (o *Orchestrator) monitorDispatcher(ctx context.Context, d DispatcherIface) {
	for {
		select {
		case recoverFatalSig := <-d.RecoverFromFatal():
			o.l.Info("recover from fatal signal received", zap.String("dispatcher", d.Name()))

			switch o.state {
			case Idle, Running, Unknown:
				o.l.Warn("commanded recover from fatal when orchestrator is not in fatal error state",
					zap.String("state", o.state.String()),
					zap.String("dispatcher", d.Name()))
			case FatalError:
				o.recoverFatalSig <- recoverFatalSig
			}
		case startSig := <-d.Start():
			o.l.Info("start signal received",
				zap.String("dispatcher", d.Name()),
				zap.String("test id", startSig.TestId.String()))

			switch o.state {
			case Idle, Running:
				o.addTestToQueue(startSig)
			case FatalError:
				o.l.Warn("orchestrator is in fatal error state, must recover from fatal error",
					zap.String("dispatcher", d.Name()))
			case Unknown:
				o.l.Warn("orchestrator is in unknown state, this should not happen")
			}
		case cancelTestSignal := <-d.CancelTest():
			o.l.Info("cancel test signal received",
				zap.String("dispatcher", d.Name()),
				zap.String("test id", cancelTestSignal.TestId.String()))

			if cancelTestSignal.TestId == o.currentTest {
				o.cancelCurrentTest()

				continue
			}

			go o.removeTestFromQueue(cancelTestSignal.TestId)
		case <-ctx.Done():
			o.l.Info("context done signal received")

			return
		}
	}
}

func (o *Orchestrator) addTestToQueue(startSig StartSignal) {
	o.testQueueMtx.Lock()
	defer o.testQueueMtx.Unlock()
	o.testQueue = append(o.testQueue, startSig)
}

func (o *Orchestrator) removeTestFromQueue(testId TestId) {
	for i := 0; i < len(o.testQueue); i++ {
		if o.testQueue[i].TestId == testId {
			o.testQueueMtx.Lock()
			o.testQueue = append(o.testQueue[:i], o.testQueue[i+1:]...)
			o.testQueueMtx.Unlock()
			return
		}
	}

}

func (o *Orchestrator) dequeueNextTest() (StartSignal, bool) {
	o.testQueueMtx.Lock()
	defer o.testQueueMtx.Unlock()

}
