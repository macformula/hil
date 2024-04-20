package stflash

import (
	"context"
	"github.com/macformula/hil/fwutils"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_resetCmd = "reset"
	_flashCmd = "st-flash"

	_serialArg = "--serial"
	_writeArg  = "write"
	_writeAddr = "0x8000000"

	_loggerName     = "flasher"
	_defaultTimeout = 2 * time.Second
)

// enforce interface implementation
var _ fwutils.FlasherIface = NewFlasher(zap.Logger{})

type Flasher struct {
	currentBoardId string
	boardActive    bool

	l *zap.Logger
}

// NewFlasher returns a st-flash flasher
func NewFlasher(l zap.Logger) *Flasher {
	return &Flasher{
		l: l.Named(_loggerName),
	}
}

// Connect establishes a connection with an STM32
func (f *Flasher) Connect(boardId string) error {
	f.l.Info("checking for active target")

	// Create a context with a timeout of 2 seconds
	ctx, cancel := context.WithTimeout(context.Background(), _defaultTimeout)
	defer cancel()

	openCmd := exec.CommandContext(ctx, _resetCmd, _serialArg, boardId)

	_, err := openCmd.CombinedOutput() // search for that serial number... different states for each ecu
	if err != nil {
		return errors.Wrap(err, "connect to stm32")
	} // uhubctl //also check for no failed to enter SWO

	//f.currentBoardId = strings.TrimSpace(string(output)
	f.boardActive = true

	f.l.Info("target found", zap.String("board id", f.currentBoardId))

	return nil
}

// Flash uses the stlink driver to flash the target with a provided binary
func (f *Flasher) Flash(binaryPath string) error {
	if f.boardActive {
		f.l.Info("attempting to flash")

		flashCmd := exec.Command(_flashCmd, _writeArg, binaryPath, _writeAddr)

		_, err := flashCmd.CombinedOutput()
		if err != nil {
			return errors.Wrap(err, "flash stm32")
		}

		f.l.Info("flash successful")

	} else {
		return errors.New("target is not connected")
	}
	return nil
}

// String returns the flasher type
func (f *Flasher) String() string {
	return "st-flash"
}

// Disconnect closes the connection with the target
func (f *Flasher) Disconnect() error {
	f.currentBoardId = ""
	f.boardActive = false

	return nil
}
