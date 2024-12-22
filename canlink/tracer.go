package canlink

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/macformula/hil/utils"
	"github.com/macformula/hil/canlink/writer"
	"github.com/pkg/errors"
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
)

const (
	_defaultTimeout    = 30 * time.Minute
	_frameBufferLength = 10
	_loggerName        = "can_tracer"

	// format for 24-hour clock with minutes, seconds, and 4 digits
	// of precision after decimal (period and colon delimiter)
	_messageTimeFormat = "15:04:05.0000"

	// format for 24-hour clock with minutes and seconds (period delimiter)
	_filenameTimeFormat = "15-04-05"

	// format for year, month and day with two digits each (period delimiter)
	_filenameDateFormat = "2006-01-02"
)

// TracerOption is a type for functions operating on Tracer
type TracerOption func(*Tracer)

// Tracer listens on a CAN bus and records all traffic during a specified period
type Tracer struct {
	l          *zap.Logger
	stop       chan struct{}
	frameCh    chan writer.TimestampedFrame
	cachedData []writer.TimestampedFrame
	err        *utils.ResettableError
	isRunning  bool
	receiver   *socketcan.Receiver

	traceDir     string
	writers     []writer.WriterIface

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
	conn net.Conn,
	opts ...TracerOption) *Tracer {

	tracer := &Tracer{
		l:            l.Named(_loggerName),
		cachedData:   []writer.TimestampedFrame{},
		err:          utils.NewResettaleError(),
		timeout:      _defaultTimeout,
		canInterface: canInterface,
		traceDir:     traceDir,
		busName:      canInterface,
		conn:         conn,
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

// WithWriters sets the slice of writers which manage trace files
func WithWriters(writer []writer.WriterIface) TracerOption {
	return func(t *Tracer) {
		t.writers = writer
	}
}

// Open opens a receiver and spawns a fetchData routine
func (t *Tracer) Open(ctx context.Context) error {
	t.l.Info("creating socketcan connection")

	t.receiver = socketcan.NewReceiver(t.conn)

	t.l.Info("canlink receiver created")

	// IMPORTANT: frameCh must be open before isRunning is set
	t.frameCh = make(chan writer.TimestampedFrame, _frameBufferLength)

	go t.fetchData(ctx)

	return nil
}

// StartTrace starts receiving CAN frames
func (t *Tracer) StartTrace(ctx context.Context) error {
	if !t.isRunning {
		t.l.Info("received start command")

		t.stop = make(chan struct{})

		for _, writer := range t.writers {
			err := writer.CreateTraceFile(t.traceDir, t.busName)
			if err != nil {
				return errors.Wrap(err, "create trace file")
			}
			t.l.Info(fmt.Sprintf("created trace file in %s", t.traceDir))
		}

		go t.receiveData(ctx)

		t.isRunning = true

		t.l.Info("tracer is running")
	}

	return nil
}

// StopTrace stops the receiving and caching frames, then dumps them to a file as long as the tracer is not stopped
func (t *Tracer) StopTrace() error {
	if t.isRunning {
		t.l.Info("sending stop signal")

		// IMPORTANT: must set isRunning to false before closing frameCh
		t.isRunning = false

		close(t.stop)

		t.cachedData = nil
	}

	return nil
}

// Close closes the receiver
func (t *Tracer) Close() error {
	close(t.frameCh)

	err := t.StopTrace()
	if err != nil {
		return errors.Wrap(err, "stop trace")
	}

	err = t.receiver.Close()
	if err != nil {
		return errors.Wrap(err, "close socketcan receiver")
	}

	for _, writer := range t.writers {
		err = writer.CloseTraceFile()
		if err != nil {
			return errors.Wrap(err, "close trace file")
		}
	}

	return nil
}

// Error returns the error set during trace execution
func (t *Tracer) Error() error {
	return t.err.Err()
}

// fetchData fetches CAN frames from the socket and sends them over a buffered channel
func (t *Tracer) fetchData(ctx context.Context) {

	timeFrame := writer.TimestampedFrame{}

	for t.receiver.Receive() {
		select {
		case <-ctx.Done():
			t.l.Info("context deadline exceeded")
			return
		case _, ok := <-t.frameCh:
			if !ok {
				t.l.Info("frame channel closed, exiting fetch routine")
				return
			}
		default:
			timeFrame.Frame = t.receiver.Frame()
			timeFrame.Time = time.Now()

			if t.isRunning {
				t.frameCh <- timeFrame
			}
		}
	}
}

// receiveData listens on the buffered channel for incoming CAN frames, parses them, then writes them to the trace files.
func (t *Tracer) receiveData(ctx context.Context) {

	timeout := time.After(t.timeout)

	for {
		select {
		case <-ctx.Done():
			t.l.Info("context deadline exceeded")
			return
		case <-t.stop:
			t.l.Info("stop signal received")
			return
		case <-timeout:
			t.l.Info("maximum trace time reached")
			t.StopTrace()
			return
		case receivedFrame := <-t.frameCh:
			t.l.Info("frame recieved")
			t.writeFrameToFile(&receivedFrame)
		}
	}
}

func (t *Tracer) writeFrameToFile(frame *writer.TimestampedFrame) error {
	for _, writer := range t.writers {
		err := writer.WriteFrameToFile(frame)
		if err != nil {
			return errors.Wrap(err, "write frame to file")
		}
	}

	return nil
}
