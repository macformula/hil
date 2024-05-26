package flow

import (
	"context"
	"fmt" //delete later !!!
	"time"

	"github.com/ethereum/go-ethereum/event"
	"github.com/google/uuid"
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

	cancelCurrentTest context.CancelFunc

	testCanceled bool

	failedTags []Tag
	testErrors []error
}

// NewSequencer returns a Sequencer object reference.
func NewSequencer(rp ResultProcessorIface, l *zap.Logger) *Sequencer {
	return &Sequencer{
		l:            l.Named(_loggerName),
		progressFeed: event.Feed{},
		fatalErr:     utils.NewResettaleError(),
		regularErr:   utils.NewResettaleError(),
		rp:           rp,
	}
}

// SubscribeToProgress subscribes to the progress of the Sequencer across its Sequence runs.
// The Progress channel gets updated whenever there is new information available.
func (s *Sequencer) SubscribeToProgress(progCh chan Progress) event.Subscription {
	return s.progressFeed.Subscribe(progCh)
}

// Open will be called at the start of the app.
func (s *Sequencer) Open(ctx context.Context) error {
	err := s.rp.Open(ctx)
	if err != nil {
		return errors.Wrap(err, "open")
	}

	return nil
}

// Run will run the sequence provided. FatalError must be called after Run to check for any non-recoverable errors.
func (s *Sequencer) Run(
	ctx context.Context,
	seq Sequence,
	cancelTest chan struct{},
	testId uuid.UUID) (bool, []Tag, []error, error) {
	if len(seq.States) == 0 {
		return false, nil, []error{errors.New("sequence cannot be empty")}, errors.New("sequence cannot be empty")
	}

	s.progress = Progress{
		CurrentState:  nil,
		StateDuration: make([]time.Duration, 0),
		StatePassed:   make([]bool, 0),
		StateIndex:    0,
		Sequence:      seq,
	}

	s.failedTags = []Tag{}
	fmt.Println("Sequencer.Run: Starting sequence:", testId, " | ", seq.Name, " | ", seq.Desc, " | ", time.Now()) // Start of sequence
	store := NewStore()
	report := &Report{
		ID:           testId,
		SequenceName: seq.Name,
		DateTime:     time.Now(),
		Description:  seq.Desc, // Assuming seq has a Desc field
	}
	store.CreateReport(testId, report)

	isPassing, err := s.runSequence(ctx, seq, cancelTest, testId)
	if err != nil {
		testErrors := append(s.testErrors, errors.Wrap(err, "run sequence"))
		s.testErrors = []error{}

		return false, s.failedTags, testErrors, errors.Wrap(err, "run sequence")
	}

	testErrors := s.testErrors
	s.testErrors = []error{}
	return isPassing, s.failedTags, testErrors, nil
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

func (s *Sequencer) Close() error {
	s.l.Info("closing sequencer")
	err := s.rp.Close()
	if err != nil {
		return errors.Wrap(err, "result processor close")
	}

	return nil
}

func (s *Sequencer) runSequence(ctx context.Context, seq Sequence, cancelTest chan struct{}, testId uuid.UUID) (bool, error) {
	for idx, state := range seq.States {
		s.progress.CurrentState = state
		s.progress.StateIndex = idx

		_ = s.progressFeed.Send(s.progress)

		s.l.Info("starting next state", zap.String("state", state.Name()))
		s.runState(ctx, cancelTest, state)
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

	_ = s.progressFeed.Send(s.progress)

	passingTest, err := s.rp.CompleteTest(ctx, testId, seq.Name)
	s.l.Info("results processor CompleteTest call complete")
	if err != nil {
		return false, errors.Wrap(err, "complete test")
	}

	return passingTest, nil
}

func (s *Sequencer) runState(ctx context.Context, cancelTest chan struct{}, state State) {
	var (
		timeoutCtx context.Context
		startTime  time.Time
	)

	startTime = time.Now()

	timeoutCtx, s.cancelCurrentTest = context.WithTimeout(ctx, state.Timeout())
	defer s.cancelCurrentTest()

	go s.monitorCancelSignal(ctx, cancelTest)

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
		fmt.Println("Sequencer.runState: Setup failed for state:", state.Name(), "with error:", s.regularErr.Err(), s.fatalErr.Err()) // State setup failed
		return
	}

	timeoutCtx, s.cancelCurrentTest = context.WithTimeout(ctx, state.Timeout())
	defer s.cancelCurrentTest()

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

		s.testErrors = append(s.testErrors, s.regularErr.Err())

		err := s.rp.SubmitError(ctx, s.regularErr.Err())
		if err != nil {
			return false, errors.Wrap(err, "encountered error")
		}
	}

	if s.fatalErr.Err() != nil {
		statePassed = false

		err := s.rp.SubmitError(ctx, s.fatalErr.Err())
		if err != nil {
			return false, errors.Wrap(err, "encountered error")
		}
	}

	results := state.GetResults()

	for tag, value := range results {
		isPassing, err := s.rp.SubmitTag(ctx, tag.ID, value)
		if err != nil {
			return false, errors.Wrap(err, "submit tag")
		}
		if !isPassing {
			statePassed = false
			s.failedTags = append(s.failedTags, tag)
		}
	}

	s.progress.StatePassed = append(s.progress.StatePassed, statePassed)

	switch {
	// If test canceled should not continue to next states.
	case s.testCanceled:
		s.testCanceled = false
		continueSequence = false
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

	fmt.Println("Sequencer.processResults: State", state.Name(), "passed:", statePassed, " ", state.GetResults()) // State pass/fail
	s.regularErr.Reset()

	return continueSequence, nil
}

func (s *Sequencer) monitorCancelSignal(ctx context.Context, cancelTest chan struct{}) {
	select {
	case <-ctx.Done():
	case <-cancelTest:
		s.testCanceled = true
		s.cancelCurrentTest()
	}
}
