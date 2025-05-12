package state

import (
	"time"

	"go.uber.org/zap"

	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/macformula"
)

// GetSequences returns all the sequences that can be run.
func GetSequences(a *macformula.App, l *zap.Logger) []flow.Sequence {
	ret := make([]flow.Sequence, len(_sequenceConstructors))

	for i, seqConstructor := range _sequenceConstructors {
		ret[i] = seqConstructor(a, l)
	}

	return ret
}

type sequenceConstructor = func(a *macformula.App, l *zap.Logger) flow.Sequence

var _sequenceConstructors = []sequenceConstructor{
	newLvControllerSequence,
	newTracerSequence,
	newSleepSequence,
	newDoNothingSequence,
	newBasicIoSequence,
}

func newLvControllerSequence(a *macformula.App, l *zap.Logger) flow.Sequence {
	return flow.Sequence{
		Name: "Lv Controller Sequence ⚡",
		Desc: "Tests the lv controller.",
		States: []flow.State{
			newSetup(a, l),
			newLvStartup(a, l),
			newCleanup(a, l),
		},
	}
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

func newDoNothingSequence(a *macformula.App, l *zap.Logger) flow.Sequence {
	return flow.Sequence{
		Name: "Do Nothing 🥱",
		Desc: "Wow... it does nothing",
		States: []flow.State{
			newSetup(a, l),
			newNothing(),
			newNothing(),
			newNothing(),
			newNothing(),
			newNothing(),
			newNothing(),
			newNothing(),
			newCleanup(a, l),
		},
	}
}

func newSleepSequence(a *macformula.App, l *zap.Logger) flow.Sequence {
	return flow.Sequence{
		Name: "Sleeper 💤",
		Desc: "zzz",
		States: []flow.State{
			newSetup(a, l),
			newSleep(1 * time.Second),
			newSleep(5 * time.Second),
			newSleep(2 * time.Second),
			newSleep(1 * time.Second),
			newCleanup(a, l),
		},
	}
}

func newBasicIoSequence(a *macformula.App, l *zap.Logger) flow.Sequence {
	return flow.Sequence{
		Name: "BasicIo 🧪",
		Desc: "Test sequence for BasicIo firmware project.",
		States: []flow.State{
			newSetup(a, l),
			newBasicIo(a, l),
			newCleanup(a, l),
		},
	}
}
