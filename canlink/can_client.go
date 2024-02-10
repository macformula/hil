package canlink

import (
	"context"
	"net"

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
}

// NewCANClient creates a new CANClient with socketcan connection.
func NewCANClient(messages MessagesDescriptor, conn net.Conn) CANClient {
	return CANClient{
		md: messages,
		rx: socketcan.NewReceiver(conn),
	}
}

// Read is a blocking function for reading a single CAN message. If given multiple possible message types to read, it
// will return the first message received from those types.
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
			for _, msgToRead := range msgsToRead {
				if frame.ID == msgToRead.Frame().ID {
					msg, err := c.md.UnmarshalFrame(frame)
					if err != nil {
						return nil, err
					}
					return msg, nil
				}
			}
		}
	}
}
