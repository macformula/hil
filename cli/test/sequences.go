package test

import (
	"time"

	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/flow/test"
)

var (
	Sequences = []flow.Sequence{DoNothingSequence, SleepSequence, FatalErrorSequence, ErrorSequence}
)

var DoNothingSequence = flow.Sequence{
	Name: "Do Nothing 🥱",
	Desc: "Wow... it does nothing",
	States: []flow.State{
		&test.DoNothingState{},
		&test.DoNothingState{},
		&test.DoNothingState{},
		&test.DoNothingState{},
	},
}

var SleepSequence = flow.Sequence{
	Name: "Sleeper 💤",
	Desc: "zzz",
	States: []flow.State{
		&test.SleepState{SleepTime: 1 * time.Second},
		&test.SleepState{SleepTime: 5 * time.Second},
		&test.SleepState{SleepTime: 1 * time.Second},
		&test.SleepState{SleepTime: 2 * time.Second},
	},
}

var FatalErrorSequence = flow.Sequence{
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

var ErrorSequence = flow.Sequence{
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
