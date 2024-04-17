package main

import (
	"context"
	"time"

	"github.com/macformula/hil/can_gen/VEH_CAN"
	"github.com/macformula/hil/canlink"
	"github.com/pkg/errors"
	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
)

const (
	_frame1Count = 6
	_frame2Count = 4
)

func main() {
	cfg := zap.NewDevelopmentConfig()
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	conn, err := socketcan.DialContext(context.Background(), "can", "can0")
	if err != nil {
		panic(err)
	}

	client := canlink.NewCANClient(VEH_CAN.Messages(), conn)
	tx := socketcan.NewTransmitter(conn)

	client.Open()

	// First Test
	go send(tx, VEH_CAN.NewContactor_Feedback().Frame(), 1, time.Second)
	msg, err := client.Read(context.Background(), VEH_CAN.NewContactor_Feedback(), VEH_CAN.NewPack_SOC())
	if err != nil {
		logger.Error("client read", zap.Error(err))
	}

	if msg.Frame().ID != VEH_CAN.NewContactor_Feedback().Frame().ID {
		logger.Error("client read", zap.Error(errors.New("incorrect CAN frame was read")))
	}

	// Second Test
	go send(tx, VEH_CAN.NewContactor_States().Frame(), 1, time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	msg, err = client.Read(ctx, VEH_CAN.NewPack_Current_Limits())
	if err != nil {
		logger.Error("client read", zap.Error(err))
	}

	if msg != nil {
		logger.Error("client read", zap.Error(errors.New("message should not have been read")))
	}

	// Third Test
	err = client.StartTracking()
	if err != nil {
		logger.Error("start tracking", zap.Error(err))
	}

	go send(tx, VEH_CAN.NewContactor_States().Frame(), _frame1Count, time.Millisecond*10)
	go send(tx, VEH_CAN.NewPack_SOC().Frame(), _frame2Count, time.Millisecond*10)
	time.Sleep(time.Second)

	data, err := client.StopTracking()
	if err != nil {
		logger.Error("stop tracking", zap.Error(err))
	}

	// Verify correct number of frames were sent
	if data[VEH_CAN.NewContactor_States().Frame().ID] != _frame1Count || data[VEH_CAN.NewPack_SOC().Frame().ID] != _frame2Count {
		logger.Error("tracking data", zap.Error(errors.New("incorrect number of frames were sent")))
	}

	err = client.Close()
	if err != nil {
		logger.Error("client close", zap.Error(err))
	}
}

func send(tx *socketcan.Transmitter, frame can.Frame, n int, delay time.Duration) {
	for i := 0; i < n; i++ {
		time.Sleep(delay)
		if err := tx.TransmitFrame(context.Background(), frame); err != nil {
			panic(err)
		}
	}
}
