package main

import (
	"context"
	"fmt"
	"github.com/macformula/hil/can"
	"github.com/macformula/hil/config"
	"go.uber.org/zap"
	"time"
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
	defer logger.Sync()

	tracer := can.NewTracer(
		config.CanInterface,
		config.TracerDirectory,
		logger,
		can.WithBusName(config.BusName),
		can.WithTimeout(3*time.Second))

	err = tracer.Open(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error with opening tracer %e", err))
	}

	err = tracer.StartTrace(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error with starting tracer %e", err))
	}

	time.Sleep(5 * time.Second)

	err = tracer.StopTrace()
	if err != nil {
		logger.Error(fmt.Sprintf("error with tracer: %e", err))
	}

	time.Sleep(5 * time.Second)

	err = tracer.StartTrace(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error with starting tracer %e", err))
	}

	time.Sleep(2 * time.Second)

	err = tracer.StopTrace()
	if err != nil {
		logger.Error(fmt.Sprintf("error with tracer: %e", err))
	}

	err = tracer.Close()
	if err != nil {
		logger.Error(fmt.Sprintf("error with closing tracer %e", err))
	}

	fmt.Printf("Error is: %v", tracer.Error())

}
