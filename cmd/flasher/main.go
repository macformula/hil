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
		fwutils.FrontController: "303636454646353035323737353034383637313432313134",
		fwutils.LvController:    "066EFF505277504867142114",
	}

	flasher := stflash.NewFlasher(*logger, ecuMap)

	err = flasher.PowerCycleStm(fwutils.LvController)
	if err != nil {
		panic(errors.Wrap(err, "power cycle"))
	}

	err = flasher.Connect(fwutils.FrontController)
	if err != nil {
		panic(errors.Wrap(err, "open flasher"))
	}

	err = flasher.Flash("/opt/macfe/bin/PRINTF_TEST.bin")
	if err != nil {
		panic(errors.Wrap(err, "flash stm32"))
	}

}
