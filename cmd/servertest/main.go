package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/macformula/hil/config"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/httpdispatcher"
	"github.com/macformula/hil/orchestrator"
	results "github.com/macformula/hil/results/client"
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

	//firebase setup code
	err = config.FirebaseDB().Connect()
	if err != nil {
		fmt.Errorf("error connecting to firebase db: %v", err)
		return
	}
	// Initiate our store
	store := flow.NewStore()
	testID := uuid.NewString()
	err = store.Create(&flow.Tag{
		ID:          testID,
		Description: "value: true",
	})
	tag, err := store.GetByID(testID)
	fmt.Println("Tag: %v\n", tag)
	if err != nil {
	}

	//rp := test.NewSimpleResultProcessor(logger)
	resultProcessor := results.NewResultProcessor(logger,
		_resultProcessorAddr,
		results.WithPushReportsToGithub(),
		//results.WithServerAutoStart(_configPath, _resultServerPath),
	)
	sequencer := flow.NewSequencer(resultProcessor, logger)
	//cliDispatcher := cli.NewCliDispatcher(test.Sequences, logger)
	//simpleDispatcher := test.NewSimpleDispatcher(logger, 5*time.Second, 10*time.Second)
	server := httpdispatcher.NewServerDispatcher(test.Sequences, httpdispatcher.NewHttpServer(logger), logger)
	o := orchestrator.NewOrchestrator(sequencer, logger, server) //cliDispatcher, simpleDispatcher)

	err = o.Open(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "orchestrator open"))
	}

	defer func() {
		panicMsg := recover()

		if panicMsg != nil {
			logger.Error("panic recovered", zap.Any("panic", panicMsg))
		}

		err = o.Close()
		if err != nil {
			logger.Error("orchestrator close", zap.Error(err))
		}

		if panicMsg != nil {
			panic(panicMsg)
		}
	}()

	err = o.Run(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "orchestrator run"))
	}

	logger.Info("shutdown main program")
}
