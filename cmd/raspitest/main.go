package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/macformula/hil/iocontrol/raspi"
)

const (
	helpText = "commands: high | low | read | help | exit"
)

func main() {
	boardPin := flag.Uint("pin", 0, "Raspberry Pi board pin number (1-40)")
	flag.Parse()

	if *boardPin == 0 || *boardPin > 40 {
		log.Fatalf("pin is required and must be between 1 and 40: --pin N (board numbering)")
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer logger.Sync()

	ctrl := raspi.NewController(logger)

	ctx := context.Background()
	if err := ctrl.Open(ctx); err != nil {
		log.Fatalf("open raspi controller: %v", err)
	}
	defer func() {
		if err := ctrl.Close(); err != nil {
			logger.Warn("failed to close raspi controller", zap.Error(err))
		}
	}()

	pin := raspi.NewDigitalPin(uint8(*boardPin))

	fmt.Printf("raspitest ready on pin P1-%d\n%s\n", *boardPin, helpText)

	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("P1-%d> ", *boardPin)
		if !reader.Scan() {
			break
		}
		cmd := strings.ToLower(strings.TrimSpace(reader.Text()))
		if cmd == "" {
			continue
		}

		switch cmd {
		case "high", "1", "on":
			if err := ctrl.SetDigital(pin, true); err != nil {
				fmt.Printf("error setting HIGH: %v\n", err)
				continue
			}
			fmt.Println("-> HIGH")
		case "low", "0", "off":
			if err := ctrl.SetDigital(pin, false); err != nil {
				fmt.Printf("error setting LOW: %v\n", err)
				continue
			}
			fmt.Println("-> LOW")
		case "read", "get":
			level, err := ctrl.ReadDigital(pin)
			if err != nil {
				fmt.Printf("error reading pin: %v\n", err)
				continue
			}
			state := "LOW"
			if level {
				state = "HIGH"
			}
			fmt.Printf("<- %s\n", state)
		case "help", "?":
			fmt.Println(helpText)
		case "exit", "quit":
			fmt.Println("bye")
			return
		default:
			fmt.Printf("unknown command %q (%s)\n", cmd, helpText)
		}
	}

	if err := reader.Err(); err != nil {
		log.Fatalf("stdin read error: %v", err)
	}
}
