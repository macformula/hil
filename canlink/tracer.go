package canlink

import (
	"net"
	"time"

	"github.com/macformula/hil/canlink/writer"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
)

const (
	_defaultTimeout    = 30 * time.Minute
	_frameBufferLength = 10
	_loggerName        = "can_tracer"
)

// TracerOption is a type for functions operating on Tracer
type TracerOption func(*Tracer)

// Tracer listens on a CAN bus and records all traffic
type Tracer struct {
	l          *zap.Logger
	frameCh    chan writer.TimestampedFrame // This channel will be provided by the bus manager
	err        *utils.ResettableError
	receiver   *socketcan.Receiver

	traceDir string
	writer  writer.WriterIface

	canInterface string
	timeout      time.Duration
	busName      string
	conn         net.Conn
}

// NewTracer returns a new Tracer
func NewTracer(
	canInterface string,
	traceDir string,
	l *zap.Logger,
	frameCh chan writer.TimestampedFrame,
	opts ...TracerOption) *Tracer {

	tracer := &Tracer{
		l:            l.Named(_loggerName),
		err:          utils.NewResettaleError(),
		timeout:      _defaultTimeout,
		canInterface: canInterface,
		traceDir:     traceDir,
		busName:      canInterface,
		frameCh: frameCh,
	}
	
	for _, o := range opts {
		o(tracer)
	}

	tracer.writer.CreateTraceFile("traces", tracer.busName)
	go tracer.handleIncomingFrames()

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

// WithWriter sets the writer to manage the trace file
func WithWriter(writer writer.WriterIface) TracerOption {
	return func(t *Tracer) {
		t.writer = writer
	}
}

// Close closes the trace file
func (t *Tracer) Close() error {
	err := t.writer.CloseTraceFile()
	if err != nil {
		return errors.Wrap(err, "close trace file")
	}

	return nil
}

// Error returns the error set during trace execution
func (t *Tracer) Error() error {
	return t.err.Err()
}

// traceIncomingFrames listens to the frames in the frame channel and writes them to a file
func (t *Tracer) handleIncomingFrames() {
	timeout := time.After(t.timeout)

	for {
		select {
		case <-timeout:
			t.l.Info("maximum trace time reached")
			return
		case receivedFrame := <-t.frameCh:
			t.l.Info("frame recieved")
			t.writer.WriteFrameToFile(&receivedFrame)
		}
	}
}
