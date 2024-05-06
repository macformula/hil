package state

import (
	"github.com/macformula/hil/macformula"
	"go.uber.org/zap"
	"time"

	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/test"
)

func GetSequences(a *macformula.App, l *zap.Logger) []flow.Sequence {
	ret := make([]flow.Sequence, len(_sequenceConstructors))

	for i, seqConstructor := range _sequenceConstructors {
		ret[i] = seqConstructor(a, l)
	}

	return ret
}

type sequenceConstructor = func(a *macformula.App, l *zap.Logger) flow.Sequence

var _sequenceConstructors = []sequenceConstructor{
	newTracerSequence,
	newDoNothingSequence,
	newSleepSequence,
	newFatalErrorSequence,
	newPanicSequence,
	newErrorSequence,
}

func newTracerSequence(a *macformula.App, l *zap.Logger) flow.Sequence {
	return flow.Sequence{
		Name: "Can Tracer ✍️",
		Desc: "Obtains a can trace",
		States: []flow.State{
			newSetup(a, l),
			newSleep(10 * time.Second),
			newCleanup(a, l),
		},
	}
}

func newDoNothingSequence(_ *macformula.App, _ *zap.Logger) flow.Sequence {
	return flow.Sequence{
		Name: "Do Nothing 🥱",
		Desc: "Wow... it does nothing",
		States: []flow.State{
			&test.DoNothingState{},
			&test.DoNothingState{},
			&test.DoNothingState{},
			&test.DoNothingState{},
			&test.DoNothingState{},
			&test.DoNothingState{},
			&test.DoNothingState{},
			&test.DoNothingState{},
			&test.DoNothingState{},
			&test.DoNothingState{},
			&test.DoNothingState{},
			&test.DoNothingState{},
		},
	}
}
func newSleepSequence(_ *macformula.App, _ *zap.Logger) flow.Sequence {
	return flow.Sequence{
		Name: "Sleeper 💤",
		Desc: "zzz",
		States: []flow.State{
			&test.SleepState{SleepTime: 1 * time.Second},
			&test.SleepState{SleepTime: 5 * time.Second},
			&test.SleepState{SleepTime: 1 * time.Second},
			&test.SleepState{SleepTime: 2 * time.Second},
		},
	}
}
func newFatalErrorSequence(_ *macformula.App, _ *zap.Logger) flow.Sequence {
	return flow.Sequence{
		Name: "Fatal Error 💀",
		Desc: "This will fatal... duh",
		States: []flow.State{
			&test.SleepState{SleepTime: 2 * time.Second},
			&test.SleepState{SleepTime: 3 * time.Second},
			&test.SleepState{SleepTime: 1 * time.Second},
			&test.RunFatalErrorState{},
			&test.SleepState{SleepTime: 2 * time.Second},
			&test.SleepState{SleepTime: 3 * time.Second},
		},
	}
}

func newPanicSequence(_ *macformula.App, _ *zap.Logger) flow.Sequence {
	return flow.Sequence{
		Name: "Panic 😨",
		Desc: "This will panic the hil app.",
		States: []flow.State{
			&test.SleepState{SleepTime: 1 * time.Second},
			&test.SleepState{SleepTime: 1 * time.Second},
			&test.SleepState{SleepTime: 1 * time.Second},
			&test.PanicState{},
			&test.SleepState{SleepTime: 2 * time.Second},
			&test.SleepState{SleepTime: 3 * time.Second},
		},
	}
}
func newErrorSequence(_ *macformula.App, _ *zap.Logger) flow.Sequence {
	return flow.Sequence{
		Name: "Normal Error",
		Desc: "ERROR ERROR ERROR",
		States: []flow.State{
			&test.SleepState{SleepTime: 2 * time.Second},
			&test.SleepState{SleepTime: 3 * time.Second},
			&test.RunErrorState{},
			&test.SleepState{SleepTime: 2 * time.Second},
			&test.SleepState{SleepTime: 3 * time.Second},
		},
	}
}
