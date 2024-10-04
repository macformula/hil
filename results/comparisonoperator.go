package results

//go:generate enumer -type=ComparisonOperator
type ComparisonOperator int

const (
	// Eq is Equal to
	Eq ComparisonOperator = iota
	// Gele is Greater than or Equal to and Less than or Equal to
	Gele
	// Gtlt is Greater than and Less than
	Gtlt
	// Gt is Greater than
	Gt
	// Lt is Less than
	Lt
	// Ge is Greater than
	Ge
	// Le is Less than
	Le
	// Log is logging the value without comparison
	Log
)
