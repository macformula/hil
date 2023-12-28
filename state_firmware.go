package hil

import (
	"github.com/macformula/hil/utils"
)

const (
	_stateName = "firmware"
)

type Firmware struct {
	err utils.ResettableError
}
