package dispatcher

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
	"go.uber.org/zap"
)

type model struct {
	l *zap.Logger

	list    list.Model
	spinner spinner.Model

	Quitting bool

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
	testItem              item
	orchestratorWorking   bool
	fatalErr              error
}

// Init is the first function that will be called.
func (c *model) Init() tea.Cmd {
	return c.spinner.Tick
}

// Update is called when a message is received.
func (c *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" {
			c.Quitting = true
			c.quit <- struct{}{}
			return c, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	switch c.currentScreen {
	case Idle:
		return updateIdle(msg, c)
	case Running:
		return updateRunning(msg, c)
	case FatalError:
		return updateFatal(msg, c)
	case Results:
		return updateResults(msg, c)
	default:
		return updateIdle(msg, c)
	}
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (c *model) View() string {
	var s string

	switch c.currentScreen {
	case Idle:
		s = idleView(c)
	case Running:
		s = runningView(c)
	case FatalError:
		s = fatalView(c)
	case Results:
		s = resultsView(c)
	default:
		idleView(c)
	}
	return indent.String("\n"+s+"\n\n", 2)
}

func updateIdle(msg tea.Msg, m *model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.testToRun = uuid.New()
				m.startChan <- orchestrator.StartSignal{
					TestId:   m.testToRun,
					Seq:      getSequence(i),
					Metadata: getMetaData(i),
				}
				m.testItem = i
			}

			m.currentScreen = Running
			m.results = make([]result, showLastResults)
			return m, m.spinner.Tick
		case "ctrl+c":
			m.Quitting = true
			m.quit <- struct{}{}
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func updateRunning(msg tea.Msg, m *model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			testId := m.testToRun
			m.cancelChan <- orchestrator.CancelTestSignal{TestId: testId}
			return m, nil
		default:
			return m, nil
		}
	default:
		return m, nil
	}
}

func updateFatal(msg tea.Msg, m *model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		return m, m.spinner.Tick
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.currentScreen = Idle
			m.fatalChan <- orchestrator.RecoverFromFatalSignal{}
			return m, m.spinner.Tick
		case "ctrl+c":
			m.Quitting = true
			m.quit <- struct{}{}
			return m, tea.Quit
		}
	}

	return m, m.spinner.Tick
}

func updateResults(msg tea.Msg, m *model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.currentScreen = Idle
			return m, m.spinner.Tick
		case "ctrl+c":
			m.Quitting = true
			m.quit <- struct{}{}
			return m, tea.Quit
		default:
			return m, nil
		}
	default:
		return m, nil
	}
}

// Sub-views
func idleView(m *model) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, docStyle.Render(m.list.View()), "\n"+currentRunningTestView(m))
}

func runningView(m *model) string {
	s := "\n" +
		m.spinner.View()

	if m.currentRunningTestId == m.testToRun {
		s += fmt.Sprintf(" Work relevant to you (%s)...\n\n", m.statusSignal.Progress.Sequence.Name)
	} else {
		s += " Work relevant to you...\n\n"
	}

	for _, res := range m.results {
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

	state := m.statusSignal.Progress.CurrentState
	if state != nil && m.orchestratorWorking {
		s += fmt.Sprintf("%s currently running...\n", state.Name())
	}

	s += helpStyle(fmt.Sprintf("\nCurrent test running: %s\n", m.testItem.sequence.Name))
	s += helpStyle(fmt.Sprintf("\nTest_ID: %s\n", m.testToRun.String()))
	s += helpStyle(fmt.Sprintf("\nCtrl+c to cancel the test\n"))

	if m.Quitting {
		s += "\n"
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, indent.String(s, 1), currentRunningTestView(m))
}

func currentRunningTestView(m *model) string {
	s := "\n" + m.spinner.View()

	if m.orchestratorWorking {
		s += " Orchestrator doing some work...\n\n"
	} else {
		s += " Orchestrator currently idle...\n\n"
	}

	for _, res := range m.currentRunningResults {
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

	state := m.statusSignal.Progress.CurrentState
	if state != nil && m.orchestratorWorking {
		s += fmt.Sprintf("%s currently running...\n", state.Name())
	}

	s += helpStyle(fmt.Sprintf("\nCurrent test running: %s\n", "Unknown name for now"))
	s += helpStyle(fmt.Sprintf("\nTest_ID: %s\n", m.currentRunningTestId.String()))
	s += helpStyle(fmt.Sprintf("\nQueue length: %d\n", m.statusSignal.QueueLength))

	if m.Quitting {
		s += "\n"
	}

	return docStyle.Render(s)
}

func fatalView(m *model) string {
	s := "\n"
	s += helpStyle(fmt.Sprintf("\n%s\nHit \"enter\" to send the fatal recovery signal (CONTACT IVAN LANGE IF YOU DO NOT KNOW HOW TO FIX PROBLEM)\n",
		title("💀💀💀 FATAL ERROR 💀💀💀")))
	s += errorStyle(fmt.Sprintf("\nERROR: %s", m.fatalErr.Error()))
	return s
}

func resultsView(m *model) string {
	var builder strings.Builder
	results := m.resultsSignal

	builder.WriteString(fmt.Sprintf("Test ID: %s\n", results.TestId.String()))
	if results.IsPassing {
		builder.WriteString(passed(fmt.Sprintf("PASSED\n\n")))
	} else {
		builder.WriteString(failed(fmt.Sprintf("FAILED\n\n")))
	}

	if results.FailedTags != nil && len(results.FailedTags) > 0 {
		builder.WriteString("Failed Tags:\n")
		for _, tag := range results.FailedTags {
			if tag.TagDescription == "" {
				builder.WriteString(fmt.Sprintf("\t🏷️ %s: %s\n", tag.TagID, "no description provided"))
			} else {
				builder.WriteString(fmt.Sprintf("\t🏷️ %s: %s\n", tag.TagID, tag.TagDescription))
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

	//builder.WriteString(fmt.Sprintf("\nProgram quits in %s seconds\n", colorFg(strconv.Itoa(m.Ticks), "79")))
	builder.WriteString(helpStyle("\nPress enter to go back to Main Menu\n"))

	return builder.String()
}

// NEW UTILS FILE FOR THIS
const (
	showLastResults = 5
)

var (
	term   = termenv.EnvColorProfile()
	failed = makeFgStyle("#ff0000")
	passed = makeFgStyle("#008000")
	title  = makeFgStyle("#ffffff")
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render

type item struct {
	sequence flow.Sequence
}

// Title is the sequence name.
func (i item) Title() string {
	return i.sequence.Name
}

// Description is the sequence description.
func (i item) Description() string {
	return i.sequence.Desc
}

// FilterValue is the sequences name.
func (i item) FilterValue() string {
	return i.sequence.Name
}

type result struct {
	duration time.Duration
	desc     string
	passed   bool
	name     string
}

// Utils Functions
func getSequence(i item) flow.Sequence {
	return i.sequence
}

func getMetaData(i item) map[string]string {
	metaData := make(map[string]string)
	// No metadata required yet
	return metaData
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}
