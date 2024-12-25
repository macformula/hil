package canlink

import (
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// Ascii object provides utilities for writing frames to trace files in ascii format
type Ascii struct {
	fileExtention string
}

func NewAscii() *Ascii {
	return &Ascii{
		fileExtention: ".asc",
	}
}

func (a *Ascii) GetFileExtention() string {
	return a.fileExtention
}


// FrameToString converts a timestamped frame into a string, for file writing
func (a *Ascii) FrameToString(l *zap.Logger, timestampedFrame *TimestampedFrame) string {
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