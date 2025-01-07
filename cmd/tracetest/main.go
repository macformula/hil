package main

// THIS FILE IS NOT UP TO DATE

import (
	"context"
	"os"
	"os/signal"
	"time"

	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"

	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/macformula/cangen/vehcan"
)

const (
	// Can bus select
	_busName  = "veh"
	_canIface = "vcan0"

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

	loggerConfig := zap.NewDevelopmentConfig()
	logger, err := loggerConfig.Build()

	conn, err := socketcan.DialContext(context.Background(), "can", _canIface)
	if err != nil {
		logger.Error("failed to create socket can connection",
			zap.String("can_interface", _canIface),
			zap.Error(err),
		)
		return
	}

	manager := canlink.NewBusManager(logger, &conn)

	tracerJsonl := canlink.NewTracer(
		_canIface,
		logger,
		&canlink.Jsonl{},
		canlink.WithTimeout(100*time.Second),
		canlink.WithFileName("trace_sample"),
	)
	tracerText := canlink.NewTracer(
		_canIface,
		logger,
		&canlink.Text{},
		canlink.WithTimeout(100*time.Second),
		canlink.WithFileName("trace_sample2"),
	)

	manager.Register(tracerJsonl)
	defer tracerJsonl.Close()

	manager.Register(tracerText)
	defer tracerText.Close()

	manager.Start(ctx)

	for {

	}
}

func waitForSigTerm(stop chan struct{}, logger *zap.Logger) {
	// Create a channel to receive signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Wait for a signal on the channel
	<-sig

	logger.Info("received sig term")

	// Send a value to the stop channel to signal shutdown
	close(stop)
}

func startSendMessageRoutine(
	ctx context.Context, stop chan struct{}, msgPeriod time.Duration, cc *canlink.CanClient, l *zap.Logger) {
	packState := vehcan.NewPack_State()
	packState.SetPopulated_Cells(_numCells)
	packState.SetPack_Current(0)

	ctrStates := vehcan.NewContactorStates()
	ctrStates.SetPackPositive(0)
	ctrStates.SetPackNegative(0)
	ctrStates.SetPackPrecharge(0)

	ticker := time.NewTicker(msgPeriod)
	closeContactors := time.After(_closeContactorDur)

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return
		case <-stop:
			return
		case <-closeContactors:
			ctrStates.SetPackPositive(1)
			ctrStates.SetPackNegative(1)
		case <-ticker.C:
			// +cellDeviation on even, -cellDeviation on odd
			cellVoltageDeviation := _cellVoltageAbsDeviation * float64(i%2+1) * (-1)
			packState.SetAvg_Cell_Voltage(_cellVoltage + cellVoltageDeviation)

			packVoltage := float64(packState.Populated_Cells()) * packState.Avg_Cell_Voltage()
			packState.SetPack_Inst_Voltage(packVoltage)

			// Set pack current if the contactors are closed
			if ctrStates.PackPositive() > 0 && ctrStates.PackNegative() > 0 {
				// +packCurrentDeviation on even i, -packCurrentDeviation on odd i
				packCurrentDeviation := _packCurrentDeviation * float64(i%2+1) * (-1)
				packCurrentIncr := float64(msgPeriod/time.Second) * _packCurrentIncrPerSec
				packCurrent := clamp(packState.Pack_Current()+packCurrentIncr, _minPackCurrent, _maxPackCurrent)
				packCurrent += packCurrentDeviation
				packState.SetPack_Current(packCurrent)
			} else {
				packState.SetPack_Current(0)
			}

			err := cc.Send(ctx, packState)
			if err != nil {
				l.Error("failed to send pack state", zap.Error(err))

				return
			}

			err = cc.Send(ctx, ctrStates)
			if err != nil {
				l.Error("failed to send contactor states", zap.Error(err))

				return
			}
		}
	}
}

func clamp(value, min, max float64) float64 {
	if value > max {
		return max
	} else if value < min {
		return min
	}

	return value
}