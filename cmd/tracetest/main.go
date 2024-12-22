package main

import (
	"context"
	"fmt"
	"time"

	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"

	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/canlink/writer"
	"github.com/macformula/hil/macformula/cangen/vehcan"
	"github.com/pkg/errors"
)

const (
	// Can bus select
	_busName  = "veh"
	_canIface = "vcan0"

	// Env config
	_timeFormat        = "2006-01-02_15-04-05"
	_logFilenameFormat = ".logs.asc"
	_traceDir          = "./traces"
	_logLevel          = zap.DebugLevel
)

func main() {
	ctx := context.Background()

	logFileName := _logFilenameFormat

	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.OutputPaths = []string{logFileName}
	loggerConfig.Level = zap.NewAtomicLevelAt(_logLevel)
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(errors.Wrap(err, "failed to create logger"))
	}

	conn, err := socketcan.DialContext(context.Background(), "can", _canIface)
	if err != nil {
		logger.Error("failed to create socket can connection",
			zap.String("can_interface", _canIface),
			zap.Error(err),
		)

		return
	}

	canClient := canlink.NewCanClient(vehcan.Messages(), conn, logger)

	writers := make([]tracewriters.TraceWriter, 0)
	writers = append(writers, tracewriters.NewAsciiWriter(logger))
	writers = append(writers, tracewriters.NewJsonWriter(logger))

	tracer := canlink.NewTracer(
		_canIface,
		_traceDir,
		logger,
		conn,
		canlink.WithBusName(_busName),
		canlink.WithTraceWriters(writers))

	err = canClient.Open()
	if err != nil {
		logger.Error("open can client", zap.Error(err))

		return
	}

	err = tracer.Open(ctx)
	if err != nil {
		logger.Error("open tracer", zap.Error(err))
		return
	}

	err = tracer.StartTrace(ctx)
	if err != nil {
		logger.Error("start trace", zap.Error(err))
	}

	fmt.Println("-------------- Starting Test --------------")
	fmt.Println("-------------- CTRL-C to Stop -------------")
	time.Sleep(6 * time.Second)

	fmt.Println("-------------- Test Complete --------------")

	logger.Info("closing trace test")

	err = tracer.StopTrace()
	if err != nil {
		logger.Error("stop trace", zap.Error(err))
	}

	err = tracer.Close()
	if err != nil {
		logger.Error("close tracer", zap.Error(err))
	}

	if tracer.Error() != nil {
		logger.Error("tracer error", zap.Error(tracer.Error()))
	}
}

