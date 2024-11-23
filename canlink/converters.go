package canlink

import (
	"fmt"
	"strconv"
	"strings"
	"encoding/json"

	"go.uber.org/zap"
)

const (
	_decimal = 10
)

// Should convert timestamped frames into desired file output as a string
type ConvertToString func(*zap.Logger, *TimestampedFrame) string

// Converts timestamped frames into strings for file writing 
func ConvertToAscii(l *zap.Logger, timestampedFrame *TimestampedFrame) string {
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
func ConvertToJson(l *zap.Logger, timestampedFrame *TimestampedFrame) string {
	type Json struct {
		Time  string `json:"time"`
		Id string `json:id`
		Length string `json:"length"`
	}

	
}