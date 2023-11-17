package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/macformula/hil/canlink"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func main() {
	testDirectory := "./actual_log" // Update with the actual path
	canInterface := "vcan0"         // Update with the appropriate CAN interface

	fmt.Println("Starting program")
	rawJSON, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	cfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder

	logger := zap.Must(cfg.Build())

	// Create a logger for the tracer
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
	time.Sleep(5 * time.Second)

	tracer.StopTrace()

	// tracer.Close()

}
