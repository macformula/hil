package canlink

import (
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_asciiFileSuffix     = "asc"
	_asciiFileLoggerName = "ascii_file"
	_decimal             = 10
	_hex                 = 16
)

// Ascii stores the information required to create each ASCII file.
type Ascii struct {
	l    *zap.Logger
	file *os.File
}

// NewAscii returns a new Ascii with its own ASCII file
func NewAscii(l *zap.Logger) *Ascii {
	return &Ascii{
		l: l.Named(_asciiFileLoggerName),
	}
}

// dumpToFile takes a CAN frame and writes it to an ASCII file
func (a *Ascii) dumpToFile(frames []TimestampedFrame, traceDir, busName string) error {
	a.l.Info("dumping data to ascii file")

	f, err := createTraceFile(traceDir, busName, _asciiFileSuffix)
	if err != nil {
		a.l.Error("failed to create trace file",
			zap.String("trace_dir", traceDir),
			zap.Error(err),
		)

		return errors.Wrap(err, "create trace file")
	}

	for _, frame := range frames {
		_, err = f.WriteString(a.formatFrame(&frame) + "\n")
		if err != nil {
			return errors.Wrap(err, "write string")
		}
	}

	return nil
}

// parseFrame concatenates the frame components in a standardized format
func (a *Ascii) formatFrame(timestampedFrame *TimestampedFrame) string {
	var builder strings.Builder

	_, err := builder.WriteString(timestampedFrame.Time.Format(_messageTimeFormat))
	if err != nil {
		a.l.Error(err.Error())
	}

	_, err = builder.WriteString(" " + strconv.FormatUint(uint64(timestampedFrame.Frame.ID), _decimal))
	if err != nil {
		a.l.Error(err.Error())
	}

	_, err = builder.WriteString(" Rx")
	if err != nil {
		a.l.Error(err.Error())
	}

	_, err = builder.WriteString(" " + strconv.FormatUint(uint64(timestampedFrame.Frame.Length), _decimal))
	if err != nil {
		a.l.Error(err.Error())
	}

	for i := uint8(0); i < timestampedFrame.Frame.Length; i++ {
		builder.WriteString(" " + strconv.FormatUint(uint64(timestampedFrame.Frame.Data[i]), _hex))
		if err != nil {
			a.l.Error(err.Error())
		}
	}

	return builder.String()
}
