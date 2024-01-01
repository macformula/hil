package dispatcher

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/macformula/hil/orchestrator"
	"go.uber.org/zap"
	"io"
	"log"
	"os"
)

type cliInterface interface {
	io.Closer
	Open(context.Context) error
	Start() chan orchestrator.StartSignal
	CancelTest() chan orchestrator.CancelTestSignal
	Status() chan orchestrator.StatusSignal
	RecoverFromFatal() chan orchestrator.RecoverFromFatalSignal
	Results() chan orchestrator.ResultsSignal
	Quit() chan struct{}
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
		statusChan:            make(chan orchestrator.StatusSignal, 2),
		cancelChan:            make(chan orchestrator.CancelTestSignal),
		fatalChan:             make(chan orchestrator.RecoverFromFatalSignal),
		currentScreen:         Idle,
		spinner:               sp,
		results:               make([]result, showLastResults),
		currentRunningResults: make([]result, showLastResults),
		quit:                  make(chan struct{}),
		Ticks:                 _timeAFK,
	}

	return &model
}

func getItems() []list.Item {
	return []list.Item{
		item{title: "All Tests", desc: "Run all Test Suites"},
		item{title: "AMK Test Suite", desc: "Runs all tests regarding the motor"},
		item{title: "BMS Test Suite", desc: "Runs all tests regarding the battery"},
	}
}

func (c *model) Close() error { // doesnt work
	c.Quitting = true
	fmt.Println("quitting")
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
	//c.program = p

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

func (c *model) Quit() chan struct{} {
	return c.quit
}

func (c *model) monitorDispatcher(ctx context.Context) {
	for {
		select {
		case status := <-c.statusChan:
			c.l.Info("status signal received")

			c.statusSignal = status
			c.currentRunningTestId = status.TestId

			log.Printf("%s Inside monitorDispatcher %s, %p", status, status.OrchestratorState, &c)

			c.currentRunningResults = make([]result, showLastResults)
			c.results = make([]result, showLastResults)
			progress := status.Progress

			if status.OrchestratorState == orchestrator.Idle {
				c.orchestratorWorking = false
				continue
			} else if status.OrchestratorState == orchestrator.FatalError {
				c.orchestratorWorking = false
				c.currentScreen = FatalError
				continue
			}
			c.orchestratorWorking = true
			// populate currentRunningResults
			for i, passed := range progress.StatePassed {
				duration := progress.StateDuration[i]

				desc := "Passed"
				if !passed {
					desc = "Failed"
				}

				c.currentRunningResults = append(c.currentRunningResults[1:], result{
					duration: duration,
					desc:     desc,
				})

				if c.currentRunningTestId == c.testToRun {
					c.results = append(c.results[1:], result{
						duration: duration,
						desc:     desc,
					})
				}
			}
		case results := <-c.resultsChan:
			c.l.Info("results signal received")

			c.resultsSignal = results
			if results.TestId == c.testToRun {
				c.currentScreen = Results
			}
			c.Ticks = _timeAFK
			c.results = make([]result, showLastResults)

			log.Printf("%s Inside monitorDispatcher resultsChan %s, %p", c.resultsSignal, c.currentScreen, &c)
		case <-ctx.Done():
			c.l.Info("context done signal received")

			c.Close()
			return
		}
	}
}