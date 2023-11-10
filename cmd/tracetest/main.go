package main

import (
	"context"
	"github.com/pkg/errors"
	"time"

	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/config"
	"go.uber.org/zap"
)

const (
	_configPath = "/opt/macfe/etc/config.yml"
)

func main() {
	config, err := config.NewConfig(_configPath)
	if err != nil {
		panic("error reading config file")
	}

	ctx := context.Background()

	cfg := zap.NewDevelopmentConfig()
	logger, err := cfg.Build()

	tracer := canlink.NewTracer(
		config.CanInterface,
		config.TracerDirectory,
		logger,
		canlink.WithBusName(config.BusName),
		canlink.WithTimeout(3*time.Second))

	err = tracer.Open(ctx)
	if err != nil {
		logger.Error("open tracer", zap.Error(err))
		return
	}

	err = tracer.StartTrace(ctx)
	if err != nil {
		logger.Error("start trace", zap.Error(err))
	}

	time.Sleep(5 * time.Second)

	err = tracer.StopTrace()
	if err != nil {
		logger.Error("stop trace", zap.Error(err))
	}

	time.Sleep(2 * time.Second)

	err = tracer.StartTrace(ctx)
	if err != nil {
		logger.Error("start trace", zap.Error(err))
	}

	time.Sleep(2 * time.Second)

	err = tracer.StopTrace()
	if err != nil {
		logger.Error("stop trace", zap.Error(err))
	}

	err = tracer.Close()
	if err != nil {
		logger.Error("close tracer", zap.Error(err))
	}

	if tracer.Error() != "" {
		logger.Error("tracer error", zap.Error(errors.New(tracer.Error())))
	}
}
