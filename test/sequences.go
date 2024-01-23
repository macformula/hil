package test

import (
	"time"

	"github.com/macformula/hil/flow"
)

var (
	Sequences = []flow.Sequence{DoNothingSequence, SleepSequence, FatalErrorSequence, ErrorSequence}
)

var DoNothingSequence = flow.Sequence{
	Name: "Do Nothing 🥱",
	Desc: "Wow... it does nothing",
	States: []flow.State{
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

var SleepSequence = flow.Sequence{
	Name: "Sleeper 💤",
	Desc: "zzz",
	States: []flow.State{
		&SleepState{SleepTime: 1 * time.Second},
		&SleepState{SleepTime: 5 * time.Second},
		&SleepState{SleepTime: 1 * time.Second},
		&SleepState{SleepTime: 2 * time.Second},
	},
}

var FatalErrorSequence = flow.Sequence{
	Name: "Fatal Error 💀",
	Desc: "This will fatal... duh",
	States: []flow.State{
		&SleepState{SleepTime: 2 * time.Second},
		&SleepState{SleepTime: 3 * time.Second},
		&SleepState{SleepTime: 1 * time.Second},
		&RunFatalErrorState{},
		&SleepState{SleepTime: 2 * time.Second},
		&SleepState{SleepTime: 3 * time.Second},
	},
}

var ErrorSequence = flow.Sequence{
	Name: "Normal Error",
	Desc: "ERROR ERROR ERROR",
	States: []flow.State{
		&SleepState{SleepTime: 2 * time.Second},
		&SleepState{SleepTime: 3 * time.Second},
		&RunErrorState{},
		&SleepState{SleepTime: 2 * time.Second},
		&SleepState{SleepTime: 3 * time.Second},
	},
}
