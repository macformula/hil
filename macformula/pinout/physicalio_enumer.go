// Code generated by "enumer -type=PhysicalIo physicalio.go"; DO NOT EDIT.

package pinout

import (
	"fmt"
	"strings"
)

const _PhysicalIoName = "UnknownPhysicalIoLvEnableButtonReadyToDriveButtonHvilOkAcceleratorPedalPosition1AcceleratorPedalPosition2AccumulatorCurrentLvController3v3RefVoltageIndicatorLedIndicatorButton"

var _PhysicalIoIndex = [...]uint8{0, 17, 31, 49, 55, 80, 105, 123, 148, 160, 175}

const _PhysicalIoLowerName = "unknownphysicaliolvenablebuttonreadytodrivebuttonhvilokacceleratorpedalposition1acceleratorpedalposition2accumulatorcurrentlvcontroller3v3refvoltageindicatorledindicatorbutton"

func (i PhysicalIo) String() string {
	if i < 0 || i >= PhysicalIo(len(_PhysicalIoIndex)-1) {
		return fmt.Sprintf("PhysicalIo(%d)", i)
	}
	return _PhysicalIoName[_PhysicalIoIndex[i]:_PhysicalIoIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _PhysicalIoNoOp() {
	var x [1]struct{}
	_ = x[UnknownPhysicalIo-(0)]
	_ = x[LvEnableButton-(1)]
	_ = x[ReadyToDriveButton-(2)]
	_ = x[HvilOk-(3)]
	_ = x[AcceleratorPedalPosition1-(4)]
	_ = x[AcceleratorPedalPosition2-(5)]
	_ = x[AccumulatorCurrent-(6)]
	_ = x[LvController3v3RefVoltage-(7)]
	_ = x[IndicatorLed-(8)]
	_ = x[IndicatorButton-(9)]
}

var _PhysicalIoValues = []PhysicalIo{UnknownPhysicalIo, LvEnableButton, ReadyToDriveButton, HvilOk, AcceleratorPedalPosition1, AcceleratorPedalPosition2, AccumulatorCurrent, LvController3v3RefVoltage, IndicatorLed, IndicatorButton}

var _PhysicalIoNameToValueMap = map[string]PhysicalIo{
	_PhysicalIoName[0:17]:         UnknownPhysicalIo,
	_PhysicalIoLowerName[0:17]:    UnknownPhysicalIo,
	_PhysicalIoName[17:31]:        LvEnableButton,
	_PhysicalIoLowerName[17:31]:   LvEnableButton,
	_PhysicalIoName[31:49]:        ReadyToDriveButton,
	_PhysicalIoLowerName[31:49]:   ReadyToDriveButton,
	_PhysicalIoName[49:55]:        HvilOk,
	_PhysicalIoLowerName[49:55]:   HvilOk,
	_PhysicalIoName[55:80]:        AcceleratorPedalPosition1,
	_PhysicalIoLowerName[55:80]:   AcceleratorPedalPosition1,
	_PhysicalIoName[80:105]:       AcceleratorPedalPosition2,
	_PhysicalIoLowerName[80:105]:  AcceleratorPedalPosition2,
	_PhysicalIoName[105:123]:      AccumulatorCurrent,
	_PhysicalIoLowerName[105:123]: AccumulatorCurrent,
	_PhysicalIoName[123:148]:      LvController3v3RefVoltage,
	_PhysicalIoLowerName[123:148]: LvController3v3RefVoltage,
	_PhysicalIoName[148:160]:      IndicatorLed,
	_PhysicalIoLowerName[148:160]: IndicatorLed,
	_PhysicalIoName[160:175]:      IndicatorButton,
	_PhysicalIoLowerName[160:175]: IndicatorButton,
}

var _PhysicalIoNames = []string{
	_PhysicalIoName[0:17],
	_PhysicalIoName[17:31],
	_PhysicalIoName[31:49],
	_PhysicalIoName[49:55],
	_PhysicalIoName[55:80],
	_PhysicalIoName[80:105],
	_PhysicalIoName[105:123],
	_PhysicalIoName[123:148],
	_PhysicalIoName[148:160],
	_PhysicalIoName[160:175],
}

// PhysicalIoString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func PhysicalIoString(s string) (PhysicalIo, error) {
	if val, ok := _PhysicalIoNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _PhysicalIoNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to PhysicalIo values", s)
}

// PhysicalIoValues returns all values of the enum
func PhysicalIoValues() []PhysicalIo {
	return _PhysicalIoValues
}

// PhysicalIoStrings returns a slice of all String values of the enum
func PhysicalIoStrings() []string {
	strs := make([]string, len(_PhysicalIoNames))
	copy(strs, _PhysicalIoNames)
	return strs
}

// IsAPhysicalIo returns "true" if the value is listed in the enum definition. "false" otherwise
func (i PhysicalIo) IsAPhysicalIo() bool {
	for _, v := range _PhysicalIoValues {
		if i == v {
			return true
		}
	}
	return false
}