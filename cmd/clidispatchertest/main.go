package main

import (
	"context"
	"time"

	"github.com/macformula/hil/dispatcher"
	"github.com/macformula/hil/flow"
	ftest "github.com/macformula/hil/flow/test"
	"github.com/macformula/hil/orchestrator"
	otest "github.com/macformula/hil/orchestrator/test"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_loggerName = "main.log"
)

func main() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{_loggerName}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	rp := ftest.NewSimpleResultProcessor(logger)
	s := flow.NewSequencer(rp, logger)
	d := dispatcher.NewCliDispatcher(logger)
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
