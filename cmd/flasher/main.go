package main

import (
	"github.com/macformula/hil/fwutils"
	"github.com/macformula/hil/fwutils/stflash"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"stdout"}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	ecuMap := map[fwutils.Ecu]string{
		fwutils.FrontController: "000F00205632500A20313236",
		fwutils.LvController:    "0006002A5632500920313236",
	}

	flasher := stflash.NewFlasher(*logger, ecuMap)

	err = flasher.Connect(fwutils.FrontController)
	if err != nil {
		panic(errors.Wrap(err, "open flasher"))
	}

	err = flasher.PowerCycleStm(fwutils.FrontController)
	if err != nil {
		panic(errors.Wrap(err, "power cycle"))
	}

	err = flasher.Disconnect()

	err = flasher.Connect(fwutils.LvController)
	if err != nil {
		panic(errors.Wrap(err, "open flasher"))
	}

	err = flasher.PowerCycleStm(fwutils.LvController)
	if err != nil {
		panic(errors.Wrap(err, "power cycle"))
	}

	err = flasher.Disconnect()

	//err = flasher.Flash("/opt/macfe/bin/PRINTF_TEST.bin")
	//if err != nil {
	//	panic(errors.Wrap(err, "flash stm32"))
	//}
}
