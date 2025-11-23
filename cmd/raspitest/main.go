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

const helpText = "commands: high | low | read | help | exit"

func main() {
	boardPin := flag.Uint("pin", 0, "Raspberry Pi board pin number (1-40)")
	// one-shot mode (optional)
	oneSet := flag.String("set", "", "one-shot: set pin: high|low")
	oneGet := flag.Bool("get", false, "one-shot: read pin")
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
	if err := ctrl.Open(context.Background()); err != nil {
		log.Fatalf("open raspi controller: %v", err)
	}
	defer func() {
		if err := ctrl.Close(); err != nil {
			logger.Warn("failed to close raspi controller", zap.Error(err))
		}
	}()

	pin := raspi.NewDigitalPin(uint8(*boardPin))

	// One-shot path
	if *oneSet != "" || *oneGet {
		switch {
		case *oneSet != "":
			switch strings.ToLower(*oneSet) {
			case "high", "1", "on":
				if err := ctrl.WriteDigital(pin, true); err != nil {
					log.Fatalf("set high: %v", err)
				}
				fmt.Printf("P1-%d -> HIGH\n", *boardPin)
			case "low", "0", "off":
				if err := ctrl.WriteDigital(pin, false); err != nil {
					log.Fatalf("set low: %v", err)
				}
				fmt.Printf("P1-%d -> LOW\n", *boardPin)
			default:
				log.Fatalf("unknown --set value: %q (use high|low)", *oneSet)
			}
		case *oneGet:
			level, err := ctrl.ReadDigital(pin)
			if err != nil {
				log.Fatalf("read: %v", err)
			}
			if level {
				fmt.Printf("P1-%d <- HIGH\n", *boardPin)
			} else {
				fmt.Printf("P1-%d <- LOW\n", *boardPin)
			}
		}
		return
	}

	// Interactive path: read from the controlling TTY
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		log.Fatalf("open /dev/tty: %v", err)
	}
	defer tty.Close()

	fmt.Printf("raspitest ready on pin P1-%d\n%s\n", *boardPin, helpText)

	reader := bufio.NewScanner(tty)
	for {
		fmt.Fprintf(tty, "P1-%d> ", *boardPin)
		if !reader.Scan() {
			break
		}
		cmd := strings.ToLower(strings.TrimSpace(reader.Text()))
		if cmd == "" {
			continue
		}

		switch cmd {
		case "high", "1", "on":
			if err := ctrl.WriteDigital(pin, true); err != nil {
				fmt.Fprintf(tty, "error setting HIGH: %v\n", err)
				continue
			}
			fmt.Fprintln(tty, "-> HIGH")
		case "low", "0", "off":
			if err := ctrl.WriteDigital(pin, false); err != nil {
				fmt.Fprintf(tty, "error setting LOW: %v\n", err)
				continue
			}
			fmt.Fprintln(tty, "-> LOW")
		case "read", "get":
			level, err := ctrl.ReadDigital(pin)
			if err != nil {
				fmt.Fprintf(tty, "error reading pin: %v\n", err)
				continue
			}
			if level {
				fmt.Fprintln(tty, "<- HIGH")
			} else {
				fmt.Fprintln(tty, "<- LOW")
			}
		case "help", "?":
			fmt.Fprintln(tty, helpText)
		case "exit", "quit":
			fmt.Fprintln(tty, "bye")
			return
		default:
			fmt.Fprintf(tty, "unknown command %q (%s)\n", cmd, helpText)
		}
	}

	if err := reader.Err(); err != nil {
		log.Fatalf("tty read error: %v", err)
	}
}
