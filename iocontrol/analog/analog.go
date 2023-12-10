package analog

import (
	"context"
	"github.com/ethereum/go-ethereum/event"
	"github.com/pkg/errors"
	"time"
)

type Direction int

const (
	Input Direction = iota
	Output
)

type DataPoint struct {
	Value float64
	err error
}

type Pin interface {
	SetDirection(context.Context, Direction) error
	Read(context.Context) (float64, error)
	StartReadStream(sampleRate int) (chan DataPoint, error)
	SubscribeToStream(chan DataPoint) (event.Subscription, error)
	StopReadStream()
	Write(context.Context, float64) error
}

type MockPin struct {
	sampTime     time.Duration
	stop         chan struct{}
	started      bool
	currentValue float64
}

func (m *MockPin)  Read() (float64, error) {
	if m.started {
		return m.currentValue
	}

	return m.read()
}

func (m *MockPin)  read() (float64, error) {
	// hw level call
	return 5.0, nil
}

func (m *MockPin) StartReadStream() (chan DataPoint, error) {
	if m.started {
		return nil, errors.New("stream already started")
	}

	dataCh := make(chan DataPoint)
	go func(dataCh chan DataPoint) {
		for {
			select {
			case <-m.stop:
				m.started = false
				return
			case <-time.After(m.sampTime):
				data, err := m.read()

				m.currentValue = data

				dataCh <- DataPoint{
					Value: data,
					err:   err,
				}
			}
		}
	}(dataCh)
	
	return dataCh, nil
}


ch = pin.start()

data <- ch

pin.Stop()
...
ch = pin.Start()
ch1 = ch

read1(ch1)
read2(ch)
data <- ch

pin.Stop()
...

