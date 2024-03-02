package main

import (
	"context"
	"fmt"
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
	formattedTime := time.Now().Format("2006.01.02_15.04.05")
	fileName := fmt.Sprintf("/opt/macfe/traces/logs/file_%s.log", formattedTime)
	cfg.OutputPaths = []string{fileName}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	tracer := canlink.NewTracer(
		config.CanInterface,
		config.TracerDirectory,
		logger,
		make([]canlink.FileType, 3, 3), // NEEDA FIX THIS
		canlink.WithBusName(config.BusName),
		canlink.WithTimeout(3*time.Second),
		canlink.WithAscii(),
		canlink.WithMcap(),
		canlink.WithCSV())

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

	if tracer.Error() != nil {
		logger.Error("tracer error", zap.Error(tracer.Error()))
	}
}
