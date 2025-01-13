package canlink

import (
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// Text object provides utilities for writing frames to trace files in text format
type Text struct{}

func (a *Text) GetFileExtension() string {
	return "txt"
}

// FrameToString converts a timestamped frame into a string, for file writing
func (a *Text) FrameToString(l *zap.Logger, timestampedFrame *TimestampedFrame) string {
	var builder strings.Builder

	write := func(s string) {
		_, err := builder.WriteString(s)
		if err != nil {
			l.Error(err.Error())
		}
	}

	write(timestampedFrame.Time.Format(_messageTimeFormat))
	write(" " + strconv.FormatUint(uint64(timestampedFrame.Frame.ID), _decimal))
	write(" Rx")
	write(" " + strconv.FormatUint(uint64(timestampedFrame.Frame.Length), _decimal))

	for i := uint8(0); i < timestampedFrame.Frame.Length; i++ {
		write(" " + fmt.Sprintf("%02X", timestampedFrame.Frame.Data[i]))
	}

	return builder.String()
}
