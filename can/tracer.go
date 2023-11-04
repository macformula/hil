package can

import (
	"context"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	_defaultTimeout    = 30 * time.Minute
	_defaultTickPeriod = 1 * time.Millisecond
	_streamQueueLength = 10
	_outputFileType    = ".asc"
)

type TracerOption func(*Tracer)

type Tracer struct {
	l          *zap.Logger
	stop       chan struct{}
	cachedData []string

	timeout      time.Duration
	directory    string
	canInterface string
	err          *utils.ResettableError
	busName      string
	fileType     string
}

func NewTracer(
	l *zap.Logger,
	directory string,
	canInterface string,
	options ...TracerOption,
) *Tracer {

	tracer := &Tracer{
		l:            l.Named("can_tracer"),
		cachedData:   []string{},
		err:          utils.NewResettaleError(),
		timeout:      _defaultTimeout,
		canInterface: canInterface,
		directory:    directory,
		busName:      canInterface,
		fileType:     _outputFileType,
	}

	for _, o := range options {
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

func WithFileType(fileType string) TracerOption {
	return func(t *Tracer) {
		t.fileType = fileType
	}
}

func (t *Tracer) StartTrace(ctx context.Context) error {
	t.l.Info("received start command, opening bus connection")

	conn, err := socketcan.DialContext(ctx, "can", t.canInterface)
	if err != nil {
		return errors.Wrap(err, "dial into socket")
	}

	receiver := socketcan.NewReceiver(conn)

	t.l.Info("bus receiver created, starting tracer")

	t.stop = make(chan struct{})
	frame := make(chan TimestampedFrame, _streamQueueLength)

	go t.receiveData(frame, ctx)
	go t.fetchData(frame, receiver, ctx)

	return nil
}

func (t *Tracer) StopTrace() error {
	t.l.Info("sending stop signal")
	close(t.stop)

	if t.err != nil {
		return errors.Wrap(t.err, "tracer error")
	}

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

	return nil
}

// fetchData fetches CAN frames from the socket and sends them over a buffered channel
func (t *Tracer) fetchData(frameCh chan TimestampedFrame, receiver *socketcan.Receiver, ctx context.Context) {
	ticker := time.NewTicker(_defaultTickPeriod)
	timeFrame := TimestampedFrame{}

	for {
		select {
		case <-ctx.Done():
			t.l.Info("context deadline exceeded")
			return
		case <-t.stop:
			t.l.Info("stop signal received")
			return
		case <-time.After(t.timeout):
			t.l.Info("maximum trace time reached")
			t.StopTrace()
		case <-ticker.C:
			if receiver.Receive() {
				timeFrame.Frame = receiver.Frame()
				timeFrame.Time = time.Now()

				frameCh <- timeFrame
			}
		}
	}
}

// receiveData listens on the buffered channel for incoming CAN frames, parses them, then caches them
func (t *Tracer) receiveData(frameCh chan TimestampedFrame, ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			t.l.Info("context deadline exceeded")
			return
		case <-t.stop:
			t.l.Info("stop signal received")
			return
		case receivedFrame := <-frameCh:
			t.cachedData = append(t.cachedData, t.parseString(receivedFrame))
		}
	}
}

// parseString concatenates the frame components in a standardized format
func (t *Tracer) parseString(data TimestampedFrame) string {
	var builder strings.Builder

	_, err := builder.WriteString(time.Now().Format(_tracerFormat))
	if err != nil {
		t.err.Set(errors.Wrapf(err, "parse frame time"))
	}

	_, err = builder.WriteString(" " + strconv.FormatUint(uint64(data.Frame.ID), 10))
	if err != nil {
		t.err.Set(errors.Wrapf(err, "parse frame id"))
	}

	_, err = builder.WriteString(" Rx")
	if err != nil {
		t.err.Set(errors.Wrapf(err, "write receiveData"))
	}

	_, err = builder.WriteString(" " + strconv.FormatUint(uint64(data.Frame.Length), 10))
	if err != nil {
		t.err.Set(errors.Wrapf(err, "parse frame length"))
	}

	for i := uint8(0); i < data.Frame.Length; i++ {
		builder.WriteString(" " + strconv.FormatUint(uint64(data.Frame.Data[i]), 16))
		if err != nil {
			t.err.Set(errors.Wrapf(err, "parse frame cached data"))
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

	_, err = builder.WriteString(time.Now().Format(_nameDateFormat) + "_")
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add date to file name")
	}

	_, err = builder.WriteString(time.Now().Format(_nameTimeFormat))
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
