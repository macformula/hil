package test

import (
	"context"
	"encoding/json"
	"github.com/macformula/hil/canlink"
	"go.uber.org/zap/zapcore"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	//"github.com/macformula/hil/canlink" // Update this import path based on your project structure
	"go.uber.org/zap"
)

func TestCanLinkIntegration(t *testing.T) {
	// Define the test directory and bus interface
	testDirectory := "actual_log" // Update with the actual path
	canInterface := "vcan0"       // Update with the appropriate CAN interface

	// Create a logger for the tracer
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
	defer logger.Sync()

	// Create a tracer instance with a timeout of 5 second
	tracer := canlink.NewTracer(
		canInterface,
		testDirectory,
		logger,
		canlink.WithTimeout(5*time.Second),
	)

	// Open the tracer
	err = tracer.Open(context.Background())
	assert.NoError(t, err)

	// Start the tracer
	err = tracer.StartTrace(context.Background())
	assert.NoError(t, err)

	// Run canplayer in a separate goroutine
	go func() {
		cmd := exec.Command("canplayer", "vcan0=vcan0", "-I", "sample_output_timestamped.log")
		err := cmd.Run()
		assert.NoError(t, err)
	}()

	// Wait for some time to allow the tracer to capture frames
	time.Sleep(3 * time.Second)

	// Stop the tracer
	err = tracer.StopTrace()
	assert.NoError(t, err)

	// Close the tracer
	err = tracer.Close()
	assert.NoError(t, err)

	// Compare the generated log file with the expected log file
	actualLogFile := testDirectory + "/" + canInterface + "_" + time.Now().Format(canlink._filenameTimeFormat) + "_" + time.Now().Format(canlink._filenameDateFormat)
	expectedLogFile := "expected_log.asc"

	actualContent, err := os.ReadFile(actualLogFile)
	assert.NoError(t, err)

	expectedContent, err := os.ReadFile(expectedLogFile)
	assert.NoError(t, err)

	assert.Equal(t, string(expectedContent), string(actualContent), "Log content mismatch")
}
