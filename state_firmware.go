package hil

import (
	"context"
	"github.com/macformula/hil/utils"
)

const (
	_stateName = "firmware"
)

type Firmware struct {
	err utils.ResettableError
}

func (f *Firmware) Name() string {
	return _stateName
}

func (f *Firmware) Start(a App, ctx context.Context) error {

}

func (f *Firmware) FatalError() error {
	//TODO implement me
	panic("implement me")
}
