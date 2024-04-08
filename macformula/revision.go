package macformula

//go:generate enumer -type=Revision revision.go
// Revision is used to differentiate between pinouts.
type Revision int

const (
	Ev5 Revision = iota
	MockTest
)
