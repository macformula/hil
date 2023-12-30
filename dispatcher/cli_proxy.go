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
}

func newCli(l *zap.Logger) *model {
	items := getItems()

	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("206"))
	llist := list.New(items, list.NewDefaultDelegate(), 0, 0)
	llist.Title = "HIL"

	model := model{
		l:             l.Named(_loggerName),
		list:          llist,
		startChan:     make(chan orchestrator.StartSignal),
		resultsChan:   make(chan orchestrator.ResultsSignal),
		statusChan:    make(chan orchestrator.StatusSignal, 2),
		currentScreen: Idle,
		spinner:       sp,
		results:       make([]result, showLastResults),
		Ticks:         _timeAFK,
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
	c.program = p

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (c *model) CancelTest() chan orchestrator.CancelTestSignal {
	return c.cancelChan
}

func (c *model) RecoverFromFatal() chan orchestrator.RecoverFromFatalSignal {
	return c.recoverChan
}

func (c *model) Status() chan orchestrator.StatusSignal {
	return c.statusChan
}

func (c *model) Results() chan orchestrator.ResultsSignal {
	return c.resultsChan
}

func (c *model) monitorDispatcher(ctx context.Context) {
	for {
		select {
		case status := <-c.statusChan:
			c.l.Info("status signal received")

			c.statusSignal = status
			c.currentScreen = ScreenState(status.OrchestratorState)

			log.Printf("%s Inside monitorDispatcher %s, %p", status, status.OrchestratorState, &c)

			c.results = make([]result, showLastResults)
			progress := status.Progress
			for i, passed := range progress.StatePassed {
				duration := progress.StateDuration[i]

				desc := "Passed"
				if !passed {
					desc = "Failed"
				}

				c.results = append(c.results[1:], result{
					duration: duration,
					desc:     desc,
				})
			}
		case results := <-c.resultsChan:
			c.l.Info("results signal received")

			c.resultsSignal = results
			c.currentScreen = Results
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
