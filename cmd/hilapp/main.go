package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/macformula/hil/macformula"
	"github.com/macformula/hil/macformula/state"
	"go.uber.org/zap/zapcore"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/cli"
	"github.com/macformula/hil/config"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
	results "github.com/macformula/hil/results/client"
	"github.com/pkg/errors"
	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
	"path/filepath"
)

const (
	_timeFormat      = "2006.01.02_15.04.05"
	_logFileFormat   = "hilapp_%s.log"
	_canNetwork      = "can"
	_vehCan          = "veh"
	_ptCan           = "pt"
	_defaultLogLevel = zap.DebugLevel
)

// These are set by the build.sh
var (
	GitCommit          string
	DirtyVsCleanCommit string
	Date               string
)

var (
	// Define flags with their types, default values, and descriptions
	configPath  = flag.String("config", "", "Path to config file")
	version     = flag.Bool("version", false, "Displays commit and date built")
	logLevelStr = flag.String("log", _defaultLogLevel.String(), "Changes the log level (debug, info, warn, error)")
)

func main() {
	// Parse command-line flags before accessing them.
	flag.Parse()

	if *version {
		fmt.Printf("Git Commit: %v (%s)\n", GitCommit, DirtyVsCleanCommit)
		fmt.Printf("Date Built: %v\n", Date)
		return
	}

	if *configPath == "" {
		fmt.Println("Missing required flag: --config")
	}

	logLevel, err := zapcore.ParseLevel(*logLevelStr)
	if err != nil {
		fmt.Printf("Invalid log level (%s)", *logLevelStr)
	}

	// Read config file.
	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		panic(errors.Errorf("new config (%s)", *configPath))
	}

	// Create Logger.
	logFileName := fmt.Sprintf(_logFileFormat, time.Now().Format(_timeFormat))
	logFilePath := filepath.Join(cfg.LogsDir, logFileName)

	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.OutputPaths = []string{logFilePath}

	loggerConfig.Level = zap.NewAtomicLevelAt(logLevel)

	logger, err := loggerConfig.Build()
	if err != nil {
		panic(errors.Wrap(err, "build logger"))
	}

	// Create results processor.
	var rpOptions []results.Option
	if cfg.ResultProcessorAutoStart {
		rpOptions = append(rpOptions, results.WithServerAutoStart(*configPath, cfg.ResultProcessorPath))
	}

	if cfg.ResultProcessorPushToGithub {
		rpOptions = append(rpOptions, results.WithPushReportsToGithub())
	}

	resultProcessor := results.NewResultProcessor(logger,
		cfg.ResultProcessorAddr,
		rpOptions...,
	)

	// Create sequencer.
	sequencer := flow.NewSequencer(resultProcessor, logger)

	// Create a context.
	ctx := context.Background()

	// Create socketcan connections.
	vehCanConn, err := socketcan.DialContext(ctx, _canNetwork, cfg.VehCanInterface)
	if err != nil {
		logger.Error("failed to setup veh can connection",
			zap.Error(errors.Wrap(err, "dial context")))
		return
	}

	ptCanConn, err := socketcan.DialContext(ctx, _canNetwork, cfg.PtCanInterface)
	if err != nil {
		logger.Error("failed to setup pt can connection",
			zap.Error(errors.Wrap(err, "dial context")))
		return
	}

	// Create can tracers.
	vehCanTracer := canlink.NewTracer(cfg.VehCanInterface,
		cfg.TraceDir,
		logger,
		vehCanConn,
		canlink.WithTimeout(time.Duration(cfg.CanTracerTimeoutMinutes)*time.Minute),
		canlink.WithBusName(_vehCan),
	)

	err = vehCanTracer.Open(ctx)
	if err != nil {
		logger.Error("failed to open veh can tracer",
			zap.Error(errors.Wrap(err, "dial context")))
		return
	}

	ptCanTracer := canlink.NewTracer(cfg.PtCanInterface,
		cfg.TraceDir,
		logger,
		ptCanConn,
		canlink.WithTimeout(time.Duration(cfg.CanTracerTimeoutMinutes)*time.Minute),
		canlink.WithBusName(_ptCan),
	)

	err = ptCanTracer.Open(ctx)
	if err != nil {
		logger.Error("failed to open pt can tracer",
			zap.Error(errors.Wrap(err, "dial context")))
		return
	}

	// Create AppState.
	appState := macformula.AppState{
		Config:       cfg,
		VehCanTracer: vehCanTracer,
		PtCanTracer:  ptCanTracer,
	}

	// Create sequences.
	sequences := state.GetSequences(&appState, logger)

	// Create command line dispatcher.
	cliDispatcher := cli.NewCliDispatcher(sequences, logger)

	// Create orchestrator.
	orch := orchestrator.NewOrchestrator(sequencer, logger, cliDispatcher)

	// Shutdown gracefully.
	defer shutdownHandler(orch, logger)

	// Open orchestrator. This also opens all objects managed by orchestrator.
	err = orch.Open(ctx)
	if err != nil {
		logger.Error("failed to open orchestrator",
			zap.Error(errors.Wrap(err, "orchestrator open")))
		return
	}

	err = orch.Run(ctx)
	if err != nil {
		logger.Error("orchestrator run error",
			zap.Error(errors.Wrap(err, "orchestrator run")))
		return
	}

	log.Info("hil app shutting down")
}

func shutdownHandler(orchestrator *orchestrator.Orchestrator, logger *zap.Logger) {
	panicMsg := recover()

	if panicMsg != nil {
		logger.Error("panic recovered", zap.Any("panic", panicMsg))
	}

	err := orchestrator.Close()
	if err != nil {
		logger.Error("orchestrator close", zap.Error(err))
	}

	if panicMsg != nil {
		panic(panicMsg)
	}
}
