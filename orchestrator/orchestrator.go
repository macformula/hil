package orchestrator

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/event"
	"github.com/google/uuid"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_loggerName = "orchestrator"
)

type Orchestrator struct {
	l           *zap.Logger
	state       State
	currentTest TestId

	sequencer   SequencerIface
	dispatchers []DispatcherIface

	progSub  event.Subscription
	progCh   chan flow.Progress
	progress flow.Progress

	testQueue []StartSignal

	shutdownSig     chan ShutdownSignal
	recoverFatalSig chan RecoverFromFatalSignal

	resultFeed event.Feed
	resultSubs []event.Subscription

	statusFeed event.Feed
	statusSubs []event.Subscription

	cancelCurrentTest chan struct{}

	testQueueMtx sync.Mutex
	progressMtx  sync.Mutex

	fatalErr *utils.ResettableError
}

func NewOrchestrator(s SequencerIface, l *zap.Logger, dispatchers ...DispatcherIface) *Orchestrator {
	ret := &Orchestrator{
		l:                 l.Named(_loggerName),
		state:             Idle,
		currentTest:       uuid.Nil,
		sequencer:         s,
		progress:          flow.Progress{},
		testQueue:         make([]StartSignal, 0),
		shutdownSig:       make(chan ShutdownSignal),
		recoverFatalSig:   make(chan RecoverFromFatalSignal),
		progCh:            make(chan flow.Progress),
		cancelCurrentTest: make(chan struct{}),
		testQueueMtx:      sync.Mutex{},
		progressMtx:       sync.Mutex{},
		fatalErr:          utils.NewResettaleError(),
		dispatchers:       dispatchers,
	}

	return ret
}

func (o *Orchestrator) Open(ctx context.Context) error {
	o.l.Info("orchestrator open")
	if len(o.dispatchers) == 0 {
		return errors.Errorf("orchestrator requires at least one dispatcher")
	}

	err := o.sequencer.Open(ctx)
	if err != nil {
		return errors.Wrap(err, "sequencer open")
	}

	// Subscribe to sequencer progress
	o.progSub = o.sequencer.SubscribeToProgress(o.progCh)

	go o.monitorProgress(ctx)

	o.resultSubs = make([]event.Subscription, len(o.dispatchers))
	o.statusSubs = make([]event.Subscription, len(o.dispatchers))

	// Setup dispatchers
	for i, d := range o.dispatchers {
		err = d.Open(ctx)
		if err != nil {
			return errors.Wrap(err, "dispatcher open")
		}

		o.resultSubs[i] = o.resultFeed.Subscribe(d.Results())
		o.statusSubs[i] = o.statusFeed.Subscribe(d.Status())

		// Monitor DispatcherIface signals
		go o.monitorDispatcher(ctx, d)
	}

	return nil
}

func (o *Orchestrator) Run(ctx context.Context) error {
	const _checkForStartSignalPeriod = 100 * time.Millisecond

	for {
		select {
		case <-time.After(_checkForStartSignalPeriod):
		case <-o.shutdownSig:
			return nil
		case <-o.recoverFatalSig:
			o.sequencer.ResetFatalError()
			o.fatalErr.Reset()
			o.state = Idle
			o.statusUpdate()
		}

		if o.state == FatalError {
			continue
		}

		// Attempt to dequeue next test, if no test queued continue.
		startSig, ok := o.dequeueNextTest()
		if !ok {
			if o.state == Running {
				o.state = Idle
				o.statusUpdate()
			}

			continue
		}

		o.currentTest = startSig.TestId

		o.state = Running
		o.statusUpdate()

		isPassing, failedTags, testErrors, err := o.sequencer.Run(ctx, startSig.Seq, o.cancelCurrentTest, o.currentTest)
		if err != nil {
			o.l.Error("sequencer run", zap.Error(errors.Wrap(err, "run")))
		}

		o.l.Info("sending results")

		o.resultFeed.Send(ResultsSignal{
			TestId:     o.currentTest,
			IsPassing:  isPassing,
			FailedTags: failedTags,
			TestErrors: testErrors,
		})

		o.resetProgress()

		err = o.sequencer.FatalError()
		if err != nil {
			o.l.Error("sequencer fatal error", zap.Error(err))

			o.fatalErr.Set(err)
			o.state = FatalError
			o.statusUpdate()
		}
	}
}

func (o *Orchestrator) Close() error {
	var resettableErr = utils.NewResettaleError()

	o.l.Info("closing orchestrator")

	for i, d := range o.dispatchers {
		err := d.Close()
		if err != nil {
			resettableErr.Set(errors.Wrap(err, "dispatcher close"))
		}

		o.resultSubs[i].Unsubscribe()
		o.statusSubs[i].Unsubscribe()
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
				o.cancelCurrentTest <- struct{}{}

				continue
			}

			o.removeTestFromQueue(cancelTestSignal.TestId)

			o.resultFeed.Send(ResultsSignal{
				TestId:     cancelTestSignal.TestId,
				IsPassing:  false,
				FailedTags: make([]flow.Tag, 0),
			})
		case <-d.Shutdown():
			o.l.Info("received shutdown signal",
				zap.String("dispatcher", d.Name()))

			if o.state == Running {
				o.cancelCurrentTest <- struct{}{}
			}

			o.shutdownSig <- ShutdownSignal{}

			return
		case <-ctx.Done():
			o.l.Info("context done signal received")

			return
		}
	}
}

func (o *Orchestrator) monitorProgress(ctx context.Context) {
	for {
		select {
		case progress := <-o.progCh:
			o.progressMtx.Lock()
			o.progress = progress
			o.progressMtx.Unlock()

			o.statusUpdate()
		case <-ctx.Done():
			return
		}
	}
}

func (o *Orchestrator) resetProgress() {
	o.progressMtx.Lock()
	defer o.progressMtx.Unlock()

	o.progress = flow.Progress{}
}

func (o *Orchestrator) addTestToQueue(startSig StartSignal) {
	o.testQueueMtx.Lock()
	defer o.testQueueMtx.Unlock()

	o.testQueue = append(o.testQueue, startSig)

	o.statusUpdate()
}

func (o *Orchestrator) removeTestFromQueue(testId TestId) {
	o.testQueueMtx.Lock()
	defer o.testQueueMtx.Unlock()

	for i := 0; i < len(o.testQueue); i++ {
		if o.testQueue[i].TestId == testId {
			o.testQueue = append(o.testQueue[:i], o.testQueue[i+1:]...)
			return
		}
	}

	o.statusUpdate()
}

func (o *Orchestrator) dequeueNextTest() (StartSignal, bool) {
	o.testQueueMtx.Lock()
	defer o.testQueueMtx.Unlock()

	if len(o.testQueue) == 0 {
		return StartSignal{}, false
	}

	nextTest := o.testQueue[0]
	o.testQueue = o.testQueue[1:]

	o.statusUpdate()

	return nextTest, true
}

func (o *Orchestrator) statusUpdate() {
	o.progressMtx.Lock()
	defer o.progressMtx.Unlock()

	o.statusFeed.Send(StatusSignal{
		OrchestratorState: o.state,
		TestId:            o.currentTest,
		Progress:          o.progress,
		QueueLength:       len(o.testQueue),
		FatalError:        o.fatalErr.Err(),
	})
}
