package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/macformula/hil/iocontrol"
	"github.com/macformula/hil/iocontrol/raspi"
	"github.com/macformula/hil/iocontrol/sil"
	"github.com/macformula/hil/iocontrol/speedgoat"
	"github.com/macformula/hil/macformula"
	"github.com/macformula/hil/macformula/pinout"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_logFileName         = "iocontrol.log"
	_speedgoatAddr       = "192.168.10.1:8001"
	_speedgoatAddrFromPi = "slrt@192.168.10.1"
	_speedgoatPw         = "slrt"
	_speedgoatScriptPath = "./sgCtrl.sh"
	_speedgoatModelName  = "speedgoat_tcp_io"
	_silPort             = 31522
)

var (
	revisionStr  = flag.String("revision", "ev5", "Revision of the system")
	useSpeedgoat = flag.Bool("use-speedgoat", false, "Use Speedgoat controller")
	useRaspi     = flag.Bool("use-raspi", false, "Use Raspi controller")
	useSil       = flag.Bool("use-sil", false, "Use Sil controller")
)

func main() {
	ctx := context.Background()

	flag.Parse()

	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{_logFileName, "stdout"}
	logger, err := cfg.Build()
	if err != nil {
		panic(errors.Wrap(err, "zap config build"))
	}
	defer logger.Sync()

	revision, err := pinout.RevisionString(*revisionStr)
	if err != nil {
		fmt.Printf("Invalid revision (%s) valid options (%v)", *revisionStr, pinout.RevisionStrings())
		return
	}

	logger.Info("starting iocheckout", zap.String("revision", revision.String()))

	var ioControlOpts []iocontrol.IOControlOption

	if *useSpeedgoat {
		sg := speedgoat.NewController(logger, _speedgoatAddr, speedgoat.WithModelAutoload(_speedgoatScriptPath, _speedgoatPw, _speedgoatAddrFromPi, _speedgoatModelName))
		ioControlOpts = append(ioControlOpts, iocontrol.WithSpeedgoat(sg))
	}

	if *useRaspi {
		rp := raspi.NewController()
		ioControlOpts = append(ioControlOpts, iocontrol.WithRaspi(rp))
	}

	if *useSil {
		s := sil.NewController(_silPort, logger)
		ioControlOpts = append(ioControlOpts, iocontrol.WithSil(s))
	}

	ioControl := iocontrol.NewIOControl(logger, ioControlOpts...)

	iocheckout := macformula.NewIoCheckout(revision, ioControl, logger)

	logger.Info("opening iocheckout")

	err = iocheckout.Open(ctx)
	if err != nil {
		logger.Error("failed to open iocheckout", zap.Error(errors.Wrap(err, "open iocheckout")))

		return
	}

	logger.Info("starting iocheckout")

	defer func(io *macformula.IoCheckout, l *zap.Logger) {
		l.Info("closing iocheckout")

		err = io.Close()
		if err != nil {
			l.Error("failed to close iocheckout", zap.Error(err))
		}
	}(iocheckout, logger)

	err = iocheckout.Start()
	if err != nil {
		logger.Error("failed iocheckout start", zap.Error(err))

		return
	}
}
