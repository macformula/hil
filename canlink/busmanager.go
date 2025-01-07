// Package canlink provides utilities to interact with a
// Controller Area Network (CAN Bus).
package canlink

import (
	"context"
	"net"
	"sync"
	"time"

	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
)

// Default buffered channel length for the broadcast and
// incoming receive channels.
//
// In order to preserve the real time broadcast requirement,
// BusManager will drop messages for any Handler that is
// unable to keep up with the broadcast rate.
const _channelBufferLength = 1000

// BusManager is a centralized node responsible for orchestrating
// all interactions with a CAN bus.
//
// It acts as a message broker supporting the transmission
// of bus traffic to registered handlers and writing frames onto the bus.
//
// BusManager uses SocketCAN on the Linux platform. Note that
// it does not manage the lifetime of the network socket connection.
//
// Example:
//
//	 package main
//
//	 import (
//	   "context"
//	   "fmt"
//	   "time"
//
//	   "go.einride.tech/can/pkg/socketcan"
//	   "go.uber.org/zap"
//	 )
//
//	 func main () {
//	   ctx := context.Background()
//
//	   loggerConfig := zap.NewDevelopmentConfig()
//	   logger, err := loggerConfig.Build()
//
//		 // Create a network connection for vcan0
//		 conn, err := socketcan.DialContext(context.Background(), "can", "vcan0")
//		 if err != nil {
//		   return
//		 }
//
//	   manager := canlink.NewBusManager(logger, conn)
//	   handler = NewHandler(...)
//
//	   broadcastChan := manager.Register(handler)
//
//	   handler.Handle(broadcastChan)
//
//	   manager.Start(ctx)
//
//	   ...
//
//	   manager.Stop()
//	   manager.Close()
type BusManager struct {
	broadcastChan map[Handler]chan TimestampedFrame

	receiver    *socketcan.Receiver
	transmitter *socketcan.Transmitter

	l         *zap.Logger
	stop      chan struct{}
	isRunning bool
	mu        sync.Mutex
}

// NewBusManager returns a BusManager object.
//
// The network connection is injected into the BusManager
// and provides the interface for a single bus.
//
// See usage example.
func NewBusManager(l *zap.Logger, conn *net.Conn) *BusManager {
	busManager := &BusManager{
		l: l.Named("bus_manager"),

		broadcastChan: make(map[Handler]chan TimestampedFrame),

		receiver:    socketcan.NewReceiver(*conn),
		transmitter: socketcan.NewTransmitter(*conn),
	}

	return busManager
}

// Register a Handler with the BusManager.
//
// Register creates a broadcast channel for the handler, then calls the handle function on a separate go routine.
// The broadcast channel is a stream of traffic received from
// the bus.
//
// The channels operate on a TimestampedFrame object.
func (b *BusManager) Register(
	handler Handler,
) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.l.Info("registering handler")

	subscription := make(chan TimestampedFrame, _channelBufferLength)
	b.broadcastChan[handler] = subscription

	b.l.Info("registered handler")

	go handler.Handle(subscription)
}

// Unregister a Handler from the BusManager.
//
// Deletes the broadcast channel that was previously
// provided from BusManager.
func (b *BusManager) Unregister(handler *Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.broadcastChan, *handler)
}

// Start the traffic broadcast
// for each of the registered handlers.
//
// The broadcast stream will begin. If handlers cannot keep up
// with the broadcast, frames for that handler will be dropped.
func (b *BusManager) Start(ctx context.Context) {
	if b.isRunning {
		b.l.Warn("bus manager is already started")
		return
	}

	b.l.Info("start broadcast and process incoming")

	b.stop = make(chan struct{})

	go b.broadcast(ctx)

	b.isRunning = true
}

// Stop the traffic broadcast and incoming frame listener.
//
// Preserves registered handlers and their assosciated channels.
func (b *BusManager) Stop() {
	if !b.isRunning {
		b.l.Warn("bus manager is already stopped")
		return
	}

	b.l.Info("stop broadcast and process incoming")

	close(b.stop)
	b.isRunning = false
}

// Close cleans up the bus network connection.
func (b *BusManager) Close() error {
	if b.isRunning {
		b.l.Info("stopping bus manager")
		b.Stop()
	}

	b.l.Info("closing socketcan receiver and transmitter")

	b.receiver.Close()
	b.transmitter.Close()

	return nil
}

func (b *BusManager) broadcast(ctx context.Context) {
	for b.receiver.Receive() {
		timeFrame := TimestampedFrame{b.receiver.Frame(), time.Now()}

		b.mu.Lock()

		for handler, ch := range b.broadcastChan {
			select {
			case <-ctx.Done():
				b.l.Info("context deadline exceeded")
				return
			case _, ok := <-ch:
				if !ok {
					b.l.Info("broadcast channel closed, exiting broadcast routine")
					return
				}
			case <-b.stop:
				b.l.Info("stop signal received")
				return
			case ch <- timeFrame:
				b.l.Info("broadcasted can frame")
			default:
				b.l.Warn("dropping frames on handler", zap.String("handler", handler.Name()))
			}
		}

		b.mu.Unlock()
	}
}

// Send transmits frames onto the connection.
func (b *BusManager) Send(ctx context.Context, frame *TimestampedFrame) {
	b.transmitter.TransmitFrame(ctx, frame.Frame)
}
