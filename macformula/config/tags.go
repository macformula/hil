package config

import "github.com/macformula/hil/flow"

type FirmwareTagCollection struct {
	FrontControllerFlashed flow.Tag
	LvControllerFlashed    flow.Tag
	TmsFlashed             flow.Tag
}

var FirmwareTags = FirmwareTagCollection{
	FrontControllerFlashed: flow.Tag{ID: "FW001", Description: "Front controller flashed."},
	LvControllerFlashed:    flow.Tag{ID: "FW002"},
	TmsFlashed:             flow.Tag{ID: "FW003"},
}

type TestTagCollection struct {
	TestTag1 flow.Tag
}

var TestTags = TestTagCollection{
	TestTag1: flow.Tag{ID: "TEST001", Description: "Test tag."},
}

type LvStartupTagCollection struct {
	PowerCycledTestBench                              flow.Tag
	TsalEnabled                                       flow.Tag
	TsalTimeToEnableMs                                flow.Tag
	RaspiEnabled                                      flow.Tag
	RaspiTimeToEnableMs                               flow.Tag
	FrontControllerEnabled                            flow.Tag
	FrontControllerTimeToEnableMs                     flow.Tag
	SpeedgoatEnabled                                  flow.Tag
	SpeedgoatTimeToEnableMs                           flow.Tag
	AccumulatorEnabled                                flow.Tag
	AccumulatorTimeToEnableMs                         flow.Tag
	MotorPrechageEnabled                              flow.Tag
	MotorPrechargeTimeToEnableMs                      flow.Tag
	MotorControllerEnabled                            flow.Tag
	MotorControllerTimeToEnable                       flow.Tag
	ShutdownCircuitEnabledBeforeCan                   flow.Tag
	ShutdownCircuitEnabledBeforeOpenContactors        flow.Tag
	ShutdownCircuitEnabled                            flow.Tag
	ShutdownCircuitTimeToEnable                       flow.Tag
	DcdcEnabledBeforeContactorsClosed                 flow.Tag
	DcdcEnabledAfterContactorsClosed                  flow.Tag
	InverterSwitchEnabledBeforeCan                    flow.Tag
	InverterSwitchEnabledBeforeClosedContactors       flow.Tag
	InverterSwitchEnabledBeforeFrontControllerCommand flow.Tag
	InverterSwitchEnabled                             flow.Tag
	InverterSwitchTimeToEnable                        flow.Tag
}

var LvStartupTags = LvStartupTagCollection{
	PowerCycledTestBench: flow.Tag{
		ID:          "LVSTART001",
		Description: "Successfully power cycled the testbench.",
	},
	TsalEnabled: flow.Tag{
		ID:          "LVSTART002",
		Description: "TSAL indicator enabled after power cycle.",
	},
	TsalTimeToEnableMs: flow.Tag{
		ID:          "LVSTART003",
		Description: "TSAL indicator time to enable after startup (ms).",
	},
	RaspiEnabled: flow.Tag{
		ID:          "LVSTART004",
		Description: "Raspi indicator enabled after tsal.",
	},
	RaspiTimeToEnableMs: flow.Tag{
		ID:          "LVSTART005",
		Description: "Raspi time to enable after tsal (ms).",
	},
	FrontControllerEnabled: flow.Tag{
		ID:          "LVSTART006",
		Description: "Front controller enabled after raspi.",
	},
	FrontControllerTimeToEnableMs: flow.Tag{
		ID:          "LVSTART007",
		Description: "Front controller time to enable after raspi (ms).",
	},
	SpeedgoatEnabled: flow.Tag{
		ID:          "LVSTART008",
		Description: "Speedgoat enabled after front controller.",
	},
	SpeedgoatTimeToEnableMs: flow.Tag{
		ID:          "LVSTART009",
		Description: "Speedgoat time to enable after front controller (ms).",
	},
	AccumulatorEnabled: flow.Tag{
		ID:          "LVSTART010",
		Description: "Accumulator enabled after speedgoat.",
	},
	AccumulatorTimeToEnableMs: flow.Tag{
		ID:          "LVSTART011",
		Description: "Accumulator time to enable after speedgoat (ms).",
	},
	MotorPrechageEnabled: flow.Tag{
		ID:          "LVSTART012",
		Description: "Motor controller precharge enabled after accumulator.",
	},
	MotorPrechargeTimeToEnableMs: flow.Tag{
		ID:          "LVSTART013",
		Description: "Motor controller precharge time to enable after accumulator (ms).",
	},
	MotorControllerEnabled: flow.Tag{
		ID:          "LVSTART014",
		Description: "Motor controller enabled after motor controller precharge.",
	},
	MotorControllerTimeToEnable: flow.Tag{
		ID:          "LVSTART015",
		Description: "Motor controller time to enable after motor controller precharge (ms).",
	},
	ShutdownCircuitEnabledBeforeCan: flow.Tag{
		ID:          "LVSTART016",
		Description: "Shutdown circuit enabled before can contactors command sent.",
	},
	ShutdownCircuitEnabledBeforeOpenContactors: flow.Tag{
		ID:          "LVSTART017",
		Description: "Shutdown circuit enabled before contactors commanded open.",
	},
	ShutdownCircuitEnabled: flow.Tag{
		ID:          "LVSTART018",
		Description: "Shutdown circuit enabled after contactors commanded open.",
	},
	ShutdownCircuitTimeToEnable: flow.Tag{
		ID:          "LVSTART019",
		Description: " Shutdown circuit time to enable after contactos commanded open (ms).",
	},
	DcdcEnabledBeforeContactorsClosed: flow.Tag{
		ID:          "LVSTART020",
		Description: "Dcdc enabled before contactors commanded closed (ms).",
	},
	DcdcEnabledAfterContactorsClosed: flow.Tag{
		ID:          "LVSTART021",
		Description: "Dcdc enabled after contactors commanded closed (ms).",
	},
	InverterSwitchEnabledBeforeCan: flow.Tag{
		ID:          "LVSTART022",
		Description: "Inverter switch enabled before can contactors command sent.",
	},
	InverterSwitchEnabledBeforeClosedContactors: flow.Tag{
		ID:          "LVSTART023",
		Description: "Inverter switch enabled before contactors commanded closed.",
	},
	InverterSwitchEnabledBeforeFrontControllerCommand: flow.Tag{
		ID:          "LVSTART024",
		Description: "Inverter switch enabled before commanded to enable by the front controller.",
	},
	InverterSwitchEnabled: flow.Tag{
		ID:          "LVSTART025",
		Description: "Inverter switch enabled after commanded to enable by the front controller.",
	},
	InverterSwitchTimeToEnable: flow.Tag{
		ID:          "LVSTART026",
		Description: "Inverter switch time to enable after commanded to enable by the front controller (ms).",
	},
}
