package main

import (
	"context"
	"time"

	"github.com/macformula/hil/cli"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
	"github.com/macformula/hil/results"
	"github.com/macformula/hil/test"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_loggerName          = "main.log"
	_resultProcessorAddr = "localhost:31763"
	_resultServerPath    = "./results/server/main.py"
	_tagsPath            = "./results/tags.yaml"
	_historicTestsPath   = "./results/historic_tests.yaml"
	_reportsDir          = "./results/reports"
)

func main() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{_loggerName}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	resultProcessor := results.NewResultAccumulator(logger,
		_tagsPath, _historicTestsPath, _reportsDir, results.NewHtmlReportGenerator())
	sequencer := flow.NewSequencer(resultProcessor, logger)
	cliDispatcher := cli.NewCliDispatcher(test.Sequences, logger)
	simpleDispatcher := test.NewSimpleDispatcher(logger, 5*time.Second, 10*time.Second)
	orchestrator := orchestrator.NewOrchestrator(sequencer, logger, cliDispatcher, simpleDispatcher)

	err = orchestrator.Open(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "orchestrator open"))
	}

	defer func() {
		panicMsg := recover()

		if panicMsg != nil {
			logger.Error("panic recovered", zap.Any("panic", panicMsg))
		}

		err = orchestrator.Close()
		if err != nil {
			logger.Error("orchestrator close", zap.Error(err))
		}

		if panicMsg != nil {
			panic(panicMsg)
		}
	}()

	err = orchestrator.Run(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "orchestrator run"))
	}

	logger.Info("shutdown main program")
}
