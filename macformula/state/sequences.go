package state

import (
	"time"

	"go.uber.org/zap"

	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/macformula"
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
	newSleepSequence,
	newDoNothingSequence,
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
			newNothing(),
			newNothing(),
			newNothing(),
			newNothing(),
			newNothing(),
			newNothing(),
			newNothing(),
		},
	}
}
func newSleepSequence(_ *macformula.App, _ *zap.Logger) flow.Sequence {
	return flow.Sequence{
		Name: "Sleeper 💤",
		Desc: "zzz",
		States: []flow.State{
			newSleep(1 * time.Second),
			newSleep(5 * time.Second),
			newSleep(2 * time.Second),
			newSleep(1 * time.Second),
		},
	}
}
