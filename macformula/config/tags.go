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
