package canlink

import (
	"strconv"

	"go.uber.org/zap"
	"encoding/json"
)

// Json object provides utilities for writing frames to trace files in json format
type Json struct {
	fileExtention string
}

func NewJson() *Json {
	return &Json{
		fileExtention: ".jsonl",
	}
}

func (a *Json) GetFileExtention() string {
	return a.fileExtention
}

// FrameToString converts a timestamped frame into a string, for file writing
func (a *Json) FrameToString(l *zap.Logger, timestampedFrame *TimestampedFrame) string {
	jsonObject := map[string]interface{}{
		"time":        timestampedFrame.Time.Format(_messageTimeFormat),
		"id":          strconv.FormatUint(uint64(timestampedFrame.Frame.ID), _decimal),
		"frameLength": strconv.FormatUint(uint64(timestampedFrame.Frame.Length), _decimal),
		"bytes":       timestampedFrame.Frame.Data,
	}

	jsonData, err := json.Marshal(jsonObject)
	if err != nil {
		l.Error(err.Error())
	}

	return string(jsonData)
}