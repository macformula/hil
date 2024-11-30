package canlink

import (
	"context"
	"testing"
	"time"

	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"

	"github.com/macformula/hil/canlink/tracewriters"
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

func TestTracer(t *testing.T) {
	tracer, logger, teardown := setup(t)
	defer teardown(t, tracer, logger)

	time.Sleep(5 * time.Second)
	
}

func setup(t *testing.T) (*Tracer, *zap.Logger, func(*testing.T, *Tracer, *zap.Logger)) {
	ctx := context.Background()
	logFileName := _logFilenameFormat

	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.OutputPaths = []string{logFileName}
	loggerConfig.Level = zap.NewAtomicLevelAt(_logLevel)
	logger, err := loggerConfig.Build()
	if err != nil {
		t.Logf("Failed to create logger. Error: %v", err)
	}

	conn, err := socketcan.DialContext(context.Background(), "can", _canIface)
	if err != nil {
		t.Fatalf("Failed to create socket can connection. Error: %v", err)
	}

	writers := make([]tracewriters.TraceWriter, 0)
	writers = append(writers, tracewriters.NewAsciiWriter(logger))
	writers = append(writers, tracewriters.NewJsonWriter(logger))

	tracer := NewTracer(
		_canIface,
		_traceDir,
		logger,
		conn,
		WithBusName(_busName),
	 	WithTraceWriters(writers))

	err = tracer.Open(ctx)
	if err != nil {
		t.Fatalf("Error opening tracer. Error: %v", err)
	}
	err = tracer.StartTrace(ctx)
	if err != nil {
		t.Fatalf("Error starting trace. Error: %v", err)
	}

	teardown := func(t *testing.T, tracer *Tracer, logger *zap.Logger) {
		logger.Info("closing trace test")

		err = tracer.StopTrace()
		if err != nil {
			t.Logf("Failed to stop trace. Error: %v", err)
		}

		err = tracer.Close()
		if err != nil {
			t.Logf("Failed to close tracer. Error: %v", err)
		}

		if tracer.Error() != nil {
			t.Logf("Tracer returned with error set. Error: %v", tracer.Error().Error())
		}
	}

	return tracer, logger, teardown
}