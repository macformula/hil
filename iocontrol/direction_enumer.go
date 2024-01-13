// Code generated by "enumer -type=Direction types.go"; DO NOT EDIT.

package iocontrol

import (
	"fmt"
	"strings"
)

const _DirectionName = "InputOutput"

var _DirectionIndex = [...]uint8{0, 5, 11}

const _DirectionLowerName = "inputoutput"

func (i Direction) String() string {
	if i < 0 || i >= Direction(len(_DirectionIndex)-1) {
		return fmt.Sprintf("Direction(%d)", i)
	}
	return _DirectionName[_DirectionIndex[i]:_DirectionIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _DirectionNoOp() {
	var x [1]struct{}
	_ = x[Input-(0)]
	_ = x[Output-(1)]
}

var _DirectionValues = []Direction{Input, Output}

var _DirectionNameToValueMap = map[string]Direction{
	_DirectionName[0:5]:       Input,
	_DirectionLowerName[0:5]:  Input,
	_DirectionName[5:11]:      Output,
	_DirectionLowerName[5:11]: Output,
}

var _DirectionNames = []string{
	_DirectionName[0:5],
	_DirectionName[5:11],
}

// DirectionString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func DirectionString(s string) (Direction, error) {
	if val, ok := _DirectionNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _DirectionNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Direction values", s)
}

// DirectionValues returns all values of the enum
func DirectionValues() []Direction {
	return _DirectionValues
}

// DirectionStrings returns a slice of all String values of the enum
func DirectionStrings() []string {
	strs := make([]string, len(_DirectionNames))
	copy(strs, _DirectionNames)
	return strs
}

// IsADirection returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Direction) IsADirection() bool {
	for _, v := range _DirectionValues {
		if i == v {
			return true
		}
	}
	return false
}