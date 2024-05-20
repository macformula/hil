package canlink

import (
	"context"
	"go.einride.tech/can"
	"go.einride.tech/can/pkg/generated"
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
	"net"

	"github.com/pkg/errors"
)

const (
	_canClientLoggerName = "can_client"
)

// MessagesDescriptor is an interface mirroring the MessagesDescriptor struct found in Einride DBCs.
type MessagesDescriptor interface {
	UnmarshalFrame(f can.Frame) (generated.Message, error)
}

// CanClient represents a connection to a CAN bus with one DBC.
type CanClient struct {
	l *zap.Logger

	md MessagesDescriptor
	rx *socketcan.Receiver
	tx *socketcan.Transmitter

	rxChan  chan can.Frame
	reading bool

	// tracker keeps track of how many CAN messages have been received (per message type)
	tracker    map[uint32]int
	tracking   bool
	stopSignal chan struct{}
}

// NewCanClient creates a new CanClient with socketcan connection.
func NewCanClient(msgs MessagesDescriptor, conn net.Conn, l *zap.Logger) *CanClient {
	return &CanClient{
		l:        l.Named(_canClientLoggerName),
		md:       msgs,
		rx:       socketcan.NewReceiver(conn),
		tx:       socketcan.NewTransmitter(conn),
		rxChan:   make(chan can.Frame),
		reading:  false,
		tracking: false,
	}
}

// UnmarshalFrame unmarshalls a CAN frame
func (c *CanClient) UnmarshalFrame(f can.Frame) (generated.Message, error) {
	return c.md.UnmarshalFrame(f)
}

// Open starts the background receiver.
func (c *CanClient) Open() error {
	go c.receive()

	return nil
}

// Close closes the socketcan receiver. This also kills the receive() goroutine.
func (c *CanClient) Close() error {
	err := c.rx.Close()
	if err != nil {
		return err
	}

	return nil
}

// receive sends received frames through rxChan, only if a frame is available and Read is in progress.
func (c *CanClient) receive() {
	for c.rx.Receive() && c.reading {
		c.rxChan <- c.rx.Frame()
	}
}

// Read is a blocking function for reading a single CAN message. If given multiple possible message types to read, it
// will return the first message received from those types. If no types are given, it will return the
// first message available.
func (c *CanClient) Read(ctx context.Context, msgsToRead ...generated.Message) (generated.Message, error) {
	c.reading = true
	defer func() {
		c.reading = false
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case frame := <-c.rxChan:
			msg, err := c.md.UnmarshalFrame(frame)
			if err != nil && !isIdNotInDatabaseError(err) {
				return nil, errors.Wrap(err, "unmarshal frame")
			} else if isIdNotInDatabaseError(err) {
				c.l.Debug("found a message we do not recognize")
				// Here we have simply read a can frame that we do not know how to unmarshal, continue to next frame.
				continue
			}

			c.l.Debug("read a message", zap.Uint32("id", msg.Frame().ID))

			// No message types were specified, return first frame of any type
			if len(msgsToRead) == 0 {
				return msg, nil
			}

			for _, msgToRead := range msgsToRead {
				if frame.ID == msgToRead.Frame().ID {
					return msg, nil
				}
			}
		}
	}
}

// Send sends a CAN frame over the bus.
func (c *CanClient) Send(ctx context.Context, msg generated.Message) error {
	frame, err := msg.MarshalFrame()
	if err != nil {
		return errors.Wrap(err, "marshal frame")
	}

	err = c.tx.TransmitFrame(ctx, frame)
	if err != nil {
		return errors.Wrap(err, "transmit frame")
	}

	return nil
}

// StartTracking initiates the tracking goroutine. This is so we can check how many CAN frames of a certain type have
// come through the CAN bus in a given time.
func (c *CanClient) StartTracking(ctx context.Context) error {
	if c.tracking {
		return errors.New("tracker is already running")
	}

	c.tracker = make(map[uint32]int)
	c.stopSignal = make(chan struct{})
	c.tracking = true

	go func(c *CanClient) {
		for {
			select {
			case <-c.stopSignal:
				c.l.Debug("received stop")
				return
			default:
				msg, err := c.Read(ctx)
				if err != nil {
					c.l.Error("failed to read message",
						zap.Error(err),
					)
					continue
				}

				if msg != nil {
					c.l.Debug("adding frame", zap.Uint32("id", msg.Frame().ID))
					c.tracker[msg.Frame().ID] += 1
				}
			}
		}
	}(c)

	return nil
}

// StopTracking stops the tracker goroutine and returns the obtained frame counts.
func (c *CanClient) StopTracking() (map[uint32]int, error) {
	if !c.tracking {
		return nil, errors.New("tracker was never started")
	}

	close(c.stopSignal)
	c.tracking = false
	return c.tracker, nil
}

// IsTracking returns whether the tracker is running or not.
func (c *CanClient) IsTracking() bool {
	return c.tracking
}
