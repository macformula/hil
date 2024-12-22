package writer

import (
	"strings"
	"fmt"
	"strconv"

	"encoding/json"
	"go.uber.org/zap"
)

// Converts timestamped frames into strings for file writing 
func convertToAscii(l *zap.Logger, timestampedFrame *TimestampedFrame) string {
	var builder strings.Builder

	_, err := builder.WriteString(timestampedFrame.Time.Format(_messageTimeFormat))
	if err != nil {
		l.Error(err.Error())
	}

	_, err = builder.WriteString(" " + strconv.FormatUint(uint64(timestampedFrame.Frame.ID), _decimal))
	if err != nil {
		l.Error(err.Error())
	}

	_, err = builder.WriteString(" Rx")
	if err != nil {
		l.Error(err.Error())
	}

	_, err = builder.WriteString(" " + strconv.FormatUint(uint64(timestampedFrame.Frame.Length), _decimal))
	if err != nil {
		l.Error(err.Error())
	}

	for i := uint8(0); i < timestampedFrame.Frame.Length; i++ {
		builder.WriteString(" " + fmt.Sprintf("%02X", timestampedFrame.Frame.Data[i]))
		if err != nil {
			l.Error(err.Error())
		}
	}

	return builder.String()
}

// Converts timestamped frames into strings for file writing 
func convertToJson(l *zap.Logger, timestampedFrame *TimestampedFrame) string {
	jsonObject := map[string]interface{}{
		"time":   timestampedFrame.Time.Format(_messageTimeFormat),
		"id":    strconv.FormatUint(uint64(timestampedFrame.Frame.ID), _decimal),
		"frameLength": strconv.FormatUint(uint64(timestampedFrame.Frame.Length), _decimal),
		"bytes": timestampedFrame.Frame.Data,
	}

	jsonData, err := json.Marshal(jsonObject)
	if err != nil {
		l.Error(err.Error())
	}

	return string(jsonData) 
}