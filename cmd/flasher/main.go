package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/macformula/hil/fwutils"
	"github.com/macformula/hil/fwutils/stflash"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	ecuToFlashStr string
	binaryPath    string
)

var _ecuToStlinkSerialNumber = map[fwutils.Ecu]string{
	fwutils.FrontController: "000F00205632500A20313236",
	fwutils.LvController:    "0006002A5632500920313236",
}

func main() {
	flag.StringVar(&ecuToFlashStr, "ecu", "", "ECU to flash (FrontController, LvController)")
	flag.StringVar(&binaryPath, "binary", "", "Path to the binary file to flash")
	flag.Parse()

	if ecuToFlashStr == "" || binaryPath == "" {
		fmt.Println("Missing required flags: --ecu and --binary")
		flag.PrintDefaults()
		return
	}

	ecuToFlash, err := fwutils.EcuString(ecuToFlashStr)
	if err != nil {
		fmt.Println("Invalid Ecu provided, options are:")
		for _, ecu := range fwutils.EcuValues() {
			fmt.Printf("\t%v\n", ecu.String())
		}
		return
	}

	// Check if binary path exists
	_, err = os.Stat(binaryPath)
	if os.IsNotExist(err) {
		fmt.Println("Binary file does not exist")
		return
	}

	// Check if binary path ends with .bin
	if !strings.HasSuffix(binaryPath, ".bin") {
		fmt.Println("Binary file should have .bin extension")
		return
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"stdout"}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	flasher := stflash.NewFlasher(logger, _ecuToStlinkSerialNumber)

	err = flasher.Connect(ecuToFlash)
	if err != nil {
		logger.Error("failed to connect",
			zap.String("ecu", ecuToFlash.String()),
			zap.Error(errors.Wrap(err, "connect")),
		)
	}

	err = flasher.PowerCycleStm(ecuToFlash)
	if err != nil {
		logger.Error("failed to power cycle stm",
			zap.String("ecu", ecuToFlash.String()),
			zap.Error(errors.Wrap(err, "power cycle stm")),
		)
	}

	err = flasher.Flash(binaryPath)
	if err != nil {
		panic(errors.Wrap(err, "flash"))
	}

	err = flasher.Disconnect()
	if err != nil {
		panic(errors.Wrap(err, "disconnect"))
	}
}
