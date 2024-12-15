package tracewriters

import (
	"os"
	"strconv"

	"encoding/json"

	"go.uber.org/zap"
	"github.com/pkg/errors"
)

// Responsible for managing the json trace file
type JsonWriter struct {
	traceFile *os.File
	l *zap.Logger
}

func NewJsonWriter(l *zap.Logger) *JsonWriter {
	return &JsonWriter{
		l: l,
		traceFile: nil,
	}
}

func (w *JsonWriter) CreateTraceFile(traceDir string, busName string) error {
	file, err := createEmptyTraceFile(traceDir, busName, "json")

	w.traceFile = file
	if err != nil {
		w.l.Info("cannot create trace file")
		return errors.Wrap(err, "creating trace file")
	}
	
	// Initializes the trace file with a "[" to ensure the json array is started
	_, err = file.WriteString("[\n")
	if err != nil {
		w.l.Info("cannot write to file")
		return errors.Wrap(err, "initialiing json file")
	}

	return nil
}

func (w *JsonWriter) WriteFrameToFile(frame *TimestampedFrame) error {
	line := w.convertToJson(frame)

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

func (w *JsonWriter) convertToJson(timestampedFrame *TimestampedFrame) string {
	jsonObject := map[string]interface{}{
		"time":   timestampedFrame.Time.Format(_messageTimeFormat),
		"id":    strconv.FormatUint(uint64(timestampedFrame.Frame.ID), _decimal),
		"frameLength": strconv.FormatUint(uint64(timestampedFrame.Frame.Length), _decimal),
		"bytes": timestampedFrame.Frame.Data,
	}

	jsonData, err := json.Marshal(jsonObject)
	if err != nil {
		w.l.Error(err.Error())
	}

	return string(jsonData) + ","
}

func (w *JsonWriter) CloseTraceFile() error {
	// this is nessesary to prevent trailing commas in the json
	info, _ := w.traceFile.Stat()
	fileLength := info.Size()

	closingBracket := []byte("\n]")

	// this checks if no messages have been written into the file
	if fileLength <= 2 {
		w.traceFile.WriteString(string(closingBracket))
		return nil
	}

	// otherwise we need to remove the trailing comma in the file
	w.traceFile.WriteAt(closingBracket, fileLength - 2)


	err := w.traceFile.Close()
	if err != nil {
		w.l.Error(err.Error())
		return errors.Wrap(err, "closing trace file")
	}

	return nil
}