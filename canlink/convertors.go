package canlink

import (
	"go.uber.org/zap"
)

// Converter provides functionality for converting timestamped frames into strings for file writing.
// Each supported trace file type must implement Converter.
type Converter interface {
	GetFileExtension() string
	FrameToString(*zap.Logger, *TimestampedFrame) string
}
