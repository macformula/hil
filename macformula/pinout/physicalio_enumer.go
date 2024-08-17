// Code generated by "enumer -type=PhysicalIo physicalio.go"; DO NOT EDIT.

package pinout

import (
	"fmt"
	"strings"
)

const _PhysicalIoName = "UnknownPhysicalIoDebugLedEnDashboardEnHvilLedEnBrakeLightEnStatusLedEnRtdsEnAccelPedalPosition1AccelPedalPosition2SuspensionTravel1SuspensionTravel2SteeringAngleStartButtonNWheelSpeedLeftAWheelSpeedLeftBWheelSpeedRightAWheelSpeedRightBMotorControllerPrechargeEnInverterSwitchEnAccumulatorEnShutdownCircuitEnTsalEnRaspiEnFrontControllerEnSpeedgoatEnMotorControllerEnImuGpsEnDcdcValidDcdcEnDcdcEnLedPowerTrainPumpEnPowertrainFanEnHvCurrentSenseLvController3v3RefVoltageFrontController3v3RefVoltageHvilDisableHvilFeedbackGlvmsDisableIndicatorLedIndicatorButton"

var _PhysicalIoIndex = [...]uint16{0, 17, 27, 38, 47, 59, 70, 76, 95, 114, 131, 148, 161, 173, 188, 203, 219, 235, 261, 277, 290, 307, 313, 320, 337, 348, 365, 373, 382, 388, 397, 413, 428, 442, 467, 495, 506, 518, 530, 542, 557}

const _PhysicalIoLowerName = "unknownphysicaliodebugledendashboardenhvilledenbrakelightenstatusledenrtdsenaccelpedalposition1accelpedalposition2suspensiontravel1suspensiontravel2steeringanglestartbuttonnwheelspeedleftawheelspeedleftbwheelspeedrightawheelspeedrightbmotorcontrollerprechargeeninverterswitchenaccumulatorenshutdowncircuitentsalenraspienfrontcontrollerenspeedgoatenmotorcontrollerenimugpsendcdcvaliddcdcendcdcenledpowertrainpumpenpowertrainfanenhvcurrentsenselvcontroller3v3refvoltagefrontcontroller3v3refvoltagehvildisablehvilfeedbackglvmsdisableindicatorledindicatorbutton"

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
	_ = x[DebugLedEn-(1)]
	_ = x[DashboardEn-(2)]
	_ = x[HvilLedEn-(3)]
	_ = x[BrakeLightEn-(4)]
	_ = x[StatusLedEn-(5)]
	_ = x[RtdsEn-(6)]
	_ = x[AccelPedalPosition1-(7)]
	_ = x[AccelPedalPosition2-(8)]
	_ = x[SuspensionTravel1-(9)]
	_ = x[SuspensionTravel2-(10)]
	_ = x[SteeringAngle-(11)]
	_ = x[StartButtonN-(12)]
	_ = x[WheelSpeedLeftA-(13)]
	_ = x[WheelSpeedLeftB-(14)]
	_ = x[WheelSpeedRightA-(15)]
	_ = x[WheelSpeedRightB-(16)]
	_ = x[MotorControllerPrechargeEn-(17)]
	_ = x[InverterSwitchEn-(18)]
	_ = x[AccumulatorEn-(19)]
	_ = x[ShutdownCircuitEn-(20)]
	_ = x[TsalEn-(21)]
	_ = x[RaspiEn-(22)]
	_ = x[FrontControllerEn-(23)]
	_ = x[SpeedgoatEn-(24)]
	_ = x[MotorControllerEn-(25)]
	_ = x[ImuGpsEn-(26)]
	_ = x[DcdcValid-(27)]
	_ = x[DcdcEn-(28)]
	_ = x[DcdcEnLed-(29)]
	_ = x[PowerTrainPumpEn-(30)]
	_ = x[PowertrainFanEn-(31)]
	_ = x[HvCurrentSense-(32)]
	_ = x[LvController3v3RefVoltage-(33)]
	_ = x[FrontController3v3RefVoltage-(34)]
	_ = x[HvilDisable-(35)]
	_ = x[HvilFeedback-(36)]
	_ = x[GlvmsDisable-(37)]
	_ = x[IndicatorLed-(38)]
	_ = x[IndicatorButton-(39)]
}

var _PhysicalIoValues = []PhysicalIo{UnknownPhysicalIo, DebugLedEn, DashboardEn, HvilLedEn, BrakeLightEn, StatusLedEn, RtdsEn, AccelPedalPosition1, AccelPedalPosition2, SuspensionTravel1, SuspensionTravel2, SteeringAngle, StartButtonN, WheelSpeedLeftA, WheelSpeedLeftB, WheelSpeedRightA, WheelSpeedRightB, MotorControllerPrechargeEn, InverterSwitchEn, AccumulatorEn, ShutdownCircuitEn, TsalEn, RaspiEn, FrontControllerEn, SpeedgoatEn, MotorControllerEn, ImuGpsEn, DcdcValid, DcdcEn, DcdcEnLed, PowerTrainPumpEn, PowertrainFanEn, HvCurrentSense, LvController3v3RefVoltage, FrontController3v3RefVoltage, HvilDisable, HvilFeedback, GlvmsDisable, IndicatorLed, IndicatorButton}

var _PhysicalIoNameToValueMap = map[string]PhysicalIo{
	_PhysicalIoName[0:17]:         UnknownPhysicalIo,
	_PhysicalIoLowerName[0:17]:    UnknownPhysicalIo,
	_PhysicalIoName[17:27]:        DebugLedEn,
	_PhysicalIoLowerName[17:27]:   DebugLedEn,
	_PhysicalIoName[27:38]:        DashboardEn,
	_PhysicalIoLowerName[27:38]:   DashboardEn,
	_PhysicalIoName[38:47]:        HvilLedEn,
	_PhysicalIoLowerName[38:47]:   HvilLedEn,
	_PhysicalIoName[47:59]:        BrakeLightEn,
	_PhysicalIoLowerName[47:59]:   BrakeLightEn,
	_PhysicalIoName[59:70]:        StatusLedEn,
	_PhysicalIoLowerName[59:70]:   StatusLedEn,
	_PhysicalIoName[70:76]:        RtdsEn,
	_PhysicalIoLowerName[70:76]:   RtdsEn,
	_PhysicalIoName[76:95]:        AccelPedalPosition1,
	_PhysicalIoLowerName[76:95]:   AccelPedalPosition1,
	_PhysicalIoName[95:114]:       AccelPedalPosition2,
	_PhysicalIoLowerName[95:114]:  AccelPedalPosition2,
	_PhysicalIoName[114:131]:      SuspensionTravel1,
	_PhysicalIoLowerName[114:131]: SuspensionTravel1,
	_PhysicalIoName[131:148]:      SuspensionTravel2,
	_PhysicalIoLowerName[131:148]: SuspensionTravel2,
	_PhysicalIoName[148:161]:      SteeringAngle,
	_PhysicalIoLowerName[148:161]: SteeringAngle,
	_PhysicalIoName[161:173]:      StartButtonN,
	_PhysicalIoLowerName[161:173]: StartButtonN,
	_PhysicalIoName[173:188]:      WheelSpeedLeftA,
	_PhysicalIoLowerName[173:188]: WheelSpeedLeftA,
	_PhysicalIoName[188:203]:      WheelSpeedLeftB,
	_PhysicalIoLowerName[188:203]: WheelSpeedLeftB,
	_PhysicalIoName[203:219]:      WheelSpeedRightA,
	_PhysicalIoLowerName[203:219]: WheelSpeedRightA,
	_PhysicalIoName[219:235]:      WheelSpeedRightB,
	_PhysicalIoLowerName[219:235]: WheelSpeedRightB,
	_PhysicalIoName[235:261]:      MotorControllerPrechargeEn,
	_PhysicalIoLowerName[235:261]: MotorControllerPrechargeEn,
	_PhysicalIoName[261:277]:      InverterSwitchEn,
	_PhysicalIoLowerName[261:277]: InverterSwitchEn,
	_PhysicalIoName[277:290]:      AccumulatorEn,
	_PhysicalIoLowerName[277:290]: AccumulatorEn,
	_PhysicalIoName[290:307]:      ShutdownCircuitEn,
	_PhysicalIoLowerName[290:307]: ShutdownCircuitEn,
	_PhysicalIoName[307:313]:      TsalEn,
	_PhysicalIoLowerName[307:313]: TsalEn,
	_PhysicalIoName[313:320]:      RaspiEn,
	_PhysicalIoLowerName[313:320]: RaspiEn,
	_PhysicalIoName[320:337]:      FrontControllerEn,
	_PhysicalIoLowerName[320:337]: FrontControllerEn,
	_PhysicalIoName[337:348]:      SpeedgoatEn,
	_PhysicalIoLowerName[337:348]: SpeedgoatEn,
	_PhysicalIoName[348:365]:      MotorControllerEn,
	_PhysicalIoLowerName[348:365]: MotorControllerEn,
	_PhysicalIoName[365:373]:      ImuGpsEn,
	_PhysicalIoLowerName[365:373]: ImuGpsEn,
	_PhysicalIoName[373:382]:      DcdcValid,
	_PhysicalIoLowerName[373:382]: DcdcValid,
	_PhysicalIoName[382:388]:      DcdcEn,
	_PhysicalIoLowerName[382:388]: DcdcEn,
	_PhysicalIoName[388:397]:      DcdcEnLed,
	_PhysicalIoLowerName[388:397]: DcdcEnLed,
	_PhysicalIoName[397:413]:      PowerTrainPumpEn,
	_PhysicalIoLowerName[397:413]: PowerTrainPumpEn,
	_PhysicalIoName[413:428]:      PowertrainFanEn,
	_PhysicalIoLowerName[413:428]: PowertrainFanEn,
	_PhysicalIoName[428:442]:      HvCurrentSense,
	_PhysicalIoLowerName[428:442]: HvCurrentSense,
	_PhysicalIoName[442:467]:      LvController3v3RefVoltage,
	_PhysicalIoLowerName[442:467]: LvController3v3RefVoltage,
	_PhysicalIoName[467:495]:      FrontController3v3RefVoltage,
	_PhysicalIoLowerName[467:495]: FrontController3v3RefVoltage,
	_PhysicalIoName[495:506]:      HvilDisable,
	_PhysicalIoLowerName[495:506]: HvilDisable,
	_PhysicalIoName[506:518]:      HvilFeedback,
	_PhysicalIoLowerName[506:518]: HvilFeedback,
	_PhysicalIoName[518:530]:      GlvmsDisable,
	_PhysicalIoLowerName[518:530]: GlvmsDisable,
	_PhysicalIoName[530:542]:      IndicatorLed,
	_PhysicalIoLowerName[530:542]: IndicatorLed,
	_PhysicalIoName[542:557]:      IndicatorButton,
	_PhysicalIoLowerName[542:557]: IndicatorButton,
}

var _PhysicalIoNames = []string{
	_PhysicalIoName[0:17],
	_PhysicalIoName[17:27],
	_PhysicalIoName[27:38],
	_PhysicalIoName[38:47],
	_PhysicalIoName[47:59],
	_PhysicalIoName[59:70],
	_PhysicalIoName[70:76],
	_PhysicalIoName[76:95],
	_PhysicalIoName[95:114],
	_PhysicalIoName[114:131],
	_PhysicalIoName[131:148],
	_PhysicalIoName[148:161],
	_PhysicalIoName[161:173],
	_PhysicalIoName[173:188],
	_PhysicalIoName[188:203],
	_PhysicalIoName[203:219],
	_PhysicalIoName[219:235],
	_PhysicalIoName[235:261],
	_PhysicalIoName[261:277],
	_PhysicalIoName[277:290],
	_PhysicalIoName[290:307],
	_PhysicalIoName[307:313],
	_PhysicalIoName[313:320],
	_PhysicalIoName[320:337],
	_PhysicalIoName[337:348],
	_PhysicalIoName[348:365],
	_PhysicalIoName[365:373],
	_PhysicalIoName[373:382],
	_PhysicalIoName[382:388],
	_PhysicalIoName[388:397],
	_PhysicalIoName[397:413],
	_PhysicalIoName[413:428],
	_PhysicalIoName[428:442],
	_PhysicalIoName[442:467],
	_PhysicalIoName[467:495],
	_PhysicalIoName[495:506],
	_PhysicalIoName[506:518],
	_PhysicalIoName[518:530],
	_PhysicalIoName[530:542],
	_PhysicalIoName[542:557],
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
