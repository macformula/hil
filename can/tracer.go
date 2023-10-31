package can

import (
	"context"
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
	"time"
)

const (
	_defaultTimeout    = 30 * time.Minute
	_defaultTickPeriod = 1 * time.Millisecond
	_streamQueueLength = 10
)

type TracerOption func(*Tracer)

type Tracer struct {
	l    *zap.Logger
	stop chan struct{}
	data []string

	timeout  time.Duration
	file     string
	receiver *socketcan.Receiver
}

func NewTracer(
	l *zap.Logger,
	file string,
	receiver *socketcan.Receiver,
	options ...TracerOption,
) *Tracer {

	tracer := &Tracer{
		l:        l.Named("can_tracer"),
		data:     []string{},
		receiver: receiver,
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
	stream := make(chan TimestampedFrame, _streamQueueLength)

	go t.receive(stream, ctx)
	go t.produce(stream, ctx)

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

func (t *Tracer) produce(data chan TimestampedFrame, ctx context.Context) {
	ticker := time.NewTicker(_defaultTickPeriod)
	timeFrame := TimestampedFrame{}

	for {
		select {
		case <-ctx.Done():
			t.l.Info("context deadline exceeded")
		case <-t.stop:
			t.l.Info("stop signal received")
		case <-time.After(t.timeout):
			t.l.Info("maximum trace time reached")
		case <-ticker.C:
			if t.receiver.Receive() {
				timeFrame.Frame = t.receiver.Frame()
				timeFrame.Time = time.Now()

				data <- timeFrame
			}
		}
	}
}

func (t *Tracer) receive(data chan TimestampedFrame, ctx context.Context) {
	timeFrame := TimestampedFrame{}

	for {
		select {
		case <-ctx.Done():
			t.l.Info("context deadline exceeded")
		case <-t.stop:
			t.l.Info("stop signal received")
		case <-time.After(t.timeout):
			t.l.Info("maximum trace time reached")
		case <-data:
			data <- timeFrame
		}
	}
}
