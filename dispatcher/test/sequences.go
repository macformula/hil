package test

import (
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/flow/test"
	"time"
)

var DoNothingSequence = flow.Sequence{
	Name: "Do Nothing 🥱",
	States: []flow.State{
		&test.DoNothingState{},
		&test.DoNothingState{},
		&test.DoNothingState{},
		&test.DoNothingState{},
	},
}

var SleepSequence = flow.Sequence{
	Name: "Sleeper 💤",
	States: []flow.State{
		&test.SleepState{SleepTime: 1 * time.Second},
		&test.SleepState{SleepTime: 5 * time.Second},
		&test.SleepState{SleepTime: 1 * time.Second},
		&test.SleepState{SleepTime: 2 * time.Second},
	},
}

var FatalErrorSequence = flow.Sequence{
	Name: "Fatal Error 💀",
	States: []flow.State{
		&test.SleepState{SleepTime: 2 * time.Second},
		&test.SleepState{SleepTime: 3 * time.Second},
		&test.SleepState{SleepTime: 1 * time.Second},
		&test.RunFatalErrorState{},
		&test.SleepState{SleepTime: 2 * time.Second},
		&test.SleepState{SleepTime: 3 * time.Second},
	},
}
