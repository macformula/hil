package flow

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/event"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_loggerName = "sequencer"
)

type Sequencer struct {
	l            *zap.Logger
	sequence     Sequence
	progress     Progress
	progressFeed event.Feed

	complete chan struct{}
	stop     chan struct{}

	fatalErr *utils.ResettableError
}

func NewSequencer(l *zap.Logger) *Sequencer {
	return &Sequencer{
		l:            l.Named(_loggerName),
		progressFeed: event.Feed{},
		fatalErr:     utils.NewResettaleError(),
	}
}

func (s *Sequencer) SubscribeToProgress(progCh chan Progress) event.Subscription {
	return s.progressFeed.Subscribe(progCh)
}

func (s *Sequencer) Run(ctx context.Context, seq Sequence) error {
	if len(seq) == 0 {
		return errors.New("sequence cannot be empty")
	}

	s.sequence = seq
	s.progress = Progress{
		CurrentState:  nil,
		StateDuration: make([]time.Duration, len(seq)),
		StateIndex:    0,
		Complete:      false,
		Sequence:      seq,
		TotalStates:   len(seq),
	}

	s.stop = make(chan struct{})
	s.complete = make(chan struct{})

	err := s.runSequence(ctx)
	if err != nil {
		return errors.Wrap(err, "run sequence")
	}

	return nil
}

func (s *Sequencer) runSequence(ctx context.Context) error {
	var (
		startTime  = time.Time{}
		nSubs      = 0
		onceErr    = utils.NewResettaleError()
		timeoutCtx context.Context
		cancel     context.CancelFunc
	)

	for idx, state := range s.sequence {
		// Check stop before setting current state
		s.progress.CurrentState = state
		s.progress.StateIndex = idx

		nSubs = s.progressFeed.Send(s.progress)
		s.l.Debug("published progress", zap.Int("num_subs", nSubs))

		s.l.Info("starting next state", zap.String("state", state.Name()))

		startTime = time.Now()

		timeoutCtx, cancel = context.WithTimeout(ctx, state.Timeout())

		err := state.Setup(timeoutCtx)
		if err != nil {
			s.fatalErr.Set(errors.Wrapf(err, "setup (%s)", state.Name()))

			// cancel timeoutCtx, do not use defer here because of for loop
			cancel()

			break
		}

		// cancel timeoutCtx, do not use defer here because of for loop
		cancel()

		timeoutCtx, cancel = context.WithTimeout(ctx, state.Timeout())

		err = state.Run(timeoutCtx)
		if err != nil {
			onceErr.Set(errors.Wrapf(err, "run (%s)", state.Name()))
		}

		// cancel timeoutCtx, do not use defer here because of for loop
		cancel()

		s.progress.StateDuration[idx] = time.Since(startTime)

		err = state.FatalError()
		if err != nil {
			s.l.Info("encountered fatal error", zap.String("state", state.Name()), zap.Error(err))
			s.fatalErr.Set(errors.Wrapf(err, "fatal error (%s)", state.Name()))

			break
		}
	}

	// Complete
	s.l.Info("sequence complete")

	// Send final update to signal that the sequence has completed
	s.progress.Complete = true
	s.progress.StateIndex++
	_ = s.progressFeed.Send(s.progress)

	return onceErr.Err()
}

func (s *Sequencer) FatalError() error {
	return s.fatalErr.Err()
}
