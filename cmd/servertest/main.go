package main

import (
	"context"
	"github.com/macformula/hil/cli"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/httpdispatcher"
	"github.com/macformula/hil/orchestrator"
	"github.com/macformula/hil/test"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_loggerName          = "main.log"
	_resultProcessorIp   = "localhost"
	_resultProcessorPort = "31763"
	_pushToGithub        = true
)

func main() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{_loggerName}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	rp := test.NewSimpleResultProcessor(logger)
	//resultProcessor := client.NewResultsClient(_resultProcessorIp, _resultProcessorPort, _pushToGithub)
	sequencer := flow.NewSequencer(rp, logger)
	cliDispatcher := cli.NewCliDispatcher(test.Sequences, logger)
	//simpleDispatcher := test.NewSimpleDispatcher(logger, 5*time.Second, 10*time.Second)
	server := httpdispatcher.NewServerDispatcher(test.Sequences, httpdispatcher.NewHttpServer(logger), logger)
	o := orchestrator.NewOrchestrator(sequencer, logger, cliDispatcher, server) //cliDispatcher, simpleDispatcher)

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
