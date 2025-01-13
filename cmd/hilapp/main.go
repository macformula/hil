package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"go.einride.tech/can/pkg/socketcan"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ethereum/go-ethereum/log"
	"github.com/macformula/hil/canlink"
	"github.com/macformula/hil/cli"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/iocontrol"
	"github.com/macformula/hil/iocontrol/sil"
	"github.com/macformula/hil/macformula"
	"github.com/macformula/hil/macformula/cangen/ptcan"
	"github.com/macformula/hil/macformula/cangen/vehcan"
	"github.com/macformula/hil/macformula/config"
	"github.com/macformula/hil/macformula/ecu/frontcontroller"
	"github.com/macformula/hil/macformula/ecu/lvcontroller"
	"github.com/macformula/hil/macformula/pinout"
	"github.com/macformula/hil/macformula/state"
	"github.com/macformula/hil/orchestrator"
	"github.com/macformula/hil/results"
	"github.com/pkg/errors"
)

const (
	_timeFormat      = "2006.01.02_15.04.05"
	_logFileFormat   = "hilapp_%s.log"
	_canNetwork      = "can"
	_vehCan          = "veh"
	_ptCan           = "pt"
	_defaultLogLevel = zap.InfoLevel
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

	// Print version and exit.
	if *version {
		fmt.Printf("Git Commit: %v (%s)\n", GitCommit, DirtyVsCleanCommit)
		fmt.Printf("Date Built: %v\n", Date)
		return
	}

	// Config path is required.
	if *configPath == "" {
		fmt.Println("Missing required flag: --config")
	}

	// Set log level.
	logLevel, err := zapcore.ParseLevel(*logLevelStr)
	if err != nil {
		fmt.Printf("Invalid log level (%s)", *logLevelStr)
	}

	// Read config file.
	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		panic(errors.Errorf("new config (%s)", *configPath))
	}

	// Get pinout revision.
	rev, err := pinout.RevisionString(cfg.Revision)
	if err != nil {
		panic(errors.Errorf("invalid revision (%s) valid options (%v)",
			cfg.Revision, pinout.RevisionStrings()))
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

	logger.Info("hil app starting", zap.Any("config", cfg))

	// Create result processor.
	resultProcessor := results.NewResultAccumulator(logger, cfg.TagsFilePath, cfg.HistoricTestsFilePath, cfg.ResultsDir,
		results.NewHtmlReportGenerator())

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

	vehBusManager := canlink.NewBusManager(logger, &vehCanConn)
	ptBusManager := canlink.NewBusManager(logger, &ptCanConn)

	// Create can tracers.
	vehCanTracer := canlink.NewTracer(
		cfg.VehCanInterface,
		logger,
		&canlink.Jsonl{},
		canlink.WithTimeout(time.Duration(cfg.CanTracerTimeoutMinutes)*time.Minute),
		canlink.WithFileName(_vehCan),
	)

	ptCanTracer := canlink.NewTracer(
		cfg.PtCanInterface,
		logger,
		&canlink.Jsonl{},
		canlink.WithTimeout(time.Duration(cfg.CanTracerTimeoutMinutes)*time.Minute),
		canlink.WithFileName(_ptCan),
	)

	// Get controllers
	var ioOpts = make([]iocontrol.IOControlOption, 0)

	switch rev {
	case pinout.Sil:
		silController := sil.NewController(cfg.SilPort, logger)
		ioOpts = append(ioOpts, iocontrol.WithSil(silController))
	default:
		panic("unconfigured revision")
	}

	// Create io controller.
	ioController := iocontrol.NewIOControl(logger, ioOpts...)

	err = ioController.Open(ctx)
	if err != nil {
		logger.Error("failed to open io controller",
			zap.Error(errors.Wrap(err, "dial context")))
		return
	}

	// Create pinout controller.
	pinoutController := pinout.NewController(rev, ioController, logger)

	err = pinoutController.Open(ctx)
	if err != nil {
		logger.Error("failed to open pinout controller",
			zap.Error(errors.Wrap(err, "dial context")))
		return
	}

	// Create testbench controller.
	testBench := macformula.NewTestBench(pinoutController, logger)

	// Create veh can client.
	vehCanClient := canlink.NewCanClient(vehcan.Messages(), vehCanConn, logger)

	// Create pt can client.
	ptCanClient := canlink.NewCanClient(ptcan.Messages(), ptCanConn, logger)

	// Create Lv Controller client.
	lvControllerClient := lvcontroller.NewClient(pinoutController, logger)

	// Create Front Controller client.
	frontControllerClient := frontcontroller.NewClient(pinoutController, vehCanClient, logger)

	// Create app object.
	app := macformula.App{
		Config:                cfg,
		VehBusManager: vehBusManager,
		PtBusManager: ptBusManager,
		VehCanTracer:          vehCanTracer,
		PtCanTracer:           ptCanTracer,
		PinoutController:      pinoutController,
		TestBench:             testBench,
		LvControllerClient:    lvControllerClient,
		FrontControllerClient: frontControllerClient,
		VehCanClient:          vehCanClient,
		PtCanClient:           ptCanClient,
	}

	// Create sequences.
	sequences := state.GetSequences(&app, logger)

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
