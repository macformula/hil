package test

import (
	"context"
	"github.com/macformula/hil/orchestrator"

	"github.com/macformula/hil/flow"
)

type simpleDispatcher struct {
	start             chan flow.Sequence
	quit              chan struct{}
	recoverFromFatal  chan struct{}
	progress          chan flow.Progress
	orchestratorState orchestrator.State
	stop              chan struct{}
}

func newSimpleDispatcher() *simpleDispatcher {
	return &simpleDispatcher{
		start:    make(chan flow.Sequence),
		quit:     make(chan struct{}),
		progress: make(chan flow.Progress),
	}
}

func (s *simpleDispatcher) Close() error {
	return nil
}

func (s *simpleDispatcher) Open(_ context.Context) error {
	return nil
}

func (s *simpleDispatcher) StartSequence(seq flow.Sequence) {
	s.start <- seq
}

func (s *simpleDispatcher) QuitSequence() {
	s.stop <- struct{}{}
	s.quit <- struct{}{}
}

func (s *simpleDispatcher) CommandRecoverFromFatal() {
}

func (s *simpleDispatcher) Start() <-chan flow.Sequence {
	return s.start
}

func (s *simpleDispatcher) Quit() <-chan struct{} {
	return s.quit
}

func (s *simpleDispatcher) RecoverFromFatal() <-chan struct{} {
	return s.recoverFromFatal
}

func (s *simpleDispatcher) Progress() chan flow.Progress {
	return s.progress
}

func (s *simpleDispatcher) OrchestratorState(state chan orchestrator.State) {
	go func(state chan orchestrator.State) {
		for {
			select {
			case s.orchestratorState = <-state:
			case <-s.stop:
				return
			}
		}
	}(state)
}
