package sequencer

import (
	"context"
	"github.com/macformula/hilmatic/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sync"
)

type Sequencer struct {
	l         *zap.Logger
	sequence  Sequence
	currState State
	currIndex int

	stop bool
	err  *utils.ResettableError
	once sync.Once
}

func NewSequencer(s Sequence, l *zap.Logger) *Sequencer {
	return &Sequencer{
		l:        l.Named("sequencer"),
		sequence: s,
		err:      utils.NewResettaleError(),
	}
}

func (s *Sequencer) Start(ctx context.Context) error {
	if len(s.sequence) == 0 {
		return errors.New("sequence cannot be empty")
	}

	s.stop = false

	go s.runSequence(ctx)

	return nil
}

func (s *Sequencer) Progress() Progress {
	return Progress{
		CurrentState: s.currState,
		StateIndex:   s.currIndex,
		TotalStates:  len(s.sequence),
	}
}

// Stop will only stop after current state completion.
func (s *Sequencer) Stop() {
	s.stop = true
}

func (s *Sequencer) runSequence(ctx context.Context) {
	for idx, state := range s.sequence {
		if s.stop {
			s.err.Set(errors.Errorf("stop called before (%s)", state.Name()))
			return
		}

		// Check stop before setting current state
		s.currState = state
		s.currIndex = idx

		s.l.Info("starting next state", zap.String("state", state.Name()))

		err := state.Start(ctx)
		if err != nil {
			s.err.Set(errors.Wrapf(err, "start (%s)", state.Name()))
		}

		err = state.FatalError()
		if err != nil {
			s.l.Info("encountered fatal error", zap.String("state", state.Name()), zap.Error(err))
			s.err.Set(errors.Wrapf(err, "fatal error (%s)", state.Name()))

			return
		}
	}
}
