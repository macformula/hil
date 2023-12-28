package hil

import "github.com/macformula/hil/flow"

// TODO: generate from tags.yaml

type FirmwareTags struct {
	FrontControllerFlashed flow.Tag
	LvControllerFlashed    flow.Tag
	TmsFlashed             flow.Tag
}

type LvNominal struct {
	TimingGood  flow.Tag
	FanPwm      flow.Tag
	PumpControl flow.Tag
}

var FwTags = FirmwareTags{
	FrontControllerFlashed: flow.Tag{"FW001"},
	LvControllerFlashed:    flow.Tag{"FW002"},
	TmsFlashed:             flow.Tag{"FW003"},
}
