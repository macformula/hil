package cli

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
	"github.com/muesli/reflow/indent"
	"go.uber.org/zap"
)

const (
	_enterKey          = "enter"
	_sigKill           = "ctrl+c"
	_quitKey           = "q"
	_escapeKey         = "esc"
	_sequenceListTitle = "HIL"
	_showLastResults   = 5
)

type cliModel struct {
	l *zap.Logger

	sequenceList list.Model
	spinner      spinner.Model

	quitting bool

	statusSignal  orchestrator.StatusSignal
	resultsSignal orchestrator.ResultsSignal

	startChan   chan orchestrator.StartSignal
	resultsChan chan orchestrator.ResultsSignal
	statusChan  chan orchestrator.StatusSignal
	fatalChan   chan orchestrator.RecoverFromFatalSignal
	cancelChan  chan orchestrator.CancelTestSignal
	quit        chan orchestrator.ShutdownSignal

	currentScreen         screenState
	results               []result
	currentRunningResults []result
	currentRunningTestId  orchestrator.TestId
	testToRun             orchestrator.TestId
	testItem              sequenceItem
	orchestratorWorking   bool
	fatalErr              error
}

func newCliModel(sequences []flow.Sequence, l *zap.Logger) *cliModel {
	var sequenceItems = make([]list.Item, len(sequences))
	for i, seq := range sequences {
		sequenceItems[i] = sequenceItem(seq)
	}

	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("206"))
	sequenceList := list.New(sequenceItems, list.NewDefaultDelegate(), 0, 0)
	sequenceList.Title = _sequenceListTitle

	cli := cliModel{
		l:                     l.Named(_loggerName),
		sequenceList:          sequenceList,
		startChan:             make(chan orchestrator.StartSignal),
		resultsChan:           make(chan orchestrator.ResultsSignal),
		statusChan:            make(chan orchestrator.StatusSignal),
		cancelChan:            make(chan orchestrator.CancelTestSignal),
		fatalChan:             make(chan orchestrator.RecoverFromFatalSignal),
		currentScreen:         Idle,
		spinner:               sp,
		results:               make([]result, _showLastResults),
		currentRunningResults: make([]result, _showLastResults),
		quit:                  make(chan orchestrator.ShutdownSignal),
	}

	return &cli
}

// Open will be called at the start of the program.
func (c *cliModel) Open(ctx context.Context) error {
	go c.run()
	go c.monitorDispatcher(ctx)

	return nil
}

// Start signal is sent to the dispatcher to start a test.
func (c *cliModel) Start() chan orchestrator.StartSignal {
	return c.startChan
}

// CancelTest will signal the dispatcher to cancel execution of the current test the Cli is trying to run
func (c *cliModel) CancelTest() chan orchestrator.CancelTestSignal {
	return c.cancelChan
}

// RecoverFromFatal is sent to signal the dispatcher that the Fatal error has been fixed.
func (c *cliModel) RecoverFromFatal() chan orchestrator.RecoverFromFatalSignal {
	return c.fatalChan
}

// Status updates the cli model on the orchestrator's status.
func (c *cliModel) Status() chan orchestrator.StatusSignal {
	return c.statusChan
}

// Results are the results of a given test.
func (c *cliModel) Results() chan orchestrator.ResultsSignal {
	return c.resultsChan
}

// Quit signals the app to shut down.
func (c *cliModel) Quit() chan orchestrator.ShutdownSignal {
	return c.quit
}

// Init is the first function that will be called.
func (c *cliModel) Init() tea.Cmd {
	return c.spinner.Tick
}

// Update is called when a message is received.
func (c *cliModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Make sure these keys always quit
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		k := keyMsg.String()
		if k == _quitKey || k == _escapeKey {
			c.quitting = true
			c.quit <- struct{}{}

			return c, tea.Quit
		}
	}

	// Hand off the message and cli to the appropriate update function for the
	// appropriate view based on the current state.
	switch c.currentScreen {
	case Idle:
		cmd = c.updateIdle(msg)

		return c, cmd
	case Running:
		cmd = c.updateRunning(msg)

		return c, cmd
	case FatalError:
		cmd = c.updateFatal(msg)

		return c, cmd
	case Results:
		cmd = c.updateResults(msg)

		return c, cmd
	default:
		cmd = c.updateIdle(msg)

		return c, cmd
	}
}

// Close will be called at the end of the program.
func (c *cliModel) Close() error { // doesnt work
	c.quitting = true

	return nil
}

// View renders the program's UI, which is just a string. The view is rendered after every Update.
func (c *cliModel) View() string {
	var s string

	switch c.currentScreen {
	case Idle:
		s = c.idleView()
	case Running:
		s = c.runningView()
	case FatalError:
		s = c.fatalView()
	case Results:
		s = c.resultsView()
	default:
		c.idleView()
	}
	return indent.String("\n"+s+"\n\n", 2)
}

func (c *cliModel) updateIdle(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msgType := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		c.sequenceList.SetSize(msgType.Width-h, msgType.Height-v)
	case tea.KeyMsg:
		switch msgType.String() {
		case _enterKey:
			seqItem, ok := c.sequenceList.SelectedItem().(sequenceItem)
			if ok {
				c.testToRun = uuid.New()
				c.startChan <- orchestrator.StartSignal{
					TestId:   c.testToRun,
					Seq:      flow.Sequence(seqItem),
					Metadata: seqItem.getMetaData(),
				}
				c.testItem = seqItem
			}

			c.currentScreen = Running
			c.results = make([]result, _showLastResults)

			return c.spinner.Tick
		case _sigKill:
			c.quitting = true
			c.quit <- struct{}{}

			return tea.Quit
		}
	case spinner.TickMsg:
		c.spinner, cmd = c.spinner.Update(msgType)

		return cmd
	default:
		return nil
	}

	c.sequenceList, cmd = c.sequenceList.Update(msg)

	return cmd
}

func (c *cliModel) updateRunning(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msgType := msg.(type) {
	case spinner.TickMsg:
		c.spinner, cmd = c.spinner.Update(msgType)

		return cmd
	case tea.KeyMsg:
		switch msgType.String() {
		case _sigKill:
			testId := c.testToRun
			c.cancelChan <- orchestrator.CancelTestSignal{TestId: testId}
		}
	}

	return nil
}

func (c *cliModel) updateFatal(msg tea.Msg) tea.Cmd {
	switch msgType := msg.(type) {
	case spinner.TickMsg:
		return c.spinner.Tick
	case tea.KeyMsg:
		switch msgType.String() {
		case _enterKey:
			c.currentScreen = Idle
			c.fatalChan <- orchestrator.RecoverFromFatalSignal{}

			return c.spinner.Tick
		case _sigKill:
			c.quitting = true
			c.quit <- struct{}{}

			return tea.Quit
		}
	}

	return c.spinner.Tick
}

func (c *cliModel) updateResults(msg tea.Msg) tea.Cmd {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case _enterKey:
			c.currentScreen = Idle

			return c.spinner.Tick
		case _sigKill:
			c.quitting = true
			c.quit <- struct{}{}

			return tea.Quit
		}
	}

	return nil
}

func (c *cliModel) idleView() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		docStyle.Render(c.sequenceList.View()),
		"\n",
		c.currentRunningTestView(),
	)
}

func (c *cliModel) runningView() string {
	s := "\n" + c.spinner.View()

	if c.currentRunningTestId == c.testToRun {
		s += fmt.Sprintf(" Work relevant to you (%s)...\n\n", c.statusSignal.Progress.Sequence.Name)
	} else {
		s += " Work relevant to you...\n\n"
	}

	for _, res := range c.results {
		if res.duration == 0 {
			s += "........................\n"
		} else {
			if res.passed {
				s += fmt.Sprintf("%s %s, finished in %s\n", passed("Passed"), res.name, res.duration)
			} else {
				s += fmt.Sprintf("%s %s, finished in %s\n", failed("Failed"), res.name, res.duration)
			}
		}
	}

	state := c.statusSignal.Progress.CurrentState
	if state != nil && c.orchestratorWorking {
		s += fmt.Sprintf("%s currently running...\n", state.Name())
	}

	s += helpStyle(fmt.Sprintf("\nCurrent test running: %s\n", c.testItem.Name))
	s += helpStyle(fmt.Sprintf("\nTest_ID: %s\n", c.testToRun.String()))
	s += helpStyle(fmt.Sprintf("\nCtrl+c to cancel the test\n"))

	if c.quitting {
		s += "\n"
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, indent.String(s, 1), c.currentRunningTestView())
}

func (c *cliModel) currentRunningTestView() string {
	s := "\n" + c.spinner.View()

	if c.orchestratorWorking {
		s += " Orchestrator doing some work...\n\n"
	} else {
		s += " Orchestrator currently idle...\n\n"
	}

	for _, res := range c.currentRunningResults {
		if res.duration == 0 {
			s += "........................\n"
		} else {
			if res.passed {
				s += fmt.Sprintf("%s %s, finished in %s\n", passed("Passed"), res.name, res.duration)
			} else {
				s += fmt.Sprintf("%s %s, finished in %s\n", failed("Failed"), res.name, res.duration)
			}
		}
	}

	state := c.statusSignal.Progress.CurrentState
	if state != nil && c.orchestratorWorking {
		s += fmt.Sprintf("%s currently running...\n", state.Name())
	}

	s += helpStyle(fmt.Sprintf("\nCurrent test running: %s\n", c.currentRunningTestId.String()))
	s += helpStyle(fmt.Sprintf("\nTest_ID: %s\n", c.currentRunningTestId.String()))
	s += helpStyle(fmt.Sprintf("\nQueue length: %d\n", c.statusSignal.QueueLength))

	if c.quitting {
		s += "\n"
	}

	return docStyle.Render(s)
}

func (c *cliModel) fatalView() string {
	s := "\n"
	s += helpStyle(fmt.Sprintf("\n%s\nHit \"enter\" to send the fatal recovery signal (CONTACT IVAN LANGE IF YOU DO NOT KNOW HOW TO FIX PROBLEM)\n",
		title("💀💀💀 FATAL ERROR 💀💀💀")))
	s += errorStyle(fmt.Sprintf("\nERROR: %s", c.fatalErr.Error()))
	return s
}

func (c *cliModel) resultsView() string {
	var builder strings.Builder

	results := c.resultsSignal

	builder.WriteString(fmt.Sprintf("Test ID: %s\n", results.TestId.String()))
	if results.IsPassing {
		builder.WriteString(passed(fmt.Sprintf("PASSED\n\n")))
	} else {
		builder.WriteString(failed(fmt.Sprintf("FAILED\n\n")))
	}

	if results.FailedTags != nil && len(results.FailedTags) > 0 {
		builder.WriteString("Failed Tags:\n")
		for _, tag := range results.FailedTags {
			if tag.Description == "" {
				builder.WriteString(fmt.Sprintf("\t🏷️ %s: %s\n", tag.ID, "no description provided"))
			} else {
				builder.WriteString(fmt.Sprintf("\t🏷️ %s: %s\n", tag.ID, tag.Description))
			}
		}
	} else {
		builder.WriteString("No failed tags.\n")
	}

	builder.WriteString("\n\n")

	if results.TestErrors != nil && len(results.TestErrors) > 0 {
		builder.WriteString("Errors:\n")
		for _, err := range results.TestErrors {
			builder.WriteString(fmt.Sprintf("\t❌  %s\n", err))
		}
	} else {
		builder.WriteString("No errors.\n")
	}

	builder.WriteString(helpStyle("\nPress enter to go back to Main Menu\n"))

	return builder.String()
}

func (c *cliModel) run() {
	p := tea.NewProgram(c, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		err = errors.Wrap(err, "run program")
		c.l.Error("failed to run program", zap.Error(err))

		return
	}
}

func (c *cliModel) monitorDispatcher(ctx context.Context) {
	for {
		select {
		case status := <-c.statusChan:
			c.l.Info("status signal received", zap.String("orchestrator state", status.OrchestratorState.String()))

			c.statusSignal = status
			c.currentRunningTestId = status.TestId

			c.fatalErr = status.FatalError

			c.currentRunningResults = make([]result, _showLastResults)
			c.results = make([]result, _showLastResults)
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

			c.results = make([]result, _showLastResults)
		case <-ctx.Done():
			c.l.Info("context done signal received")

			return
		}
	}
}
