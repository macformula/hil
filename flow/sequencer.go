package flow

import (
	"context"
	"github.com/google/uuid"
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

	fatalErr   *utils.ResettableError
	regularErr *utils.ResettableError

	rp ResultProcessorIface

	failedTags []Tag
}

// NewSequencer returns a Sequencer object reference.
func NewSequencer(rp ResultProcessorIface, l *zap.Logger) *Sequencer {
	return &Sequencer{
		l:            l.Named(_loggerName),
		progressFeed: event.Feed{},
		fatalErr:     utils.NewResettaleError(),
		rp:           rp,
	}
}

// SubscribeToProgress subscribes to the progress of the Sequencer across its Sequence runs.
// The Progress channel gets updated whenever there is new information available.
func (s *Sequencer) SubscribeToProgress(progCh chan Progress) event.Subscription {
	return s.progressFeed.Subscribe(progCh)
}

// Run will run the sequence provided. FatalError must be called after Run to check for any non-recoverable errors.
func (s *Sequencer) Run(ctx context.Context, seq Sequence, testId uuid.UUID) (bool, []Tag, error) {
	if len(seq) == 0 {
		return false, nil, errors.New("sequence cannot be empty")
	}

	s.progress = Progress{
		CurrentState:  nil,
		StateDuration: make([]time.Duration, len(seq)),
		StatePassed:   make([]bool, len(seq)),
		StateIndex:    0,
		Sequence:      seq,
	}

	s.failedTags = []Tag{}

	isPassing, err := s.runSequence(ctx, seq, testId)
	if err != nil {
		return false, s.failedTags, errors.Wrap(err, "run sequence")
	}

	return isPassing, s.failedTags, nil
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

func (s *Sequencer) runSequence(ctx context.Context, seq Sequence, testId uuid.UUID) (bool, error) {
	for idx, state := range seq {
		s.progress.CurrentState = state
		s.progress.StateIndex = idx

		_ = s.progressFeed.Send(s.progress)

		s.l.Info("starting next state", zap.String("state", state.Name()))

		s.runState(ctx, state)

		s.l.Info("processing results", zap.String("state", state.Name()))

		continueSequence, err := s.processResults(ctx, state)
		if err != nil {
			return false, errors.Wrap(err, "process results")
		}

		if !continueSequence {
			s.l.Info("stopping sequence execution early")
			break
		}
	}

	s.l.Info("sequence complete")

	passingTest, err := s.rp.CompleteTest(ctx, testId)
	if err != nil {
		return false, errors.Wrap(err, "complete test")
	}

	return passingTest, nil
}

func (s *Sequencer) runState(ctx context.Context, state State) {
	var (
		timeoutCtx context.Context
		cancel     context.CancelFunc
		startTime  time.Time
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

		s.regularErr.Set(errors.Wrapf(err, "setup (%s)", state.Name()))
	}

	// Check for fatal error after setup
	err = state.FatalError()
	if err != nil {
		s.l.Error("encountered fatal error", zap.String("state", state.Name()), zap.Error(err))

		s.fatalErr.Set(errors.Wrapf(err, "fatal error during setup (%s)", state.Name()))
	}

	// If we encounter an error during setup, return early and do not call run.
	if s.regularErr.Err() != nil || s.fatalErr.Err() != nil {
		s.progress.StateDuration = append(s.progress.StateDuration, time.Since(startTime))

		return
	}

	timeoutCtx, cancel = context.WithTimeout(ctx, state.Timeout())
	defer cancel()

	s.l.Info("running state", zap.String("state", state.Name()))

	// Run the state logic
	err = state.Run(timeoutCtx)
	if err != nil {
		s.regularErr.Set(errors.Wrapf(err, "run (%s)", state.Name()))
	}

	s.progress.StateDuration = append(s.progress.StateDuration, time.Since(startTime))

	// Check for fatal error after run
	err = state.FatalError()
	if err != nil {
		s.l.Error("encountered fatal error", zap.String("state", state.Name()), zap.Error(err))

		s.fatalErr.Set(errors.Wrapf(err, "fatal error during run (%s)", state.Name()))
	}

	return
}

func (s *Sequencer) processResults(ctx context.Context, state State) (bool, error) {
	var (
		statePassed      = true
		continueSequence bool
	)

	if s.regularErr.Err() != nil {
		statePassed = false

		err := s.rp.EncounteredError(ctx, s.regularErr.Err())
		if err != nil {
			return false, errors.Wrap(err, "encountered error")
		}
	}

	if s.fatalErr.Err() != nil {
		err := s.rp.EncounteredError(ctx, s.fatalErr.Err())
		if err != nil {
			return false, errors.Wrap(err, "encountered error")
		}
	}

	results := state.GetResults()

	for tag, value := range results {
		isPassing, err := s.rp.SubmitTag(ctx, tag.TagID, value)
		if err != nil {
			return false, errors.Wrap(err, "submit tag")
		}

		s.progress.StatePassed = append(s.progress.StatePassed, isPassing)

		if !isPassing {
			statePassed = false
			s.failedTags = append(s.failedTags, tag)
		}
	}

	switch {
	// If encountered fatal error, should not continue.
	case state.FatalError() != nil:
		continueSequence = false
	// If state passed and did not get any regular errors, continue sequence.
	case statePassed && (s.regularErr.Err() == nil):
		continueSequence = true
	// If state encountered error or did not pass, but continue on fail is true, continue sequence.
	case state.ContinueOnFail():
		continueSequence = true
	default:
		continueSequence = false
	}

	return continueSequence, nil
}
