package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/macformula/hil/flow"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// SetupSequencer is a helper function that is used in sequencer_test.go
func SetupSequencer(t *testing.T, logFileName string) (*flow.Sequencer, *zap.Logger) {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{logFileName}
	logger, err := cfg.Build()
	if err != nil {
		t.Error(errors.Wrap(err, "build logger"))
	}
	defer logger.Sync()

	sequencer := flow.NewSequencer(logger)

	go func(s *flow.Sequencer, l *zap.Logger) {
		prog := make(chan flow.Progress)
		sub := s.SubscribeToProgress(prog)
		defer sub.Unsubscribe()

		for {
			select {
			case err1 := <-sub.Err():
				l.Error("subscription error", zap.Error(err1))
			case p, ok := <-prog:
				if !ok {
					l.Error("progress channel closed")
					return
				}

				l.Info("received progress update",
					zap.String("completion", fmt.Sprintf("%f%%", float64(p.StateIndex)*100/float64(len(p.Sequence)))),
					zap.String("current_state", p.CurrentState.Name()),
				)

				// Log state durations on the last state
				if p.Complete {
					for idx, state := range p.Sequence {
						l.Info("state_duration",
							zap.String("state", state.Name()),
							zap.Int64("duration (ms)", p.StateDuration[idx].Milliseconds()),
						)
					}
				}
			}
		}
	}(sequencer, logger)

	// Give time for subscription
	time.Sleep(1 * time.Second)

	return sequencer, logger
}
