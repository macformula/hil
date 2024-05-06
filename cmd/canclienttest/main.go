package main

import (
	"context"
	"fmt"
	"github.com/macformula/hil/macformula/cangen/vehcan"
	"go.einride.tech/can/pkg/generated"
	"math"
	"time"

	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"

	"github.com/macformula/hil/canlink"
	"github.com/pkg/errors"
)

const (
	_canIface            = "can1"
	_test1MessagePeriod  = 100 * time.Millisecond
	_test1Timeout        = 10 * _test1MessagePeriod
	_test2Timeout        = 2 * time.Second
	_test3Duration       = 5 * time.Second
	_test3AcceptableDiff = 1
)

// NOTE: This requires the can interface of choice to be in loopback mode!!

func main() {
	ctx := context.Background()

	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(errors.Wrap(err, "build logger config"))
	}
	defer logger.Sync()

	conn, err := socketcan.DialContext(ctx, "can", _canIface)
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

	defer func() {
		err = canClient.Close()
		if err != nil {
			logger.Error("client close", zap.Error(err))
		}
	}()
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println("!!!!!!!REQUIRES CAN INTERFACE BE IN LOOPBACK MODE!!!!!!!")
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

	logger.Info("starting test 1")

	// First test: send a can message, read it back. Give 2 messages of interest to the read function.
	msgSent1 := vehcan.BMSBroadcast{}
	expectedMsgRead1 := msgSent1
	msgNotSent1 := vehcan.Contactor_Feedback{}
	numMsgsToSend1 := 10

	ctxTest1, cancelTest1 := context.WithTimeout(ctx, _test1Timeout)

	go startSendRoutine(ctxTest1, canClient, &msgSent1, _test1MessagePeriod, numMsgsToSend1, logger)

	actualMsgRead1, err := canClient.Read(ctxTest1, &msgNotSent1, &msgSent1)
	if err != nil {
		logger.Error("client read", zap.Error(err))

		return
	}

	cancelTest1()

	if actualMsgRead1 == nil {
		logger.Error("failed client read test1: did not read any can messages")

		return
	}

	if expectedMsgRead1.Frame().ID != actualMsgRead1.Frame().ID {
		logger.Error("failed client read test",
			zap.String("expected_msg_read", expectedMsgRead1.String()),
			zap.String("actual_msg_read", actualMsgRead1.String()),
		)

		return
	}

	logger.Info("complete test 1")
	logger.Info("starting test 2")

	// Second Test: send wrong message, expect nil when reading from canclient.
	msgSent2 := vehcan.Contactor_States{}
	tryToReadMsg2 := vehcan.Contactor_Feedback{}
	numMsgsToSend2 := int(_test2Timeout / _test1MessagePeriod)

	ctxTest2, cancelTest2 := context.WithTimeout(context.Background(), _test2Timeout)

	go startSendRoutine(ctxTest2, canClient, &msgSent2, _test1MessagePeriod, numMsgsToSend2, logger)

	msgRead2, err := canClient.Read(ctxTest2, &tryToReadMsg2)
	if err != nil {
		logger.Error("client read", zap.Error(err))

		return
	}

	cancelTest2()

	if msgRead2 != nil {
		logger.Error("read rx message when sent a different message",
			zap.Uint32("message_sent", msgSent2.Frame().ID),
			zap.Uint32("message_read", msgRead2.Frame().ID),
		)

		return
	}

	logger.Info("complete test 2")
	logger.Info("starting test 3")

	// Third Test: send messages over a certain amount of time and see if we read expected amount of messages.
	msgsToSend3 := []struct {
		msg    generated.Message
		period time.Duration
	}{
		{
			msg:    vehcan.NewContactor_States(),
			period: 10 * time.Millisecond,
		},
		{
			msg:    vehcan.NewContactor_Feedback(),
			period: 100 * time.Millisecond,
		},
		{
			msg:    vehcan.NewThermistorBroadcast(),
			period: 1000 * time.Millisecond,
		},
	}

	ctxTest3, cancelTest3 := context.WithTimeout(ctx, _test3Duration+100*time.Second)

	// Num msgs to send condition will never be met
	numMsgs := -1

	for _, periodicMessage := range msgsToSend3 {
		logger.Debug("starting send routine",
			zap.Duration("period", periodicMessage.period))
		go startSendRoutine(ctxTest3, canClient, periodicMessage.msg, periodicMessage.period, numMsgs, logger)
	}

	err = canClient.StartTracking(ctxTest3)
	if err != nil {
		logger.Error("start tracking", zap.Error(err))

		return
	}

	time.Sleep(_test3Duration)

	canIdToNumReceivedFrames, err := canClient.StopTracking()
	if err != nil {
		logger.Error("stop tracking", zap.Error(err))

		return
	}

	// Stop sending
	cancelTest3()

	fmt.Println(canIdToNumReceivedFrames)

	logger.Info("rx message map", zap.Any("id_to_num_received", canIdToNumReceivedFrames))
	// Verify correct number of frames were sent
	for _, periodicMessage := range msgsToSend3 {
		msgId := periodicMessage.msg.Frame().ID
		numReceivedFrames, ok := canIdToNumReceivedFrames[msgId]
		if !ok {
			logger.Error("no messages received when expected non-zero amount",
				zap.String("message", periodicMessage.msg.Descriptor().Name),
				zap.Uint32("msg_id", periodicMessage.msg.Frame().ID),
				zap.Duration("periodicity", periodicMessage.period),
				zap.Duration("test_duration", _test3Duration),
			)

			return
		}

		expectedNumFrames := int(math.Floor(float64(_test3Duration / periodicMessage.period)))

		if int(math.Abs(float64(expectedNumFrames-numReceivedFrames))) > _test3AcceptableDiff {
			logger.Error("message receive count outside of tollerable range",
				zap.String("message", periodicMessage.msg.Descriptor().Name),
				zap.Uint32("msg_id", periodicMessage.msg.Frame().ID),
				zap.Duration("periodicity", periodicMessage.period),
				zap.Int("expected_num_frames", expectedNumFrames),
				zap.Int("actual_num_frames", numReceivedFrames),
				zap.Int("acceptable_diff", _test3AcceptableDiff),
			)

			return
		}
	}
}

func startSendRoutine(
	ctx context.Context,
	cc *canlink.CanClient,
	msg generated.Message,
	period time.Duration,
	numMsgsToSend int,
	l *zap.Logger) {
	// Number of messages sent
	sent := 0

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.Debug("sending msg", zap.Uint32("id", msg.Frame().ID))
			err := cc.Send(ctx, msg)
			if err != nil {
				l.Error("failed to send msg",
					zap.String("message_name", msg.Descriptor().Name),
					zap.Error(err),
				)

				return
			}

			sent += 1

			if sent == numMsgsToSend {
				return
			}
		}
	}
}
