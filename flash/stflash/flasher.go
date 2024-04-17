package stflash

import (
	"github.com/macformula/hil/flash"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os/exec"
)

const (
	_infoCmd  = "st-info"
	_flashCmd = "st-flash"

	_serialArg = "--serial"

	_loggerName = "flasher"
)

// enforce interface implementation
var _ flash.FlasherIface = NewFlasher()

type Flasher struct {
	currentBoardId string

	l *zap.Logger
}

// NewFlasher returns a st-flash flasher
func NewFlasher(l zap.Logger) *Flasher {
	return &Flasher{
		l: l.Named(_loggerName),
	}
}

func (f *Flasher) Open() error {
	f.l.Info("checking for active target")

	open := exec.Command(_infoCmd, _serialArg)

	output, err := open.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "connect to stm32")
	}

	f.currentBoardId = string(output)

	f.l.Info("target found", zap.String("board id", f.currentBoardId))

	return nil
}

func (f *Flasher) Flash(binary string) error {
	//TODO implement me
	panic("implement me")
}

func (f *Flasher) String() string {
	return "st-flash"
}

func (f *Flasher) Close() error {
	//TODO implement me
	panic("implement me")
}
