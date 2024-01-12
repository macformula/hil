package cli

import (
	"context"
	"io"

	"github.com/macformula/hil/orchestrator"
)

type cliIface interface {
	io.Closer
	// Open will be called on CliDispatcher Open
	Open(context.Context) error
	// Start signal is sent by the Cli to the CliDispatcher to signal start test
	Start() chan orchestrator.StartSignal
	// CancelTest will signal the Dispatcher to cancel execution of the current test the Cli is trying to run
	CancelTest() chan orchestrator.CancelTestSignal
	// Status signal is received when the orchestrator sends a new status
	Status() chan orchestrator.StatusSignal
	// RecoverFromFatal is sent to signal the orchestrator that the Fatal error has been fixed
	RecoverFromFatal() chan orchestrator.RecoverFromFatalSignal
	// Results signal is received when the orchestrator completes or cancels a test
	Results() chan orchestrator.ResultsSignal
	// Quit is sent when the user wants to quit the Cli
	Quit() chan orchestrator.ShutdownSignal
}
