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

// Sequencer is responsible for managing the setup and execution of a Sequence.
type Sequencer struct {
	l            *zap.Logger
	progress     Progress
	progressFeed event.Feed

	fatalErr *utils.ResettableError
}

// NewSequencer returns a Sequencer object reference.
func NewSequencer(l *zap.Logger) *Sequencer {
	return &Sequencer{
		l:            l.Named(_loggerName),
		progressFeed: event.Feed{},
		fatalErr:     utils.NewResettaleError(),
	}
}

// SubscribeToProgress subscribes to the progress of the Sequencer accross its Sequence runs.
// The Progress channel gets updated whenever there is new information available.
func (s *Sequencer) SubscribeToProgress(progCh chan Progress) event.Subscription {
	return s.progressFeed.Subscribe(progCh)
}

// Run will run the sequence provided. FatalError must be called after Run to check for any non-recoverable errors.
func (s *Sequencer) Run(ctx context.Context, seq Sequence) error {
	if len(seq) == 0 {
		return errors.New("sequence cannot be empty")
	}

	s.progress = Progress{
		CurrentState:  nil,
		StateDuration: make([]time.Duration, len(seq)),
		StateIndex:    0,
		Complete:      false,
		Sequence:      seq,
		TotalStates:   len(seq),
	}

	err := s.runSequence(ctx, seq)
	if err != nil {
		return errors.Wrap(err, "run sequence")
	}

	return nil
}

func (s *Sequencer) runSequence(ctx context.Context, seq Sequence) error {
	var (
		startTime  = time.Time{}
		nSubs      = 0
		onceErr    = utils.NewResettaleError()
		timeoutCtx context.Context
		cancel     context.CancelFunc
	)
	
	for idx, state := range seq {
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
	_ = s.progressFeed.Send(s.progress)

	return onceErr.Err()
}

// FatalError indicates that there is an error that requires intervention.
// The Sequencer will stop executing all remaining states in the Sequence if it encounters a fatal error.
func (s *Sequencer) FatalError() error {
	return s.fatalErr.Err()
}
