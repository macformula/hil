package canlink

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"go.einride.tech/can"
	"go.einride.tech/can/pkg/generated"
	"go.einride.tech/can/pkg/socketcan"
)

// MessagesDescriptor is an interface mirroring the MessagesDescriptor struct found in Einride DBCs.
type MessagesDescriptor interface {
	UnmarshalFrame(f can.Frame) (generated.Message, error)
}

// CANClient represents a connection to a CAN bus with one DBC.
type CANClient struct {
	md MessagesDescriptor
	rx *socketcan.Receiver

	// tracker keeps track of how many CAN messages have been received (per message type)
	tracker    map[uint32]uint32
	tracking   bool
	stopSignal chan struct{}
}

// NewCANClient creates a new CANClient with socketcan connection.
func NewCANClient(messages MessagesDescriptor, conn net.Conn) CANClient {
	return CANClient{
		md:       messages,
		rx:       socketcan.NewReceiver(conn),
		tracking: false,
	}
}

// Read is a blocking function for reading a single CAN message. If given multiple possible message types to read, it
// will return the first message received from those types. If no types are given, it will return the first message available.
func (c *CANClient) Read(ctx context.Context, msgsToRead ...generated.Message) (generated.Message, error) {
	for {
		select { // TODO: maybe implement a max timeout here just to be safe?
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// No available frame, break out of select
			if !c.rx.Receive() {
				break
			}

			frame := c.rx.Frame()
			// No message types were specified, return first frame of any type
			if len(msgsToRead) == 0 {
				msg, err := c.md.UnmarshalFrame(frame)
				if err != nil {
					return nil, errors.Wrap(err, "unmarshal frame")
				}
				return msg, nil
			}

			for _, msgToRead := range msgsToRead {
				if frame.ID == msgToRead.Frame().ID {
					msg, err := c.md.UnmarshalFrame(frame)
					if err != nil {
						return nil, errors.Wrap(err, "unmarshal frame")
					}
					return msg, nil
				}
			}
		}
	}
}

// StartTracking initiates the tracking goroutine. This is so we can check how many CAN frames of a certain type have
// come through the CAN bus in a given time.
func (c *CANClient) StartTracking() error {
	if c.tracking {
		return errors.New("tracker is already running")
	}

	c.tracker = make(map[uint32]uint32)
	c.stopSignal = make(chan struct{})
	c.tracking = true

	go func(c *CANClient) {
		for {
			select {
			case <-c.stopSignal:
				return
			default:
				msg, err := c.Read(context.Background())
				if err != nil { // TODO: maybe log these errors?
					continue
				}
				c.tracker[msg.Frame().ID] += 1
			}
		}
	}(c)

	return nil
}

// StopTracking stops the tracker goroutine and returns the obtained frame counts.
func (c *CANClient) StopTracking() (map[uint32]uint32, error) {
	if !c.tracking {
		return nil, errors.New("tracker was never started")
	}

	close(c.stopSignal)
	c.tracking = false
	return c.tracker, nil
}

// IsTracking returns whether the tracker is running or not.
func (c *CANClient) IsTracking() bool {
	return c.tracking
}
