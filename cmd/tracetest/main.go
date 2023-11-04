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
		logger,
		config.TracerDirectory,
		config.CANInterface,
		can.WithBusName(config.BusName),
		can.WithTimeout(3*time.Second))

	tracer.StartTrace(ctx)

	time.Sleep(5 * time.Second)

	err = tracer.StopTrace()
	if err != nil {
		logger.Error(fmt.Sprintf("error with tracer: %e", err))
	}
}
