package main

import (
	"context"
	dispatcher "github.com/macformula/hil/dispatcher"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/flow/test"
	"github.com/macformula/hil/orchestrator"
	"go.uber.org/zap"
	"time"
)

func main() {
	d := dispatcher.NewCliDispatcher(zap.L())
	err := d.Open(context.Background())
	if err != nil {
		return
	}

	go mimicProgress(d)

	time.Sleep(20 * time.Second)
}

// waits on start signal, then sends 3 progress and a results
func mimicProgress(d orchestrator.Dispatcher) {
	start := d.Start()

	for {
		select {
		case signal := <-start:
			progress1 := flow.Progress{
				CurrentState:  &test.DoNothingState{},
				StateIndex:    2,
				Complete:      true,
				Sequence:      flow.Sequence{},
				StatePassed:   []bool{true},
				StateDuration: []time.Duration{time.Second},
			}

			progress2 := flow.Progress{
				CurrentState:  &test.DoNothingState{},
				StateIndex:    2,
				Complete:      true,
				Sequence:      flow.Sequence{},
				StatePassed:   []bool{true, true},
				StateDuration: []time.Duration{time.Second, 2 * time.Second},
			}

			progress3 := flow.Progress{
				CurrentState:  &test.DoNothingState{},
				StateIndex:    2,
				Complete:      true,
				Sequence:      flow.Sequence{},
				StatePassed:   []bool{true, true, false},
				StateDuration: []time.Duration{time.Second, 2 * time.Second, 500 * time.Millisecond},
			}

			d.Status() <- orchestrator.StatusSignal{
				OrchestratorState: orchestrator.Running,
				TestId:            signal.TestId,
				Progress:          progress1,
				QueueLength:       1,
				FatalError:        nil,
			}
			time.Sleep(1 * time.Second)
			//time.Sleep(progress1.StateDuration[len(progress1.StateDuration)-1])

			d.Status() <- orchestrator.StatusSignal{
				OrchestratorState: orchestrator.Running,
				TestId:            signal.TestId,
				Progress:          progress2,
				QueueLength:       1,
				FatalError:        nil,
			}
			//time.Sleep(progress1.StateDuration[len(progress2.StateDuration)-1])
			time.Sleep(5 * time.Second)
			d.Status() <- orchestrator.StatusSignal{
				OrchestratorState: orchestrator.Running,
				TestId:            signal.TestId,
				Progress:          progress3,
				QueueLength:       1,
				FatalError:        nil,
			}
			//time.Sleep(progress1.StateDuration[len(progress3.StateDuration)-1])
			//d.Close()
			time.Sleep(1 * time.Second)

			d.Results() <- orchestrator.ResultsSignal{
				TestId:     signal.TestId,
				IsPassing:  true,
				FailedTags: nil,
			}

			time.Sleep(3 * time.Second)
		}
	}
}
