package main

import (
	"context"
	"encoding/json"
	"fmt"

	firebase "firebase.google.com/go"
	//"firebase.google.com/go/storage"

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
	app, _, _ := config.SetupFirebase()
	type Data struct {
		Message string `json:"message"`
		Value   int    `json:"value"`
	}
	data := Data{Message: "Hello from Go!", Value: 42}
	jsonData, _ := json.Marshal(data) // Convert to JSON bytes
	err = sendToStorage(app, jsonData)
	if err != nil {
		panic(errors.Wrap(err, "failed to send to storage"))
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

func sendToStorage(app *firebase.App, jsonData []byte) error {
	ctx := context.Background()

	// Get a Storage Client
	client, err := app.Storage(ctx)
	if err != nil {
		return fmt.Errorf("error getting Storage client: %v", err)
	}

	// Specify Bucket and Object Name
	//bucketName := "your-bucket-name" // Replace with your actual bucket name
	objectName := "data.json" // Or a dynamically generated name

	// Get a Bucket Handle
	bucket, err := client.DefaultBucket()
	if err != nil {
		return fmt.Errorf("error getting default bucket: %v", err)
	}

	// Create a Writer to Upload
	wc := bucket.Object(objectName).NewWriter(ctx)
	wc.ContentType = "application/json"
	if _, err := wc.Write(jsonData); err != nil {
		return fmt.Errorf("error writing to bucket: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("error closing writer: %v", err)
	}

	fmt.Println("JSON data successfully written to Firebase Storage:", objectName)
	return nil
}
