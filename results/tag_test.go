package results

import (
	"testing"
)

func TestTag_IsPassing(t *testing.T) {
	tagDb := map[string]Tag{
		"integerEquality": {
			Description:   "Integer equality check",
			CompOp:        Eq,
			UpperLimit:    0, // not used
			LowerLimit:    0, // not used
			ExpectedValue: 5,
			Unit:          "",
		},
		"floatRangeInclusive": {
			Description:   "Float range check (inclusive)",
			CompOp:        Gele,
			UpperLimit:    10.0,
			LowerLimit:    5.0,
			ExpectedValue: 0.0, // not used
			Unit:          "unit",
		},
		"integerRangeExclusive": {
			Description:   "Integer range check (exclusive)",
			CompOp:        Gtlt,
			UpperLimit:    10,
			LowerLimit:    0,
			ExpectedValue: 5, // not used
			Unit:          "unit",
		},
		"greaterThan": {
			Description:   "Greater than check",
			CompOp:        Gt,
			UpperLimit:    0, // not used
			LowerLimit:    5,
			ExpectedValue: 0, // not used
			Unit:          "unit",
		},
		"lessThan": {
			Description:   "Less than check",
			CompOp:        Lt,
			UpperLimit:    5,
			LowerLimit:    0, // not used
			ExpectedValue: 0, // not used
			Unit:          "unit",
		},
		"greaterThanOrEqual": {
			Description:   "Greater than or equal check",
			CompOp:        Ge,
			UpperLimit:    0, // not used
			LowerLimit:    5,
			ExpectedValue: 0, // not used
			Unit:          "unit",
		},
		"lessThanOrEqual": {
			Description:   "Less than or equal check",
			CompOp:        Le,
			UpperLimit:    5,
			LowerLimit:    0, // not used
			ExpectedValue: 0, // not used
			Unit:          "unit",
		},
		"logOnly": {
			Description:   "Log only (always passes)",
			CompOp:        Log,
			UpperLimit:    0, // not used
			LowerLimit:    0, // not used
			ExpectedValue: 0, // not used
			Unit:          "unit",
		},
		"stringEquality": {
			Description:   "String equality check",
			CompOp:        Eq,
			UpperLimit:    "", // not used
			LowerLimit:    "", // not used
			ExpectedValue: "test",
			Unit:          "",
		},
		"booleanEquality": {
			Description:   "Boolean equality check",
			CompOp:        Eq,
			UpperLimit:    false, // not used
			LowerLimit:    false, // not used
			ExpectedValue: true,
			Unit:          "",
		},
	}

	testCases := []struct {
		name     string
		tagID    string
		value    any
		expected bool
	}{
		{"integerEquality passing", "integerEquality", 5, true},
		{"integerEquality failing", "integerEquality", 6, false},
		{"floatRangeInclusive passing lower bound", "floatRangeInclusive", 5.0, true},
		{"floatRangeInclusive passing upper bound", "floatRangeInclusive", 10.0, true},
		{"floatRangeInclusive passing middle", "floatRangeInclusive", 7.5, true},
		{"floatRangeInclusive failing below", "floatRangeInclusive", 4.9, false},
		{"floatRangeInclusive failing above", "floatRangeInclusive", 10.1, false},
		{"integerRangeExclusive passing", "integerRangeExclusive", 5, true},
		{"integerRangeExclusive failing lower bound", "integerRangeExclusive", 0, false},
		{"integerRangeExclusive failing upper bound", "integerRangeExclusive", 10, false},
		{"greaterThan passing", "greaterThan", 6, true},
		{"greaterThan failing", "greaterThan", 5, false},
		{"lessThan passing", "lessThan", 4, true},
		{"lessThan failing", "lessThan", 5, false},
		{"greaterThanOrEqual passing equal", "greaterThanOrEqual", 5, true},
		{"greaterThanOrEqual passing above", "greaterThanOrEqual", 6, true},
		{"greaterThanOrEqual failing", "greaterThanOrEqual", 4, false},
		{"lessThanOrEqual passing equal", "lessThanOrEqual", 5, true},
		{"lessThanOrEqual passing below", "lessThanOrEqual", 4, true},
		{"lessThanOrEqual failing", "lessThanOrEqual", 6, false},
		{"logOnly always passing", "logOnly", 100, true},
		{"stringEquality passing", "stringEquality", "test", true},
		{"stringEquality failing", "stringEquality", "wrong", false},
		{"booleanEquality passing", "booleanEquality", true, true},
		{"booleanEquality failing", "booleanEquality", false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tag := tagDb[tc.tagID]
			result, err := tag.IsPassing(tc.value)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}
