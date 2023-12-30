package dispatcher

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
	"io"
	"os"
)

type cliInterface interface {
	io.Closer
	Open(context.Context) error
	Start() chan orchestrator.StartSignal
	CancelTest() chan orchestrator.CancelTestSignal
	Status() chan orchestrator.StatusSignal
	RecoverFromFatal() chan orchestrator.RecoverFromFatalSignal
}

func newCli() *model {
	items := []list.Item{
		item{title: "All Tests", desc: "Run all Test Suites"},
		item{title: "AMK Test Suite", desc: "Runs all tests regarding the motor"},
		item{title: "BMS Test Suite", desc: "Runs all tests regarding the battery"},
	}

	status := orchestrator.StatusSignal{
		OrchestratorState: orchestrator.Idle,
		TestId:            orchestrator.TestId{},
		Progress:          flow.Progress{},
		QueueLength:       0,
		FatalError:        nil,
	}

	const showLastResults = 5

	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("206"))
	llist := list.New(items, list.NewDefaultDelegate(), 0, 0)
	llist.Title = "HIL"

	model := model{
		list:        llist,
		startChan:   make(chan orchestrator.StartSignal),
		resultsChan: make(chan orchestrator.ResultsSignal),
		statusChan:  make(chan orchestrator.StatusSignal),
		status:      status,
		spinner:     sp,
		results:     make([]result, showLastResults),
	}

	return &model
}

func (c model) Close() error { // doesnt work
	c.Quitting = true
	fmt.Println("quitting")
	return nil
}

func (c model) Open(ctx context.Context) error {
	go c.run()
	return nil
}

func (c model) Start() chan orchestrator.StartSignal {
	return c.startChan
}

func (c model) run() {
	f, _ := tea.LogToFile("debug.log", "debug")
	defer f.Close()

	p := tea.NewProgram(c, tea.WithAltScreen())
	c.program = p

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (c model) CancelTest() chan orchestrator.CancelTestSignal {
	return c.cancelChan
}

func (c model) RecoverFromFatal() chan orchestrator.RecoverFromFatalSignal {
	return c.recoverChan
}

func (c model) Status() chan orchestrator.StatusSignal {
	return c.statusChan
}
