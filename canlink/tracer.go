package canlink

import (
	"context"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"
	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	_defaultTimeout    = 30 * time.Minute
	_frameBufferLength = 10
	_defaultFileType   = ".asc"
	_loggerName        = "can_tracer"

	// format for 24-hour clock with minutes, seconds, and 4 digits
	// of precision after decimal (period and colon delimiter)
	_messageTimeFormat = "15:04:05.0000"

	// format for 24-hour clock with minutes and seconds (period delimiter)
	_filenameTimeFormat = "15.04.05"

	// format for year, month and day with two digits each (period delimiter)
	_filenameDateFormat = "2006.01.02"
)

type TracerOption func(*Tracer)

type Tracer struct {
	l          *zap.Logger
	stop       chan struct{}
	frameCh    chan TimestampedFrame
	cachedData []string
	err        *utils.ResettableError
	isRunning  bool
	receiver   *socketcan.Receiver

	directory    string
	canInterface string
	timeout      time.Duration
	busName      string
	fileType     string
}

func NewTracer(
	canInterface string,
	directory string,
	l *zap.Logger,
	opts ...TracerOption) *Tracer {
	tracer := &Tracer{
		l:            l.Named(_loggerName),
		cachedData:   []string{},
		err:          utils.NewResettaleError(),
		timeout:      _defaultTimeout,
		fileType:     _defaultFileType,
		canInterface: canInterface,
		directory:    directory,
		busName:      canInterface,
	}

	for _, o := range opts {
		o(tracer)
	}

	return tracer
}

func WithTimeout(timeout time.Duration) TracerOption {
	return func(t *Tracer) {
		t.timeout = timeout
	}
}

func WithBusName(name string) TracerOption {
	return func(t *Tracer) {
		t.busName = name
	}
}

func (t *Tracer) Open(ctx context.Context) error {
	t.l.Info("creating socketcan connection")

	conn, err := socketcan.DialContext(ctx, "canlink", t.canInterface)
	if err != nil {
		return errors.Wrap(err, "dial into socket")
	}

	t.receiver = socketcan.NewReceiver(conn)

	t.l.Info("canlink receiver created")

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

		t.l.Info("tracer is running")

		t.isRunning = true
	}

	return nil
}

// StopTrace stops the receiving and caching frames, then dumps them to a file as long as the tracer is not stopped
func (t *Tracer) StopTrace() error {
	if t.isRunning {
		t.l.Info("sending stop signal")

		t.isRunning = false
		close(t.stop)

		t.l.Info("getting file name")
		file, err := t.getFile()
		if err != nil {
			return errors.Wrap(err, "get pointer to file")
		}

		t.l.Info("dumping to file")
		err = t.dumpToFile(file)
		if err != nil {
			return errors.Wrap(err, "dump cached contents to file")
		}

		t.cachedData = nil
	}

	return nil
}
func (t *Tracer) Close() error {
	err := t.receiver.Close()
	if err != nil {
		return errors.Wrap(err, "close socketcan receiver")
	}

	close(t.frameCh)

	return nil
}

func (t *Tracer) Error() error {
	return t.err
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
			t.cachedData = append(t.cachedData, t.parseString(receivedFrame))
		}
	}
}

// parseString concatenates the frame components in a standardized format
func (t *Tracer) parseString(data TimestampedFrame) string {
	var builder strings.Builder

	_, err := builder.WriteString(time.Now().Format(_messageTimeFormat))
	if err != nil {
		t.err.Set(errors.Wrap(err, "parse frame time"))
	}

	_, err = builder.WriteString(" " + strconv.FormatUint(uint64(data.Frame.ID), 10))
	if err != nil {
		t.err.Set(errors.Wrap(err, "parse frame id"))
	}

	_, err = builder.WriteString(" Rx")
	if err != nil {
		t.err.Set(errors.Wrap(err, "write receiveData"))
	}

	_, err = builder.WriteString(" " + strconv.FormatUint(uint64(data.Frame.Length), 10))
	if err != nil {
		t.err.Set(errors.Wrap(err, "parse frame length"))
	}

	for i := uint8(0); i < data.Frame.Length; i++ {
		builder.WriteString(" " + strconv.FormatUint(uint64(data.Frame.Data[i]), 16))
		if err != nil {
			t.err.Set(errors.Wrap(err, "parse frame cached data"))
		}
	}

	return builder.String()
}

// dumpToFile writes all the frames to a passed in file
func (t *Tracer) dumpToFile(file *os.File) error {
	for _, value := range t.cachedData {
		_, err := file.WriteString(value + "\n")
		if err != nil {
			return errors.Wrap(err, "write string to file")
		}
	}

	return nil
}

// getFile makes the file name based on parameters provided
func (t *Tracer) getFile() (*os.File, error) {
	var file *os.File
	var builder strings.Builder

	_, err := builder.WriteString(t.directory + "/")
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add directory to filepath")
	}

	_, err = builder.WriteString(t.busName + "_")
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add bus name to file name")
	}

	_, err = builder.WriteString(time.Now().Format(_filenameDateFormat) + "_")
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add date to file name")
	}

	_, err = builder.WriteString(time.Now().Format(_filenameTimeFormat))
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add time to file name")
	}

	_, err = builder.WriteString(t.fileType)
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add file type to file name")
	}

	file, err = os.Create(builder.String())
	if err != nil {
		return &os.File{}, errors.Wrap(err, "create file")
	}

	return file, nil
}

func (t *Tracer) frameCatcher(fr can.Frame) {

}
