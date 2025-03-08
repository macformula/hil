package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"

	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/macformula/cangen/vehcan"
	"github.com/pkg/errors"
)

const (
	_canIface = "can0"
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

	fmt.Println("Using can interface: ", _canIface)

	conn, err := socketcan.DialContext(ctx, "can", _canIface)
	if err != nil {
		logger.Error("failed to dial context",
			zap.String("can_interface", _canIface),
			zap.Error(err),
		)
		return
	}

	busManager := canlink.NewBusManager(logger, &conn)

	busManager.Start(ctx)

	defer func() {
		err = busManager.Close()
		if err != nil {
			logger.Error("bus manager", zap.Error(err))
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
	packPositive := readInt("Enter value for PackPositive: ")
	packNegative := readInt("Enter value for PackNegative: ")
	packPrecharge := readInt("Enter value for PackPrecharge: ")

	contactors := vehcan.NewContactorStates()

	contactors.SetPackNegative(packNegative)
	contactors.SetPackPositive(packPositive)
	contactors.SetPackPrecharge(packPrecharge)

	err = busManager.Send(ctx, contactors)
	if err != nil {
		logger.Error("failed to send contactors message", zap.Error(err))
		return
	}
}
