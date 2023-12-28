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

// SubscribeToProgress subscribes to the progress of the Sequencer across its Sequence runs.
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
	}

	err := s.runSequence(ctx, seq)
	if err != nil {
		return errors.Wrap(err, "run sequence")
	}

	return nil
}

// FatalError indicates that there is an error that requires intervention.
// The Sequencer will stop executing all remaining states in the Sequence if it encounters a fatal error.
func (s *Sequencer) FatalError() error {
	return s.fatalErr.Err()
}

// ResetFatalError sets the fatal error to nil.
func (s *Sequencer) ResetFatalError() {
	s.fatalErr.Reset()
}

func (s *Sequencer) runSequence(ctx context.Context, seq Sequence) error {
	var (
		onceErr = utils.NewResettaleError()
		nSubs   = 0
	)

	for idx, state := range seq {
		s.progress.CurrentState = state
		s.progress.StateIndex = idx

		nSubs = s.progressFeed.Send(s.progress)
		s.l.Debug("published progress", zap.Int("num_subs", nSubs))

		s.l.Info("starting next state", zap.String("state", state.Name()))

		regularErr := s.runState(ctx, state)

		err := s.processResults(ctx, state, regularErr)
		if err != nil {
			return errors.Wrap(err, "process results")
		}
		//TODO: figure out how to handle this error

	}

	// Complete
	s.l.Info("sequence complete")

	// Send final update to signal that the sequence has completed
	s.progress.Complete = true
	_ = s.progressFeed.Send(s.progress)

	return onceErr.Err()
}

func (s *Sequencer) runState(ctx context.Context, state State) error {
	var (
		timeoutCtx context.Context
		cancel     context.CancelFunc
		startTime  time.Time
		regularErr = utils.NewResettaleError()
	)

	startTime = time.Now()

	timeoutCtx, cancel = context.WithTimeout(ctx, state.Timeout())
	defer cancel()

	s.l.Info("setting up state", zap.String("state", state.Name()))

	// Set up the state for execution
	err := state.Setup(timeoutCtx)
	if err != nil {
		s.l.Error("received error during setup",
			zap.String("state", state.Name()),
			zap.Error(err))

		regularErr.Set(errors.Wrapf(err, "setup (%s)", state.Name()))
	}

	// Check for fatal error after setup
	err = state.FatalError()
	if err != nil {
		s.l.Error("encountered fatal error", zap.String("state", state.Name()), zap.Error(err))

		s.fatalErr.Set(errors.Wrapf(err, "fatal error during setup (%s)", state.Name()))
	}

	// If we encounter an error during setup, do not call run.
	if regularErr.Err() != nil {
		s.progress.StateDuration = append(s.progress.StateDuration, time.Since(startTime))

		return regularErr.Err()
	}

	timeoutCtx, cancel = context.WithTimeout(ctx, state.Timeout())
	defer cancel()

	s.l.Info("running state", zap.String("state", state.Name()))

	// Run the state logic
	err = state.Run(timeoutCtx)
	if err != nil {
		regularErr.Set(errors.Wrapf(err, "run (%s)", state.Name()))
	}

	s.progress.StateDuration = append(s.progress.StateDuration, time.Since(startTime))

	// Check for fatal error after run
	err = state.FatalError()
	if err != nil {
		s.l.Error("encountered fatal error", zap.String("state", state.Name()), zap.Error(err))

		s.fatalErr.Set(errors.Wrapf(err, "fatal error during run (%s)", state.Name()))
	}

	return regularErr.Err()
}

func (s *Sequencer) processResults(ctx context.Context, state State, err error) error {
	// TODO: integrate ResultProcessor
	return nil
}
