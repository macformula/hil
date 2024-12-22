package writer

import (
	"os"
	
	"go.uber.org/zap"
	"github.com/pkg/errors"
)

// Responsible for managing the trace file
type Writer struct {
	traceFile *os.File
	fileExtention string
	convertToString func(*zap.Logger, *TimestampedFrame) string
	l *zap.Logger
}

func NewWriter(l *zap.Logger, fileExtention string) *Writer {
	convertToStringMapping := map[string]func(*zap.Logger, *TimestampedFrame) string{
        ".jsonl": convertToJson,
        ".asc": convertToAscii,
    }
	return &Writer{
		l: l,
		fileExtention: fileExtention,
		convertToString: convertToStringMapping[fileExtention],
		traceFile: nil,
	}
}

func (w *Writer) CreateTraceFile(traceDir string, busName string) error {
	file, err := createEmptyTraceFile(traceDir, busName, w.fileExtention)
	w.traceFile = file
	if err != nil {
		w.l.Info("cannot create trace file")
		return errors.Wrap(err, "creating trace file")
	}

	return nil
}

func (w *Writer) WriteFrameToFile(frame *TimestampedFrame) error {
	line := w.convertToString(w.l, frame)

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


func (w *Writer) CloseTraceFile() error {
	err := w.traceFile.Close()
	if err != nil {
		w.l.Error(err.Error())
		return errors.Wrap(err, "closing trace file")
	}

	return nil
}