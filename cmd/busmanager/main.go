package main

import (
	"context"
	"fmt"
	"time"

	//"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"

	"github.com/macformula/hil/canlink"
)

type Handler struct {
	name string,
	Handle(chan TimestampedFrame, chan struct{}) error
}

func NewHandler() *Handler {
	handler := &Handler{
		name: "testHandler",
	}

	return handler
}

func (h *Handler) Handle (broadcast chan canlink.TimestampedFrame, stopChan chan struct{}) error {
	go func() {
		for {
			select {
			case frame := <-broadcast:
				fmt.Println("RECEIVED: ", frame.Frame)
			default:
			}
		}
	}()

	go func() {
		var i byte

		for {
			time.Sleep(2 * time.Millisecond)

			frame := canlink.TimestampedFrame{}
			copy(frame.Frame.Data[:], []byte{i})
			frame.Time = time.Now()

			i = i + 1
		}
	}()

	return nil
}

func (h *Handler) Name() string {
	return "Handler"
}

func main() {
	ctx := context.Background()

	loggerConfig := zap.NewDevelopmentConfig()
	logger, err := loggerConfig.Build()

	conn, err := socketcan.DialContext(context.Background(), "can", "vcan0")
	if err != nil {
		logger.Error("failed to create socket can connection",
			zap.String("can_interface", "vcan0"),
			zap.Error(err),
		)
		return
	}

	manager := canlink.NewBusManager(logger, &conn)
	handler := NewHandler()

	manager.Register(handler)
	

	manager.Start(ctx)

	for {
	}
}
