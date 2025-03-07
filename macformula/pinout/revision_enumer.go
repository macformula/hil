// Code generated by "enumer -type=Revision revision.go"; DO NOT EDIT.

package pinout

import (
	"fmt"
	"strings"
)

const _RevisionName = "Ev5MockTestSilSgTest"

var _RevisionIndex = [...]uint8{0, 3, 11, 14, 20}

const _RevisionLowerName = "ev5mocktestsilsgtest"

func (i Revision) String() string {
	if i < 0 || i >= Revision(len(_RevisionIndex)-1) {
		return fmt.Sprintf("Revision(%d)", i)
	}
	return _RevisionName[_RevisionIndex[i]:_RevisionIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _RevisionNoOp() {
	var x [1]struct{}
	_ = x[Ev5-(0)]
	_ = x[MockTest-(1)]
	_ = x[Sil-(2)]
	_ = x[SgTest-(3)]
}

var _RevisionValues = []Revision{Ev5, MockTest, Sil, SgTest}

var _RevisionNameToValueMap = map[string]Revision{
	_RevisionName[0:3]:        Ev5,
	_RevisionLowerName[0:3]:   Ev5,
	_RevisionName[3:11]:       MockTest,
	_RevisionLowerName[3:11]:  MockTest,
	_RevisionName[11:14]:      Sil,
	_RevisionLowerName[11:14]: Sil,
	_RevisionName[14:20]:      SgTest,
	_RevisionLowerName[14:20]: SgTest,
}

var _RevisionNames = []string{
	_RevisionName[0:3],
	_RevisionName[3:11],
	_RevisionName[11:14],
	_RevisionName[14:20],
}

// RevisionString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func RevisionString(s string) (Revision, error) {
	if val, ok := _RevisionNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _RevisionNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Revision values", s)
}

// RevisionValues returns all values of the enum
func RevisionValues() []Revision {
	return _RevisionValues
}

// RevisionStrings returns a slice of all String values of the enum
func RevisionStrings() []string {
	strs := make([]string, len(_RevisionNames))
	copy(strs, _RevisionNames)
	return strs
}

// IsARevision returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Revision) IsARevision() bool {
	for _, v := range _RevisionValues {
		if i == v {
			return true
		}
	}
	return false
}
