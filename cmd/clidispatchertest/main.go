package main

import (
	"context"
	"time"

	"github.com/macformula/hil/cli"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
	"github.com/macformula/hil/results/client"
	"github.com/macformula/hil/test"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_loggerName          = "main.log"
	_resultProcessorAddr = "localhost:31763"
	_configPath          = "./config/hil-config/config.yaml"
	_resultServerPath    = "./results/server/main.py"
)

func main() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{_loggerName}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	resultProcessor := results.NewResultProcessor(logger,
		_resultProcessorAddr,
		results.WithPushReportsToGithub(),
		results.WithServerAutoStart(_configPath, _resultServerPath),
	)
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

		err = orchestrator.Close()
		if err != nil {
			logger.Error("orchestrator close", zap.Error(err))
		}

		panic(panicMsg)
	}()

	err = orchestrator.Run(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "orchestrator run"))
	}

	logger.Info("shutdown main program")
}
