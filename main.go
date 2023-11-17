package main

import (
	"context"
	"github.com/macformula/hil/canlink"
	"go.uber.org/zap"
	"time"
)

func main() {
	testDirectory := "actual_log" // Update with the actual path
	canInterface := "vcan0"       // Update with the appropriate CAN interface

	// Create a logger for the tracer
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Create a tracer instance with a timeout of 5 second
	tracer := canlink.NewTracer(
		canInterface,
		testDirectory,
		logger,
		canlink.WithTimeout(5*time.Second),
	)

	// Open the tracer
	tracer.Open(context.Background())

	// Start the tracer
	tracer.StartTrace(context.Background())
	time.Sleep(5*time.Second)
	// tracer.StopTrace()

	// tracer.Close()
}
