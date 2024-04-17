package main

import (
	"github.com/macformula/hil/flash/stflash"
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

	err = flasher.Connect()
	if err != nil {
		panic(errors.Wrap(err, "open flasher"))
	}

	//err = flasher.Flash("/opt/macfe/bin/PRINTF_TEST.bin")
	//if err != nil {
	//	panic(errors.Wrap(err, "flash stm32"))
	//}

}
