package main

import (
	"github.com/macformula/hil/iocontrol"
	"github.com/macformula/hil/iocontrol/raspi"
	"github.com/macformula/hil/iocontrol/speedgoat"
	"github.com/macformula/hil/macformula"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_logFileName   = "iocontrol.log"
	_revision      = macformula.Ev5
	_speedgoatAddr = "192.168.10.1:8001"
)

func main() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{_logFileName}
	logger, err := cfg.Build()
	if err != nil {
		panic(errors.Wrap(err, "zap config build"))
	}
	defer logger.Sync()

	rpiController := raspi.NewController()
	sgController := speedgoat.NewController(logger, _speedgoatAddr)

	ioControl := iocontrol.NewIOControl(logger, iocontrol.WithRaspi(rpiController), iocontrol.WithSpeedgoat(sgController))

	logger.Info("opening iocontrol")

	err = ioControl.Open()
	if err != nil {
		panic(errors.Wrap(err, "open iocontrol"))
	}

	ioCheckout := macformula.NewIoCheckout(_revision, ioControl, logger)

	logger.Info("starting iocheckout")

	err = ioCheckout.Start()
	if err != nil {
		logger.Error("iocheckout start", zap.Error(err))
	}

	return
}
