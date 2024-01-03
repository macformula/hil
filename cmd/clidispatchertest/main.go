package main

import (
	"context"
	"github.com/macformula/hil/dispatcher"
	"github.com/macformula/hil/flow"
	ftest "github.com/macformula/hil/flow/test"
	"github.com/macformula/hil/orchestrator"
	otest "github.com/macformula/hil/orchestrator/test"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

const (
	_loggerName = "main.log"
)

var mutex chan struct{}

func main() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{_loggerName}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	rp := ftest.NewSimpleResultProcessor(logger)
	s := flow.NewSequencer(rp, logger)
	d := dispatcher.NewCliDispatcher(logger)
	otest.NewSimpleDispatcher(5*time.Second, 10*time.Second)
	o := orchestrator.NewOrchestrator(s, logger, d)

	err = o.Open(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "orchestrator open"))
	}

	err = o.Run(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "orchestrator run"))
	}

	logger.Info("shutdown main program")
}

//func main() {
//	d := dispatcher.NewCliDispatcher(zap.L())
//	err := d.Open(context.Background())
//	if err != nil {
//		return
//	}
//
//	mutex = make(chan struct{}, 1)
//	mutex <- struct{}{}
//	go mimicProgress(d)
//
//	time.Sleep(10 * time.Second)
//	// give fatal error status
//	d.Status() <- orchestrator.StatusSignal{
//		OrchestratorState: orchestrator.FatalError,
//		TestId:            orchestrator.TestId{},
//		Progress:          flow.Progress{},
//		QueueLength:       1,
//		FatalError:        nil,
//	}
//	<-d.Quit()
//	time.Sleep(1 * time.Second)
//}
//
//// waits on start signal, then sends 3 progress and a results
//func mimicProgress(d orchestrator.DispatcherIface) {
//	start := d.Start()
//	cancel := d.CancelTest()
//
//	for {
//		select {
//		case signal := <-start:
//			go mimicRunning(d, signal)
//		case signal := <-cancel:
//			log.Printf("Cancelled signal %s", signal)
//		}
//	}
//}
//
//func mimicRunning(d orchestrator.DispatcherIface, signal orchestrator.StartSignal) {
//	<-mutex
//	progress1 := flow.Progress{
//		CurrentState:  &test.DoNothingState{},
//		StateIndex:    1,
//		Sequence:      flow.Sequence{},
//		StatePassed:   []bool{true},
//		StateDuration: []time.Duration{time.Second},
//	}
//
//	progress2 := flow.Progress{
//		CurrentState:  &test.DoNothingState{},
//		StateIndex:    2,
//		Sequence:      flow.Sequence{},
//		StatePassed:   []bool{true, true},
//		StateDuration: []time.Duration{time.Second, 2 * time.Second},
//	}
//
//	progress3 := flow.Progress{
//		CurrentState:  &test.DoNothingState{},
//		StateIndex:    2,
//		Sequence:      flow.Sequence{},
//		StatePassed:   []bool{true, true, false},
//		StateDuration: []time.Duration{time.Second, 2 * time.Second, 500 * time.Millisecond},
//	}
//
//	progress4 := flow.Progress{
//		CurrentState:  &test.DoNothingState{},
//		StateIndex:    3,
//		Sequence:      flow.Sequence{},
//		StatePassed:   []bool{true, true, false, true},
//		StateDuration: []time.Duration{time.Second, 2 * time.Second, 500 * time.Millisecond, 3 * time.Second},
//	}
//
//	progress5 := flow.Progress{
//		CurrentState:  &test.DoNothingState{},
//		StateIndex:    3,
//		Sequence:      flow.Sequence{},
//		StatePassed:   []bool{true, true, false, true, false},
//		StateDuration: []time.Duration{time.Second, 2 * time.Second, 500 * time.Millisecond, 3 * time.Second, 750 * time.Millisecond},
//	}
//
//	time.Sleep(1 * time.Second)
//	d.Status() <- orchestrator.StatusSignal{
//		OrchestratorState: orchestrator.Running,
//		TestId:            signal.TestId,
//		Progress:          progress1,
//		QueueLength:       1,
//		FatalError:        nil,
//	}
//
//	time.Sleep(2 * time.Second)
//	d.Status() <- orchestrator.StatusSignal{
//		OrchestratorState: orchestrator.Running,
//		TestId:            signal.TestId,
//		Progress:          progress2,
//		QueueLength:       1,
//		FatalError:        nil,
//	}
//
//	time.Sleep(500 * time.Millisecond)
//	d.Status() <- orchestrator.StatusSignal{
//		OrchestratorState: orchestrator.Running,
//		TestId:            signal.TestId,
//		Progress:          progress3,
//		QueueLength:       1,
//		FatalError:        nil,
//	}
//
//	time.Sleep(3 * time.Second)
//	d.Status() <- orchestrator.StatusSignal{
//		OrchestratorState: orchestrator.Running,
//		TestId:            signal.TestId,
//		Progress:          progress4,
//		QueueLength:       1,
//		FatalError:        nil,
//	}
//
//	time.Sleep(750 * time.Millisecond)
//	d.Status() <- orchestrator.StatusSignal{
//		OrchestratorState: orchestrator.Running,
//		TestId:            signal.TestId,
//		Progress:          progress5,
//		QueueLength:       1,
//		FatalError:        nil,
//	}
//
//	time.Sleep(1 * time.Second)
//	d.Status() <- orchestrator.StatusSignal{
//		OrchestratorState: orchestrator.Idle,
//		TestId:            orchestrator.TestId{},
//		Progress:          flow.Progress{},
//		QueueLength:       0,
//		FatalError:        nil,
//	}
//
//	time.Sleep(1 * time.Second)
//
//	d.Results() <- orchestrator.ResultsSignal{
//		TestId:     signal.TestId,
//		IsPassing:  true,
//		FailedTags: nil,
//	}
//
//	time.Sleep(3 * time.Second)
//	mutex <- struct{}{}
//}
