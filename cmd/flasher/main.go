package main

import (
	"github.com/macformula/hil/flash/stflash"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_loggerName = "main.log"
)

func main() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{_loggerName}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	flasher := stflash.NewFlasher(*logger)

	err = flasher.Open()
	if err != nil {
		panic(errors.Wrap(err, "open flasher"))
	}
}
