package can

import (
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
	"time"
)

const (
	_defaultTimeout      = 30 * time.Minute
	_defaultSamplePeriod = 1 * time.Millisecond
)

type Tracer struct {
	l    *zap.Logger
	stop chan struct{}

	samplePeriod time.Duration
	timeout      time.Duration
	file         string
	receiver     *socketcan.Receiver
}

func NewTracer(
	l *zap.Logger,
	file string,
	receiver socketcan.Receiver,
	options ...func(*Tracer),
) *Tracer {

	tracer := &Tracer{
		l: l.Named("can_tracer"),
	}

	for _, o := range options {
		o(tracer)
	}

	return tracer
}

func WithTimeout(timeout time.Duration) func(*Tracer) {
	return func(t *Tracer) {
		t.timeout = timeout
	}
}
