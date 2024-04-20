package main

import (
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

	flasher := stflash.NewFlasher(*logger)

	err = flasher.Connect("303636454646353035323737353034383637313432313134")
	if err != nil {
		panic(errors.Wrap(err, "open flasher"))
	}

	//builder := fwutils.NewBuilder()

	//err = flasher.Flash("/opt/macfe/bin/PRINTF_TEST.bin")
	//if err != nil {
	//	panic(errors.Wrap(err, "flash stm32"))
	//}

}
