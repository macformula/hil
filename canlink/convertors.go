package canlink

import (
	"go.uber.org/zap"
)

// Converter provides functionality for converting timestamped frames into strings for file writing. 
// Each supported trace file type must impliment Converter.
type Converter interface {
	GetFileExtention() string
	FrameToString(*zap.Logger, *TimestampedFrame) string
}
