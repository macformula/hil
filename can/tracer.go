package can

import (
	"context"
	"github.com/macformula/hil/utils"
	"github.com/pkg/errors"
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

const (
	_defaultTimeout    = 30 * time.Minute
	_defaultTickPeriod = 1 * time.Millisecond
	_streamQueueLength = 10
)

type TracerOption func(*Tracer)

type Tracer struct {
	l          *zap.Logger
	stop       chan struct{}
	cachedData []string

	timeout  time.Duration
	file     string
	receiver *socketcan.Receiver
	err      *utils.ResettableError
}

func NewTracer(
	l *zap.Logger,
	file string,
	receiver *socketcan.Receiver,
	options ...TracerOption,
) *Tracer {

	tracer := &Tracer{
		l:          l.Named("can_tracer"),
		cachedData: []string{},
		receiver:   receiver,
		err:        utils.NewResettaleError(),
		timeout:    _defaultTimeout,
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

func (t *Tracer) StartTrace(ctx context.Context) error {
	t.l.Info("received start command, beginning trace")

	t.stop = make(chan struct{})
	frame := make(chan TimestampedFrame, _streamQueueLength)

	go t.receive(frame, ctx)
	go t.produce(frame, ctx)

	return nil
}

func (t *Tracer) StopTrace() error {
	t.l.Info("sending stop signal")
	close(t.stop)

	defer t.l.Sync()

	// TODO: DUMP CACHED DATA HERE
	// TODO: add resettable error that may have come up

	return nil
}

func (t *Tracer) produce(frame chan TimestampedFrame, ctx context.Context) {
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
			if t.receiver.Receive() {
				timeFrame.Frame = t.receiver.Frame()
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

	_, err := builder.WriteString(time.Now().Format(_timeFormat))
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
