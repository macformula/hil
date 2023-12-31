package test

import (
	"context"
	"github.com/macformula/hil/flow"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"testing"
	"time"
)

// TestSequencer_Nominal tests the nominal sequence cases where no errors are returned and all states finish executing
// by their required deadlines.
func TestSequencer_Nominal(t *testing.T) {
	seqs := []flow.Sequence{
		{
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&SleepState{1 * time.Second},
		},
		{
			&SleepState{100 * time.Millisecond},
			&SleepState{100 * time.Millisecond},
			&SleepState{100 * time.Millisecond},
			&SleepState{100 * time.Millisecond},
		},
		{
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
			&DoNothingState{},
		},
	}

	sequencer, l := SetupSequencer(t, "nominal_test.log")

	for _, seq := range seqs {
		err := sequencer.Run(context.Background(), seq)
		if err != nil {
			l.Error("sequencer run", zap.Error(err))
			t.Error(errors.Wrap(err, "sequencer run"))
		}

		err = sequencer.FatalError()
		if err != nil {
			l.Error("received fatal error", zap.Error(err))
			t.Error(errors.Wrap(err, "received fatal error"))
		}
	}

	// Give time to flush logger
	time.Sleep(1 * time.Second)
}

// TestSequencer_FatalError tests the case where a state returns a fatal error.
func TestSequencer_FatalError(t *testing.T) {
	tests := []struct {
		seq                 flow.Sequence
		expectedMaxStateIdx int
	}{
		{
			seq: flow.Sequence{
				&RunFatalErrorState{}, // Make sure the last two states do not run
			},
			expectedMaxStateIdx: 0,
		},
		{
			seq: flow.Sequence{
				&SetupFatalErrorState{}, // Make sure the last two states do not run
				&DoNothingState{},
				&DoNothingState{},
				&DoNothingState{},
			},
			expectedMaxStateIdx: 0,
		},
		{
			seq: flow.Sequence{
				&DoNothingState{},
				&DoNothingState{},
				&SetupFatalErrorState{}, // Make sure the last two states do not run
				&DoNothingState{},
				&DoNothingState{},
				&DoNothingState{},
			},
			expectedMaxStateIdx: 2,
		},
		{
			seq: flow.Sequence{
				&RunFatalErrorState{}, // Make sure the last two states do not run
			},
			expectedMaxStateIdx: 0,
		},
		{
			seq: flow.Sequence{
				&RunFatalErrorState{}, // Make sure the last two states do not run
				&DoNothingState{},
				&DoNothingState{},
				&DoNothingState{},
			},
			expectedMaxStateIdx: 0,
		},
		{
			seq: flow.Sequence{
				&DoNothingState{},
				&DoNothingState{},
				&SetupFatalErrorState{}, // Make sure the last two states do not run
				&DoNothingState{},
				&DoNothingState{},
				&DoNothingState{},
			},
			expectedMaxStateIdx: 2,
		},
	}

	sequencer, l := SetupSequencer(t, "fatal_run_error_test.log")

	prog := make(chan flow.Progress)
	sub := sequencer.SubscribeToProgress(prog)
	defer sub.Unsubscribe()

	var (
		maxStateIdx = -1
		stop        = make(chan struct{})
	)

	go func(maxStateIdx *int) {
		for {
			select {
			case <-stop:
				return
			case p := <-prog:
				*maxStateIdx = p.StateIndex
			case <-sub.Err():
				return
			}
		}
	}(&maxStateIdx)

	for _, test := range tests {
		err := sequencer.Run(context.Background(), test.seq)
		if err != nil {
			l.Error("sequencer run", zap.Error(err))
		}

		err = sequencer.FatalError()
		if err != nil {
			l.Error("received fatal error", zap.Error(err))
		} else {
			t.Error("expected fatal error")
		}

		if test.expectedMaxStateIdx != maxStateIdx {
			t.Errorf("expected test to be cut short after %d states, ran till %d states",
				test.expectedMaxStateIdx+1,
				maxStateIdx+1)
		}
	}

	close(stop)
	// Give time to flush logger
	time.Sleep(1 * time.Second)
}

// TestSequencer_Timeout tests the case where a run or setup call hangs forever.
func TestSequencer_Timeout(t *testing.T) {
	seqs := []flow.Sequence{
		{
			&DoNothingState{},
			&DoNothingState{},
			&RunForeverState{},
			&DoNothingState{},
			&DoNothingState{},
		},
		{
			&DoNothingState{},
			&DoNothingState{},
			&SetupForeverState{},
			&DoNothingState{},
			&DoNothingState{},
		},
	}
	sequencer, l := SetupSequencer(t, "timeout_test.log")

	for _, seq := range seqs {
		err := sequencer.Run(context.Background(), seq)
		if err != nil {
			l.Error("sequencer run", zap.Error(err))
		}

		err = sequencer.FatalError()
		if err != nil {
			l.Error("received fatal error", zap.Error(err))
		}
	}

	// Give time to flush logger
	time.Sleep(1 * time.Second)
}

// TestSequencer_Timeout tests the case where a run or setup call hangs forever.
func TestSequencer_EmptySequence(t *testing.T) {
	seq := flow.Sequence{}

	sequencer, _ := SetupSequencer(t, "timeout_test.log")

	err := sequencer.Run(context.Background(), seq)
	if err == nil {
		t.Error("expected error")
	}

	// Give time to flush logger
	time.Sleep(1 * time.Second)
}
