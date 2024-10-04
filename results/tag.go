package results

import (
	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
)

// Tag is a single test tag
type Tag struct {
	Description   string `yaml:"description"`
	CompOpString  string `yaml:"compOp"`
	CompOp        ComparisonOperator
	UpperLimit    any    `yaml:"upperLimit,omitempty"`
	LowerLimit    any    `yaml:"lowerLimit,omitempty"`
	ExpectedValue any    `yaml:"expectedValue,omitempty"`
	Unit          string `yaml:"unit"`
}

// IsPassing checks if the value passes the tag
func (t *Tag) IsPassing(value any) (bool, error) {
	if t.CompOp == Log {
		return true, nil
	}

	switch v := value.(type) {
	case bool:
		return isPassingBool(v, t.CompOp, t.ExpectedValue)
	case int:
		return isPassingNumeric(v, t.CompOp, t.ExpectedValue, t.UpperLimit, t.LowerLimit)
	case float64:
		return isPassingNumeric(v, t.CompOp, t.ExpectedValue, t.UpperLimit, t.LowerLimit)
	case string:
		return isPassingString(v, t.CompOp, t.ExpectedValue)
	default:
		return false, errors.Errorf("unsupported type (%T)", value)
	}
}

func isPassingBool(value bool, compOp ComparisonOperator, expectedValue any) (bool, error) {
	if compOp != Eq {
		return false, errors.New("boolean values only support equality comparison")
	}

	expected, ok := expectedValue.(bool)
	if !ok {
		return false, errors.New("expectedValue must be boolean for boolean comparison")
	}

	return value == expected, nil
}

func isPassingString(value string, compOp ComparisonOperator, expectedValue any) (bool, error) {
	if compOp != Eq {
		return false, errors.New("string values only support equality comparison")
	}

	expected, ok := expectedValue.(string)
	if !ok {
		return false, errors.New("expectedValue must be string for string comparison")
	}

	return value == expected, nil
}

func isPassingNumeric[T constraints.Ordered](
	value T,
	compOp ComparisonOperator,
	expectedValue, upperLimit, lowerLimit any,
) (bool, error) {
	switch compOp {
	case Eq:
		expected, ok := expectedValue.(T)
		if !ok {
			return false, errors.New("expectedValue type mismatch")
		}
		return value == expected, nil

	case Gele, Gtlt:
		upper, ok1 := upperLimit.(T)
		lower, ok2 := lowerLimit.(T)
		if !ok1 || !ok2 {
			return false, errors.New("limit values type mismatch")
		}
		if compOp == Gele {
			return value >= lower && value <= upper, nil
		}
		return value > lower && value < upper, nil

	case Gt, Ge:
		lower, ok := lowerLimit.(T)
		if !ok {
			return false, errors.New("lowerLimit type mismatch")
		}
		if compOp == Gt {
			return value > lower, nil
		}
		return value >= lower, nil

	case Lt, Le:
		upper, ok := upperLimit.(T)
		if !ok {
			return false, errors.New("upperLimit type mismatch")
		}
		if compOp == Lt {
			return value < upper, nil
		}
		return value <= upper, nil

	default:
		return false, errors.Errorf("unknown comparison operator (%v)", compOp.String())
	}
}
