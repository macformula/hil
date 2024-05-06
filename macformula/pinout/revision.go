package pinout

// Revision is used to differentiate between pinouts.
//
//go:generate enumer -type=Revision "revision.go"
type Revision int

const (
	Ev5 Revision = iota
	MockTest
	Sil
)
