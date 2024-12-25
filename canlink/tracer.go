package canlink

import (
	"os"
	"time"
	"fmt"

	"path/filepath"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

const (
	_defaultTimeout    = 30 * time.Minute
	_frameBufferLength = 10
	_loggerName        = "can_tracer"

	_decimal = 10

	_messageTimeFormat = "15:04:05.0000"
	_filenameTimeFormat = "15-04-05"
	_filenameDateFormat = "2006-01-02"
)


// TracerOption is a type for functions operating on Tracer
type TracerOption func(*Tracer)

// Tracer listens on a CAN bus and records all traffic
type Tracer struct {
	l          *zap.Logger
	frameCh    chan TimestampedFrame // This channel will be provided by the bus manager
	err        *utils.ResettableError

	traceDir string
	converter Converter
	traceFile *os.File

	canInterface string
	timeout      time.Duration
	busName      string
}

// NewTracer returns a new Tracer
func NewTracer(
	canInterface string,
	traceDir string,
	l *zap.Logger,
	opts ...TracerOption) *Tracer {

	tracer := &Tracer{
		l:            l.Named(_loggerName),
		err:          utils.NewResettaleError(),
		timeout:      _defaultTimeout,
		canInterface: canInterface,
		traceDir:     traceDir,
		busName:      canInterface,
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

// WithBusName sets the name of the bus for the Tracer
func WithBusName(name string) TracerOption {
	return func(t *Tracer) {
		t.busName = name
	}
}

// WithFrameChannel sets frame channel for incoming frames
func WithFrameCh(frameCh chan TimestampedFrame) TracerOption {
	return func(t *Tracer) {
		t.frameCh = frameCh
	}
}

// WithConverter sets the converter
func WithConverter(converter Converter) TracerOption {
	return func(t *Tracer) {
		t.converter = converter
	}
}

// Close closes the trace file
func (t *Tracer) Close() error {
	t.l.Info("closing trace file")
	err := t.traceFile.Close()
	if err != nil {
		t.l.Error(err.Error())
		return errors.Wrap(err, "closing trace file")
	}

	if err != nil {
		return errors.Wrap(err, "close trace file")
	}

	return nil
}

// Error returns the error set during trace execution
func (t *Tracer) Error() error {
	return t.err.Err()
}

// Handle listens to the frames in the broadcastChan and writes them to a file
func (t *Tracer) Handle(broadcastChan chan TimestampedFrame, transmitChan chan TimestampedFrame) error {
	file, err := createEmptyTraceFile(t.traceDir, t.busName, t.converter.GetFileExtention())
	t.traceFile = file
	if err != nil {
		t.l.Info("cannot create trace file")
		return errors.Wrap(err, "creating trace file")
	}

	if err != nil {
		return errors.Wrap(err, "close trace file")
	}

	timeout := time.After(t.timeout)

	func() error {
		for {
			select {
			case <-timeout:
				t.l.Info("maximum trace time reached")
				return nil
			case receivedFrame := <-broadcastChan:
				t.l.Info("frame recieved")
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

func (t *Tracer) Name() string {
	return "Tracer"
}

// createEmptyTraceFile creates an *os.File given information
func createEmptyTraceFile(dir string, busName string, fileSuffix string) (*os.File, error) {
	dateStr := time.Now().Format(_filenameDateFormat)
	timeStr := time.Now().Format(_filenameTimeFormat)

	fileName := fmt.Sprintf(
		"%s_%s_%s.%s",
		busName,
		dateStr,
		timeStr,
		fileSuffix,
	)

	filePath := filepath.Join(dir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "create file")
	}

	return file, nil
}
