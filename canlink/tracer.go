package canlink

import (
	"context"
	//"fmt"
	"net"
	"time"

	"github.com/macformula/hil/utils"
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
	_filenameTimeFormat = "15.04.05"

	// format for year, month and day with two digits each (period delimiter)
	_filenameDateFormat = "2006.01.02"
)

// TracerOption is a type for functions operating on Tracer
type TracerOption func(*Tracer)

// Tracer listens on a CAN bus and records all traffic during a specified period
type Tracer struct {
	l          *zap.Logger
	stop       chan struct{}
	frameCh    chan TimestampedFrame
	cachedData []TimestampedFrame
	err        *utils.ResettableError
	isRunning  bool
	receiver   *socketcan.Receiver

	directory    string
	canInterface string
	timeout      time.Duration
	busName      string
	types        []TraceFile
	conn         net.Conn
}

// NewTracer returns a new Tracer
func NewTracer(
	canInterface string,
	directory string,
	l *zap.Logger,
	conn net.Conn,
	opts ...TracerOption) *Tracer {

	tracer := &Tracer{
		l:            l.Named(_loggerName),
		cachedData:   []TimestampedFrame{},
		err:          utils.NewResettaleError(),
		timeout:      _defaultTimeout,
		canInterface: canInterface,
		directory:    directory,
		busName:      canInterface,
		types:        []TraceFile{},
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

// WithFiles sets different filetypes the tracer can dump CAN data to
func WithFiles(f ...TraceFile) TracerOption {
	return func(t *Tracer) {
		for _, filetype := range f {
			t.types = append(t.types, filetype)
		}
	}
}

// Open opens a receiver and spawns a fetchData routine
func (t *Tracer) Open(ctx context.Context) error {
	t.l.Info("creating socketcan connection")

	t.receiver = socketcan.NewReceiver(t.conn)

	t.l.Info("canlink receiver created")

	// IMPORTANT: frameCh must be open before isRunning is set
	t.frameCh = make(chan TimestampedFrame, _frameBufferLength)

	go t.fetchData(ctx)

	return nil
}

// StartTrace starts receiving and caching CAN frames
func (t *Tracer) StartTrace(ctx context.Context) error {
	if !t.isRunning {
		t.l.Info("received start command")

		t.stop = make(chan struct{})

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

		for _, files := range t.types {
			err := files.dumpToFile(t.cachedData)

			if err != nil {
				t.l.Error(err.Error())
			}
		}

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

	return nil
}

// Error returns the error set during trace execution
func (t *Tracer) Error() error {
	return t.err.Err()
}

// fetchData fetches CAN frames from the socket and sends them over a buffered channel
func (t *Tracer) fetchData(ctx context.Context) {

	timeFrame := TimestampedFrame{}

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

// receiveData listens on the buffered channel for incoming CAN frames, parses them, then caches them
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
			t.cachedData = append(t.cachedData, receivedFrame)
		}
	}
}
