package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/macformula/cangen/vehcan"
	"github.com/pkg/errors"
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

const (
	_canIface = "vcan0"
)

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

	reader := bufio.NewReader(os.Stdin)

	// Function to read an integer value from the user
	readInt := func(prompt string) uint8 {
		for {
			fmt.Print(prompt)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			input = strings.TrimSpace(input)
			value, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid input, please enter a valid integer")
				continue
			}
			return uint8(value)
		}
	}

	// Prompt the user for input values
	packPositive := readInt("Enter value for Pack_Positive: ")
	packNegative := readInt("Enter value for Pack_Negative: ")
	packPrecharge := readInt("Enter value for Pack_Precharge: ")

	contactors := vehcan.NewContactor_States()

	contactors.SetPack_Negative(packNegative)
	contactors.SetPack_Positive(packPositive)
	contactors.SetPack_Precharge(packPrecharge)

	err = canClient.Send(ctx, contactors)
	if err != nil {
		logger.Error("failed to send contactors message", zap.Error(err))
		return
	}
}
