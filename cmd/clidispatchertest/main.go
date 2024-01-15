package main

import (
	"context"
	"github.com/macformula/hil/results/client"
	"time"

	"github.com/macformula/hil/cli"
	dtest "github.com/macformula/hil/cli/test"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
	otest "github.com/macformula/hil/orchestrator/test"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_loggerName          = "main.log"
	_resultProcessorIp   = "localhost"
	_resultProcessorPort = "31763"
	_pushToGithub        = false
)

func main() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{_loggerName}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	rp := client.NewResultsClient(_resultProcessorIp, _resultProcessorPort, _pushToGithub)
	s := flow.NewSequencer(rp, logger)
	d := cli.NewCliDispatcher(dtest.Sequences, logger)
	d2 := otest.NewSimpleDispatcher(logger, 5*time.Second, 10*time.Second)
	o := orchestrator.NewOrchestrator(s, logger, d, d2)

	err = o.Open(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "orchestrator open"))
	}

	err = o.Run(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "orchestrator run"))
	}

	logger.Info("shutdown main program")
}
