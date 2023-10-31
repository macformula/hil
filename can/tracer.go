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

	// provide a directory, can interface, and OPTIONAL bus name (all provided by the config file)
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

	go t.receive(frame, ctx)
	go t.produce(frame, receiver, ctx)

	return nil
}

func (t *Tracer) StopTrace() error {
	t.l.Info("sending stop signal")
	close(t.stop)

	file, err := t.getFile()
	if err != nil {
		return errors.Wrap(err, "get pointer to file")
	}

	err = t.dumpToFile(file)
	if err != nil {
		return errors.Wrap(err, "dump cached contents to file")
	}
	// TODO: add resettable error that may have come up

	return nil
}

func (t *Tracer) produce(frame chan TimestampedFrame, receiver *socketcan.Receiver, ctx context.Context) {
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
		case <-time.After(_defaultTimeout):
			t.l.Info("maximum trace time reached")
			t.StopTrace()
		case <-ticker.C:
			if receiver.Receive() {
				timeFrame.Frame = receiver.Frame()
				timeFrame.Time = time.Now()

				frame <- timeFrame
			}
		}
	}
}

func (t *Tracer) receive(frame chan TimestampedFrame, ctx context.Context) {
	receivedFrame := TimestampedFrame{}

	for {
		select {
		case <-ctx.Done():
			t.l.Info("context deadline exceeded")
			return
		case <-t.stop:
			t.l.Info("stop signal received")
			return
		case <-time.After(_defaultTimeout):
			t.l.Info("maximum trace time reached")
			t.StopTrace()
		case <-frame:
			frame <- receivedFrame
			t.cachedData = append(t.cachedData, t.parseString(receivedFrame))
		}
	}
}

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
		t.err.Set(errors.Wrapf(err, "write receive"))
	}

	_, err = builder.WriteString(" " + strconv.FormatUint(uint64(data.Frame.Length), 10))
	if err != nil {
		t.err.Set(errors.Wrapf(err, "parse frame length"))
	}

	for _, v := range data.Frame.Data {
		builder.WriteString(" " + strconv.FormatUint(uint64(v), 16))
		if err != nil {
			t.err.Set(errors.Wrapf(err, "parse frame cachedData"))
		}
	}

	return builder.String()
}

func (t *Tracer) dumpToFile(file *os.File) error {
	for _, value := range t.cachedData {
		_, err := file.WriteString(value + "\n")
		if err != nil {
			return errors.Wrap(err, "write string to file")
		}
	}

	return nil
}

func (t *Tracer) getFile() (*os.File, error) {
	var file *os.File
	var builder strings.Builder

	_, err := builder.WriteString(t.directory + "/")
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add directory to filepath")
	}

	_, err = builder.WriteString(t.busName + "_")
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add busname to file name")
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
