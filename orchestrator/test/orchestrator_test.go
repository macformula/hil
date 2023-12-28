package test

import (
	"context"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/flow/test"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"testing"

	"github.com/macformula/hil/orchestrator"
)

func TestOrchestrator_Nominal(t *testing.T) {
	const (
		_logFileName = "orchestrator_nominal.log"
		_numRuns     = 10
	)

	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{_logFileName}
	logger, err := cfg.Build()
	if err != nil {
		t.Error(errors.Wrap(err, "build logger"))
	}
	defer logger.Sync()

	dispatcher := newSimpleDispatcher()

	orchestrator := orchestrator.NewOrchestrator(logger, dispatcher)

	sequence := flow.Sequence{
		&test.DoNothingState{},
		&test.DoNothingState{},
		&test.DoNothingState{},
		&test.DoNothingState{},
		&test.DoNothingState{},
	}

	err = orchestrator.Open(context.Background())
	if err != nil {
		t.Error(errors.Wrap(err, "orchestrator open"))
	}

	// This is not the typical use case
	go func() {
		for i := 0; i < _numRuns; i++ {
			dispatcher.StartSequence(sequence)

			for {
				progress := <-dispatcher.Progress()

				if progress.Complete {
					break
				}
			}
		}

		dispatcher.QuitSequence()
	}()

	err = orchestrator.Run(context.Background())
	if err != nil {
		t.Error(errors.Wrap(err, "orchestrator run"))
		return
	}

	err = orchestrator.Close()
	if err != nil {
		t.Error(errors.Wrap(err, "orchestrator close"))
		return
	}
}

func TestOrchestrator_FatalError(t *testing.T) {
	const (
		_logFileName = "orchestrator_fatal.log"
		_numRuns     = 10
	)

	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{_logFileName}
	logger, err := cfg.Build()
	if err != nil {
		t.Error(errors.Wrap(err, "build logger"))
	}
	defer logger.Sync()

	dispatcher := newSimpleDispatcher()

	orchestrator := orchestrator.NewOrchestrator(logger, dispatcher)

	sequence := flow.Sequence{
		&test.DoNothingState{},
		&test.DoNothingState{},
		&test.DoNothingState{},
		&test.DoNothingState{},
		&test.RunFatalErrorState{},
	}

	err = orchestrator.Open(context.Background())
	if err != nil {
		t.Error(errors.Wrap(err, "orchestrator open"))
	}

	// This is not the typical use case
	go func() {
		for i := 0; i < _numRuns; i++ {
			dispatcher.StartSequence(sequence)

		}

		dispatcher.QuitSequence()
	}()

	err = orchestrator.Run(context.Background())
	if err != nil {
		t.Error(errors.Wrap(err, "orchestrator run"))
		return
	}

	err = orchestrator.Close()
	if err != nil {
		t.Error(errors.Wrap(err, "orchestrator close"))
		return
	}
}
