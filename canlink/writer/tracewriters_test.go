package writer

import (
	"testing"
	"time"

	"go.einride.tech/can"
	"go.uber.org/zap"
)

// Tests a if a simple timestamped frame is converted to json with proper formatting
func TestJsonConvertToString(t *testing.T) {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	logger, err := loggerConfig.Build()

	if err != nil {
		t.Errorf("failed to build logger")
	}

	writer := NewWriter(
		logger,
		".jsonl",
	)

	time := time.Date(2024, time.November, 27, 10, 30, 45, 0, time.Local)
	timestampedFrame := TimestampedFrame{
		Frame: can.Frame{
			ID:     12,
			Length: 8,
			Data:   [8]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88},
		},
		Time: time,
	}

	output := writer.convertToString(logger, &timestampedFrame)
	expected := `{"bytes":[17,34,51,68,85,102,119,136],"frameLength":"8","id":"12","time":"10:30:45.0000"}`

	if output != expected {
		t.Fatalf(`
		writer.convertToJson(canlink.TimestampedFrame{
			Frame: can.Frame{
				ID: 12,
				Length: 8,
				Data: [8]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88},
			},
			Time: time,
		}) = 
		 %v want 
		 %v`, output, expected)
	}
}

// Tests a if a simple timestamped frame is converted to ascii with proper format
func TestAsciiConvertToString(t *testing.T) {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	logger, err := loggerConfig.Build()

	if err != nil {
		t.Errorf("failed to build logger")
	}

	writer := NewWriter(
		logger,
		".asc",
	)

	time := time.Date(2024, time.November, 27, 10, 30, 45, 0, time.Local)
	timestampedFrame := TimestampedFrame{
		Frame: can.Frame{
			ID:     12,
			Length: 8,
			Data:   [8]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88},
		},
		Time: time,
	}

	output := writer.convertToString(logger, &timestampedFrame)
	expected := `10:30:45.0000 12 Rx 8 11 22 33 44 55 66 77 88`

	if output != expected {
		t.Fatalf(`
		writer.convertToJson(canlink.TimestampedFrame{
			Frame: can.Frame{
				ID: 12,
				Length: 8,
				Data: [8]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88},
			},
			Time: time,
		}) = 
		 %v; want 
		 %v`, output, expected)
	}
}
