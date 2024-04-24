package main

import (
	"context"
	"github.com/macformula/hil/cangen/vehcan"
	"go.einride.tech/can/pkg/generated"
	"time"

	"github.com/macformula/hil/canlink"
	"github.com/pkg/errors"
	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
)

const (
	_frame1Count = 6
	_frame2Count = 4
	_canIface    = "can0"
	_txRoutinePeriod = 100*time.Millisecond
	_test2Timeout = 2*time.Second
)

func main() {
	cfg := zap.NewDevelopmentConfig()
	logger, err := cfg.Build()
	if err != nil {
		panic(errors.Wrap(err, "build logger config"))
	}
	defer logger.Sync()

	conn, err := socketcan.DialContext(context.Background(), "can", _canIface)
	if err != nil {
		logger.Error("failed to dial context",
			zap.String("can_interface", _canIface),
			zap.Error(err),
		)

		return
	}

	canClient := canlink.NewCanClient(vehcan.Messages(), conn, logger)

	err = canClient.Open()
	if err != nil {
		logger.Error("failed to open can client", zap.Error(err))

		return
	}

	// First test: send a can message, read it back. Give 2 messages of interest to the read function.
	msgSent1 := vehcan.BMSBroadcast{}
	expectedMsgRead1 := msgSent1
	msgNotSent1 := vehcan.Contactor_Feedback{}

	ctx := context.Background()
	ctxTest1, cancelTest1 := context.WithCancel(ctx)

	go startSendRoutine(ctxTest1, canClient, &msgSent1, _txRoutinePeriod, logger)

	actualMsgRead1, err := canClient.Read(context.Background(), &msgNotSent1, &msgSent1)
	if err != nil {
		logger.Error("client read", zap.Error(err))
	}

	cancelTest1()

	if expectedMsgRead1.Frame().ID != actualMsgRead1.Frame().ID {
		logger.Error("failed client read test",
			zap.String("expected_msg_read", expectedMsgRead1.String()),
			zap.String("actual_msg_read", actualMsgRead1.String()),
			zap.Error(errors.New("incorrect can frame read")))
	}

	// Second Test: send wrong message, expect nil when reading from canclient.
	msgSent2 := vehcan.Contactor_States{}
	tryToReadMsg2 := vehcan.Contactor_Feedback{}

	ctxTest2, cancel := context.WithTimeout(context.Background(), _test2Timeout)
	defer cancel()

	go startSendRoutine(ctxTest2, canClient, &msgSent2, _txRoutinePeriod, logger)


	msgRead, err := canClient.Read(ctx, &tryToReadMsg2)
	if err != nil {
		logger.Error("client read", zap.Error(err))
	}

	if msgRead != nil {
		logger.Error("read message when none sent")
	}

	// Third Test
	err = canClient.StartTracking()
	if err != nil {
		logger.Error("start tracking", zap.Error(err))
	}

	go send(tx, VEH_CAN.NewContactor_States().Frame(), _frame1Count, time.Millisecond*10)
	go send(tx, VEH_CAN.NewPack_SOC().Frame(), _frame2Count, time.Millisecond*10)
	time.Sleep(time.Second)

	data, err := canClient.StopTracking()
	if err != nil {
		logger.Error("stop tracking", zap.Error(err))
	}

	// Verify correct number of frames were sent
	if data[VEH_CAN.NewContactor_States().Frame().ID] != _frame1Count || data[VEH_CAN.NewPack_SOC().Frame().ID] != _frame2Count {
		logger.Error("tracking data", zap.Error(errors.New("incorrect number of frames were sent")))
	}

	err = canClient.Close()
	if err != nil {
		logger.Error("client close", zap.Error(err))
	}
}

func send(tx *, frame can.Frame, n int, delay time.Duration) {
	for i := 0; i < n; i++ {
		time.Sleep(delay)
		if err := tx.TransmitFrame(context.Background(), frame); err != nil {
			panic(err)
		}
	}
}

func startSendRoutine(
	ctx context.Context,
	cc *canlink.CanClient,
	msg generated.Message,
	period time.Duration,
	l *zap.Logger) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(period):
			err := cc.Send(ctx, msg)
			if err != nil {
				l.Error("failed to send msg",
					zap.String("message_name", msg.String()),
					zap.Error(err),
				)
			}
		}
	}
}

