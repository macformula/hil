package dispatcher

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/macformula/hil/dispatcher/test"
	"github.com/macformula/hil/orchestrator"
	"go.uber.org/zap"
	"io"
	"os"
)

// cliIface is responsible for managing a Cli
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

func newCli(l *zap.Logger) *model {
	items := getItems()

	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("206"))
	llist := list.New(items, list.NewDefaultDelegate(), 0, 0)
	llist.Title = "HIL"

	model := model{
		l:                     l.Named(_loggerName),
		list:                  llist,
		startChan:             make(chan orchestrator.StartSignal),
		resultsChan:           make(chan orchestrator.ResultsSignal),
		statusChan:            make(chan orchestrator.StatusSignal),
		cancelChan:            make(chan orchestrator.CancelTestSignal),
		fatalChan:             make(chan orchestrator.RecoverFromFatalSignal),
		currentScreen:         Idle,
		spinner:               sp,
		results:               make([]result, showLastResults),
		currentRunningResults: make([]result, showLastResults),
		quit:                  make(chan orchestrator.ShutdownSignal),
	}

	return &model
}

func getItems() []list.Item {
	l := []list.Item{}
	for _, s := range test.Sequences {
		l = append(l, item{
			sequence: s,
		})
	}
	return l
}

func (c *model) Close() error { // doesnt work
	c.Quitting = true
	return nil
}

func (c *model) Open(ctx context.Context) error {
	go c.run()
	go c.monitorDispatcher(ctx)

	return nil
}

func (c *model) Start() chan orchestrator.StartSignal {
	return c.startChan
}

func (c *model) run() {
	f, _ := tea.LogToFile("debug.log", "debug")
	defer f.Close()

	p := tea.NewProgram(c, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (c *model) CancelTest() chan orchestrator.CancelTestSignal {
	return c.cancelChan
}

func (c *model) RecoverFromFatal() chan orchestrator.RecoverFromFatalSignal {
	return c.fatalChan
}

func (c *model) Status() chan orchestrator.StatusSignal {
	return c.statusChan
}

func (c *model) Results() chan orchestrator.ResultsSignal {
	return c.resultsChan
}

func (c *model) Quit() chan orchestrator.ShutdownSignal {
	return c.quit
}

func (c *model) monitorDispatcher(ctx context.Context) {
	for {
		select {
		case status := <-c.statusChan:
			c.l.Info("status signal received", zap.String("orchestrator state", status.OrchestratorState.String()))

			c.statusSignal = status
			c.currentRunningTestId = status.TestId

			c.fatalErr = status.FatalError

			c.currentRunningResults = make([]result, showLastResults)
			c.results = make([]result, showLastResults)
			progress := status.Progress
			c.l.Info("progress state info",
				zap.Bools("state passed", progress.StatePassed),
				zap.Durations("state durations", progress.StateDuration),
				zap.String("testid", status.TestId.String()))

			if c.currentScreen == FatalError && status.OrchestratorState != orchestrator.FatalError {
				c.currentScreen = Idle
			}

			if status.OrchestratorState == orchestrator.Idle {
				c.orchestratorWorking = false
				continue
			} else if status.OrchestratorState == orchestrator.FatalError {
				c.orchestratorWorking = false
				c.currentScreen = FatalError
				continue
			}
			c.orchestratorWorking = true

			for i, statePassed := range progress.StatePassed {
				duration := progress.StateDuration[i]
				stateName := status.Progress.Sequence.States[i].Name()

				desc := "Passed"
				isPassed := true
				if !statePassed {
					desc = "Failed" // not useful currently
					isPassed = false
				}

				c.currentRunningResults = append(c.currentRunningResults[1:], result{
					duration: duration,
					desc:     desc,
					passed:   isPassed,
					name:     stateName,
				})

				if c.currentRunningTestId == c.testToRun {
					c.results = append(c.results[1:], result{
						duration: duration,
						desc:     desc,
						passed:   isPassed,
						name:     stateName,
					})
				}
			}
		case results := <-c.resultsChan:
			c.l.Info("results signal received",
				zap.Any("failed tags", results.FailedTags),
				zap.Bool("is passing", results.IsPassing),
				zap.Any("tagId from results", results.TestId),
				zap.Any("tagId stored", c.testToRun))

			if results.TestId == c.testToRun {
				c.resultsSignal = results
				c.currentScreen = Results
				c.testToRun = uuid.New()
			}

			c.results = make([]result, showLastResults)
		case <-ctx.Done():
			c.l.Info("context done signal received")

			c.Close()
			return
		}
	}
}
