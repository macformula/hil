package dispatcher

import (
	"github.com/macformula/hil/flow"
)

const (
	showLastResults = 5
)

type sequenceItem flow.Sequence

// Title is the sequence name.
func (i sequenceItem) Title() string {
	return i.Name
}

// Description is the sequence description.
func (i sequenceItem) Description() string {
	return i.Desc
}

// FilterValue is the sequences name.
func (i sequenceItem) FilterValue() string {
	return i.Name
}

// Metadata is not yet used.
func (i sequenceItem) getMetaData() map[string]string {
	metaData := make(map[string]string)
	return metaData
}
