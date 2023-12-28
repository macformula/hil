package flow

type Tag struct {
	TagID              string
	TagDescription     string
	ComparisonOperator string
	LowerLimit         float64
	UpperLimit         float64
	ExpectedValue      any
	Unit               string
}
