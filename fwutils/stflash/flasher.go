package stflash

import (
	"context"
	"fmt"
	"github.com/macformula/hil/fwutils"
	"os/exec"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	// stlink commands
	_stFlashCmd  = "st-flash"
	_stResetCmd  = "reset"
	_stResetArg  = "--reset"
	_stSerialArg = "--serial"
	_stWriteArg  = "write"
	_stWriteAddr = "0x8000000"

	// uhubctl commands
	_uhubCmd        = "uhubctl"
	_uhubPwrArg     = "-a"
	_uhubOff        = "off"
	_uhubOn         = "on"
	_uhubPortNumArg = "-p"
	_uhubHubArg     = "-l"
	_uhubDefaultHub = "1-1"

	_loggerName     = "flasher"
	_defaultTimeout = time.Second
)

// enforce interface implementation
var _ fwutils.FlasherIface = NewFlasher(zap.Logger{}, map[fwutils.Ecu]string{})

// Flasher implements the FlasherIface for stlink flashing
type Flasher struct {
	currentBoardId string
	boardActive    bool
	ecuSerialMap   map[fwutils.Ecu]string

	l *zap.Logger
}

// NewFlasher returns a st-flash flasher
func NewFlasher(l zap.Logger, ecuSerialMap map[fwutils.Ecu]string) *Flasher {
	return &Flasher{
		l:            l.Named(_loggerName),
		ecuSerialMap: ecuSerialMap,
	}
}

// Connect establishes a connection with an STM32
func (f *Flasher) Connect(ecu fwutils.Ecu) error {
	f.l.Info("checking for active target", zap.String("target", ecu.String()))

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = _defaultTimeout

	firstAttempt := true

	// wrap the Connect logic in a retry loop
	err := backoff.Retry(func() error {

		if !firstAttempt {
			if err := f.PowerCycleStm(ecu); err != nil {
				return err
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), _defaultTimeout)
		defer cancel()

		// Attempt to reset the target
		openCmd := exec.CommandContext(ctx, _stFlashCmd, _stSerialArg, f.ecuSerialMap[ecu], _stResetCmd)
		_, err := openCmd.CombinedOutput()
		if err != nil {
			errMsg := fmt.Sprintf("connect to stm32 with STLink ID %s", f.ecuSerialMap[ecu])
			return errors.Wrap(err, errMsg)
		}

		// If successful, set the target as active
		f.boardActive = true
		f.currentBoardId = f.ecuSerialMap[ecu]

		return nil
	}, bo)

	if err != nil {
		return errors.Wrap(err, "failed to connect after retries")
	}

	f.l.Info("target found", zap.String("target", ecu.String()), zap.String("target serial", f.currentBoardId))

	return nil
}

// Flash uses the stlink driver to flash the target with a provided binary
func (f *Flasher) Flash(binName string) error {
	if f.boardActive {
		f.l.Info("attempting to flash")

		flashCmd := exec.Command(_stFlashCmd, _stSerialArg, f.currentBoardId, _stResetArg, _stWriteArg, binName, _stWriteAddr)

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
	f.l.Info("disconnecting from target", zap.String("target id", f.currentBoardId))

	f.currentBoardId = ""
	f.boardActive = false

	return nil
}

// PowerCycleStm requires uhubctl installed and usb permissions set to not require sudo
// https://github.com/mvp/uhubctl
func (f *Flasher) PowerCycleStm(ecu fwutils.Ecu) error {
	f.l.Info("finding usb port number for target")

	stmPort, err := f.extractPortNumber(ecu)
	if err != nil {
		return errors.Wrap(err, "get usb port number for target")
	}

	f.l.Info("port found, power cycling target")

	offCmd := exec.Command(_uhubCmd, _uhubPwrArg, _uhubOff, _uhubPortNumArg, stmPort, _uhubHubArg, _uhubDefaultHub)
	_, err = offCmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "power off target")
	}

	onCmd := exec.Command(_uhubCmd, _uhubPwrArg, _uhubOn, _uhubPortNumArg, stmPort, _uhubHubArg, _uhubDefaultHub)
	_, err = onCmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "power off target")
	}

	f.l.Info("target has been restarted")

	return nil
}

func (f *Flasher) extractPortNumber(ecu fwutils.Ecu) (string, error) {
	cmd := exec.Command("uhubctl")

	output, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "error executing uhubctl command")
	}

	outputStr := string(output)

	lines := strings.Split(outputStr, "\n")

	// Iterate through each line
	for _, line := range lines {
		// Check if the line contains the serial number
		if strings.Contains(line, f.ecuSerialMap[ecu]) {
			// Extract the port number from the line
			parts := strings.Fields(line)
			portNumber := strings.TrimPrefix(parts[1], "Port")
			return portNumber, nil
		}
	}

	return "", errors.New("usb port number for target not found")
}
