package main

import (
	"context"
	"fmt"
	"github.com/macformula/hil/cmd/mcaptest/stubs"
	"github.com/macformula/hil/macformula/cangen/vehcan"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"

	"github.com/macformula/hil/canlink"
	"github.com/pkg/errors"
)

const (
	// Can bus select
	_busName  = "veh"
	_canIface = "can0"

	// Env config
	_timeFormat        = "2006-01-02_15-04-05"
	_logFilenameFormat = "./logs/tracetest_%s.log"
	_traceDir          = "./traces"
	_logLevel          = zap.DebugLevel

	// Timing
	_msgPeriod         = 100 * time.Millisecond
	_closeContactorDur = 2 * time.Second

	// Can message values
	_cellVoltage             = 3.3 // Volts
	_cellVoltageAbsDeviation = 0.1
	_numBatteryModules       = 6
	_numBricksPerModules     = 24
	_numCells                = _numBatteryModules * _numBricksPerModules
	_maxPackCurrent          = 300
	_minPackCurrent          = 300
	_packCurrentDeviation    = 2
	_packCurrentIncrPerSec   = 5 // amps per second
)

func main() {
	ctx := context.Background()

	formattedTime := time.Now().Format(_timeFormat)
	logFileName := fmt.Sprintf(_logFilenameFormat, formattedTime)

	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.OutputPaths = []string{logFileName}
	loggerConfig.Level = zap.NewAtomicLevelAt(_logLevel)
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(errors.Wrap(err, "failed to create logger"))
	}

	conn, err := stubs.DialContextStub(context.Background(), "can", _canIface)
	if err != nil {
		logger.Error("failed to create socket can connection",
			zap.String("can_interface", _canIface),
			zap.Error(err),
		)

		return
	}

	canClient := canlink.NewCanClient(vehcan.Messages(), conn, logger)
	mcap := canlink.NewMcap(canClient, logger)

	err = canClient.Open()
	if err != nil {
		logger.Error("open can client", zap.Error(err))

		return
	}

	fmt.Println("-------------- Starting Test --------------")
	fmt.Println("-------------- CTRL-C to Stop -------------")

	stop := make(chan struct{})

	go testMcap(ctx, stop, logger, mcap, _msgPeriod)

	waitForSigTerm(stop, logger)

	fmt.Println("-------------- Test Complete --------------")

	err = canClient.Close()
	if err != nil {
		logger.Error("close can client", zap.Error(err))
	}
}

func waitForSigTerm(stop chan struct{}, logger *zap.Logger) {
	// Create a channel to receive signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Wait for a signal on the channel
	select {
	case <-sig:
		logger.Info("received sig term")
		close(stop)
	case <-stop:
	}
}

func testMcap(ctx context.Context, stop chan struct{}, logger *zap.Logger, mcap *canlink.Mcap, msgPeriod time.Duration) {
	ticker := time.NewTicker(msgPeriod)
	defer ticker.Stop()

	packState := vehcan.NewPack_State()
	timeout := time.After(time.Second * 9)
	var frame []canlink.TimestampedFrame

	shouldExit := false

	for !shouldExit {
		select {
		case <-ctx.Done():
			shouldExit = true
		case <-stop:
			shouldExit = true
		case <-timeout:
			shouldExit = true
		case <-ticker.C:
			packState.SetPack_Current(packState.Pack_Current() + 1)
			packState.SetPopulated_Cells(packState.Populated_Cells() + 2)
			packState.SetAvg_Cell_Voltage(packState.Avg_Cell_Voltage() + 3)
			packState.SetPack_Inst_Voltage(packState.Pack_Inst_Voltage() + 4)
			frame = append(frame, canlink.TimestampedFrame{
				Frame: packState.Frame(),
				Time:  time.Now(),
			})
		}
	}

	err := mcap.DumpToFile(frame, "traces", "veh")
	if err != nil {
		logger.Error("mcap DumpToFile", zap.Error(err))
	}
	close(stop)
}
