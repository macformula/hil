package canlink

import (
	"strconv"

	"encoding/json"
	"go.uber.org/zap"
)

// Jsonl object provides utilities for writing frames to trace files in jsonl format.
// See https://jsonlines.org/ to view Jsonl documentation.
type Jsonl struct{}

// GetFileExtension returns the file extension
func (a *Jsonl) GetFileExtension() string {
	return "jsonl"
}

// FrameToString converts a timestamped frame into a string, for file writing
func (a *Jsonl) FrameToString(l *zap.Logger, timestampedFrame *TimestampedFrame) string {
	jsonlObject := map[string]interface{}{
		"time":        timestampedFrame.Time.Format(_messageTimeFormat),
		"id":          strconv.FormatUint(uint64(timestampedFrame.Frame.ID), _decimal),
		"frameLength": strconv.FormatUint(uint64(timestampedFrame.Frame.Length), _decimal),
		"bytes":       timestampedFrame.Frame.Data,
	}

	jsonlData, err := json.Marshal(jsonlObject)
	if err != nil {
		l.Error(err.Error())
	}

	return string(jsonlData)
}
