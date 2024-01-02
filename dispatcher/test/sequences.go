package test

import (
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/flow/test"
	"time"
)

var DoNothingSequence = flow.Sequence{
	&test.DoNothingState{},
	&test.DoNothingState{},
	&test.DoNothingState{},
	&test.DoNothingState{},
}

var SleepSequence = flow.Sequence{
	&test.SleepState{SleepTime: 5 * time.Second},
	&test.SleepState{SleepTime: 5 * time.Second},
	&test.SleepState{SleepTime: 1 * time.Second},
	&test.SleepState{SleepTime: 2 * time.Second},
}
