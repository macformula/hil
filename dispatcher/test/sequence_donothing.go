package test

import (
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/flow/test"
)

var DoNothingSequence = flow.Sequence{
	&test.DoNothingState{},
	&test.DoNothingState{},
	&test.DoNothingState{},
	&test.DoNothingState{},
}
