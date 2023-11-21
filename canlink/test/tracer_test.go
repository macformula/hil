package test

import (
	"context"
	"github.com/macformula/hil/canlink"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"testing"
	"time"
)

const (
	_canInterface = "vcan0"
	_busName      = "vcan0"
	_timeout      = 5 * time.Second
	_defaultWait
)

func TestTracer(t *testing.T) {
	// get a background context
	ctx := context.Background()

	// set up a basic logger for duration of testing
	cfg := zap.NewDevelopmentConfig()
	logger, err := cfg.Build()

	// get the current directory of the test
	directory, err := os.Getwd()
	if err != nil {
		panic("get working directory")
	}

	// Create a tracer instance with a timeout of 5 second
	tracer := canlink.NewTracer(
		_canInterface,
		directory,
		logger,
		canlink.WithBusName(_busName),
		canlink.WithTimeout(_timeout),
	)

	// Open the tracer
	err = tracer.Open(ctx)
	assert.NoError(t, err)

	// Start the tracer
	err = tracer.StartTrace(ctx)
	assert.NoError(t, err)

	// Run sample can traffic to be picked up by the tracer
	go func() {
		cmd := exec.Command("canplayer", "vcan0=vcan0", "-I", "sample_output_timestamped.log")
		err := cmd.Run()
		assert.NoError(t, err)
	}()

	time.Sleep(5 * time.Second)

	// Start the tracer
	err = tracer.StopTrace()
	assert.NoError(t, err)

	// Close the tracer
	err = tracer.Close()
	assert.NoError(t, err)

}
