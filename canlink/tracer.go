package canlink

import (
	"fmt"
	"os"
	"time"

	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

const (
	_defaultTimeout    = 30 * time.Minute
	_frameBufferLength = 10
	_loggerName        = "can_tracer"

	_defaultFileName = "" // if no file name is provided, file name will be generated when trace file is created

	_decimal = 10

	_messageTimeFormat  = "15:04:05.0000"
	_filenameTimeFormat = "15-04-05"
	_filenameDateFormat = "2006-01-02"
)

// TracerOption is a type for functions operating on Tracer
type TracerOption func(*Tracer)

// Tracer listens on a CAN bus and records all traffic
type Tracer struct {
	l   *zap.Logger
	err *utils.ResettableError

	converter Converter
	traceFile *os.File
	fileName  string

	canInterface string
	timeout      time.Duration
}

// NewTracer returns a new Tracer
func NewTracer(
	canInterface string,
	l *zap.Logger,
	converter Converter,
	opts ...TracerOption) *Tracer {

	tracer := &Tracer{
		l:            l.Named(_loggerName),
		err:          utils.NewResettaleError(),
		timeout:      _defaultTimeout,
		canInterface: canInterface,
		fileName:     _defaultFileName,
		converter:    converter,
	}

	for _, o := range opts {
		o(tracer)
	}

	return tracer
}

// WithTimeout sets the timeout for the Tracer
func WithTimeout(timeout time.Duration) TracerOption {
	return func(t *Tracer) {
		t.timeout = timeout
	}
}

// WithFileName sets the file name
func WithFileName(fileName string) TracerOption {
	return func(t *Tracer) {
		t.fileName = fileName
	}
}

// Error returns the error set during trace execution
func (t *Tracer) Error() error {
	return t.err.Err()
}

// Handle listens to the frames in the broadcastChan and writes them to a file
func (t *Tracer) Handle(broadcastChan chan TimestampedFrame, stopChan chan struct{}) error {
	err := t.createTraceFile()
	if err != nil {
		return err
	}

	timeout := time.After(t.timeout)

	func() error {
		for {
			select {
			case <-stopChan:
				t.l.Info("stopping handle")
				t.close()
				return nil
			case <-timeout:
				t.l.Info("maximum trace time reached")
				return nil
			case receivedFrame := <-broadcastChan:
				t.l.Info("frame received")
				line := t.converter.FrameToString(t.l, &receivedFrame)

				_, err := t.traceFile.WriteString(line + "\n")
				if err != nil {
					t.l.Info("cannot write to file")
					return errors.Wrap(err, "writing to trace file")
				}
			default:
			}
		}
	}()

	return nil
}

// Name returns the name of the handler.
// This value is only used for error logging
func (t *Tracer) Name() string {
	return "Tracer"
}

// GetFileName simply returns the file name of the trace file this tracer is responsible for
func (t *Tracer) GetFileName() string {
	return t.fileName
}

// close closes the trace file
func (t *Tracer) close() error {
	t.l.Info("closing trace file")
	err := t.traceFile.Close()
	if err != nil {
		t.l.Error(err.Error())
		return errors.Wrap(err, "closing trace file")
	}

	return nil
}

// createEmptyTraceFile generates empty trace file given a file name
func (t *Tracer) createEmptyTraceFile() (*os.File, error) {
	file, err := os.Create(t.fileName)
	if err != nil {
		return nil, errors.Wrap(err, "create trace file")
	}

	return file, nil
}

// createTraceFile creates a new trace file with the proper file name
func (t *Tracer) createTraceFile() error {
	if t.fileName == _defaultFileName {
		dateStr := time.Now().Format(_filenameDateFormat)
		timeStr := time.Now().Format(_filenameTimeFormat)
		t.fileName = fmt.Sprintf(
			"%s_%s.%s",
			dateStr,
			timeStr,
			t.converter.GetFileExtension(),
		)
	} else {
		t.fileName = fmt.Sprintf(
			"%s.%s",
			t.fileName,
			t.converter.GetFileExtension(),
		)
	}

	file, err := t.createEmptyTraceFile()
	t.traceFile = file

	if err != nil {
		t.l.Info("cannot create trace file")
		return errors.Wrap(err, "creating trace file")
	}

	return nil
}
