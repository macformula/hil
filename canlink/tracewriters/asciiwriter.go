package tracewriters

import (
	"os"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"github.com/pkg/errors"
)

// Responsible for managing the ascii trace file
type AsciiWriter struct {
	traceFile *os.File
	l *zap.Logger
}

func NewAsciiWriter(l *zap.Logger) *AsciiWriter {
	return &AsciiWriter{
		l: l,
		traceFile: nil,
	}
}

func (w *AsciiWriter) CreateTraceFile(traceDir string, busName string) error {
	file, err := createEmptyTraceFile(traceDir, busName, "asc")
	w.traceFile = file
	if err != nil {
		w.l.Info("cannot create trace file")
		return errors.Wrap(err, "creating trace file")
	}

	return nil
}

func (w *AsciiWriter) WriteFrameToFile(frame *TimestampedFrame) error {
	line := w.convertToAscii(frame)

    _, err := w.traceFile.WriteString(line)
    if err != nil {
		w.l.Info("cannot write to file")
        return errors.Wrap(err, "writing to trace file")
    }

	_, err = w.traceFile.WriteString("\n")
    if err != nil {
		w.l.Info("cannot write to file")
        return errors.Wrap(err, "writing to trace file")
    }

	return nil
}

// Converts timestamped frames into strings for file writing 
func (w *AsciiWriter) convertToAscii(timestampedFrame *TimestampedFrame) string {
	var builder strings.Builder

	_, err := builder.WriteString(timestampedFrame.Time.Format(_messageTimeFormat))
	if err != nil {
		w.l.Error(err.Error())
	}

	_, err = builder.WriteString(" " + strconv.FormatUint(uint64(timestampedFrame.Frame.ID), _decimal))
	if err != nil {
		w.l.Error(err.Error())
	}

	_, err = builder.WriteString(" Rx")
	if err != nil {
		w.l.Error(err.Error())
	}

	_, err = builder.WriteString(" " + strconv.FormatUint(uint64(timestampedFrame.Frame.Length), _decimal))
	if err != nil {
		w.l.Error(err.Error())
	}

	for i := uint8(0); i < timestampedFrame.Frame.Length; i++ {
		builder.WriteString(" " + fmt.Sprintf("%02X", timestampedFrame.Frame.Data[i]))
		if err != nil {
			w.l.Error(err.Error())
		}
	}

	return builder.String()
}

func (w *AsciiWriter) CloseTraceFile() error {
	err := w.traceFile.Close()
	if err != nil {
		w.l.Error(err.Error())
		return errors.Wrap(err, "closing trace file")
	}

	return nil
}