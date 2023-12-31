package dispatcher

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/macformula/hil/dispatcher/test"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
	"go.uber.org/zap"
	"log"
	"strconv"
	"strings"
	"time"
)

type model struct {
	l *zap.Logger

	list    list.Model
	spinner spinner.Model

	Choice   int
	Chosen   bool
	Ticks    int
	Frames   int
	Progress float64
	Loaded   bool
	Quitting bool

	statusSignal  orchestrator.StatusSignal
	resultsSignal orchestrator.ResultsSignal

	startChan   chan orchestrator.StartSignal
	resultsChan chan orchestrator.ResultsSignal
	statusChan  chan orchestrator.StatusSignal
	recoverChan chan orchestrator.RecoverFromFatalSignal
	cancelChan  chan orchestrator.CancelTestSignal
	quit        chan struct{} // temporary

	//program       *tea.Program
	currentScreen         screenState
	results               []result
	currentRunningResults []result
	currentRunningTestId  orchestrator.TestId // change this to Test_id
	testToRun             orchestrator.TestId
	testItem              item
	orchestratorWorking   bool
}

func (c *model) Init() tea.Cmd {
	//f := runPretendProcess(c)
	return tea.Batch(
		c.spinner.Tick,
		//f,
	)
}

// Main update function.
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
	log.Printf("%s Inside Main UPDATE %s, %p", c.resultsSignal, c.currentScreen, &c)
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

// The main view, which just calls the appropriate sub-view
func (c *model) View() string {
	var s string
	if c.Quitting {
		return "\n  See you later!\n\n"
	}

	switch c.currentScreen {
	case Idle:
		//fmt.Println("Idle state")
		s = idleView(c)
	case Running:
		//fmt.Println("Running state")
		s = runningView(c)
	case FatalError:
		//fmt.Println("FatalError state")
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
	//case frameMsg: // doesnt hit this
	//	log.Printf("%s Inside updateRunning %s", m.statusSignal, m.currentScreen)
	//	return m, frame()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.currentScreen = Idle
			m.cancelChan <- orchestrator.CancelTestSignal{TestId: m.testToRun}
			m.testToRun = orchestrator.TestId{}
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
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":

		}
	}

	return m, nil
}

func updateResults(msg tea.Msg, m *model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		//log.Printf("%s Inside updateResults %s", m.statusSignal, m.currentScreen)
		if m.Ticks == 0 {
			m.currentScreen = Idle
			return m, m.spinner.Tick
		}
		m.Ticks--
		return m, tick()
	case tea.KeyMsg:
		if msg.String() == "enter" {
			m.currentScreen = Idle
			return m, nil
		}
		return m, nil
	default:
		return m, tick()
	}
}

// Sub-views

// processFinishedMsg is sent when a pretend process completes.
type processFinishedMsg time.Duration

// pretendProcess simulates a long-running process.
func runPretendProcess(m *model) tea.Cmd {
	startTime := time.Now()
	//status := <-m.statusChan
	//m.status = status
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	return func() tea.Msg {
		return processFinishedMsg(elapsedTime)
	}
}

func idleView(m *model) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, docStyle.Render(m.list.View()), "\n"+currentRunningTestView(m))
}

func runningView(m *model) string {
	s := "\n" +
		m.spinner.View() + " Work relevant to you...\n\n"

	for _, res := range m.results {
		if res.duration == 0 {
			s += "........................\n"
		} else {
			s += fmt.Sprintf("%s Job finished in %s\n", res.desc, res.duration)
		}
	}

	s += helpStyle(fmt.Sprintf("\nCurrent test running: %s\n", m.testItem.title))
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
			s += fmt.Sprintf("%s Job finished in %s\n", res.desc, res.duration)
		}
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
	return fmt.Sprintf("fatalView")
}

func resultsView(m *model) string {
	var builder strings.Builder
	results := m.resultsSignal

	builder.WriteString(fmt.Sprintf("Test ID: %s\n", results.TestId.String()))
	builder.WriteString(fmt.Sprintf("Passing: %t\n", results.IsPassing))

	if results.FailedTags != nil && len(results.FailedTags) > 0 {
		builder.WriteString("Failed Tags:\n")
		for _, tag := range results.FailedTags {
			builder.WriteString(fmt.Sprintf("  Tag ID: %s\n", tag.TagID))
			builder.WriteString(fmt.Sprintf("  Tag Description: %s\n", tag.TagDescription))
			builder.WriteString(fmt.Sprintf("  Comparison Operator: %s\n", tag.ComparisonOperator))
			builder.WriteString(fmt.Sprintf("  Lower Limit: %f\n", tag.LowerLimit))
			builder.WriteString(fmt.Sprintf("  Upper Limit: %f\n", tag.UpperLimit))
			builder.WriteString(fmt.Sprintf("  Expected Value: %v\n", tag.ExpectedValue))
			builder.WriteString(fmt.Sprintf("  Unit: %s\n", tag.Unit))
			builder.WriteString("\n")
		}
	} else {
		builder.WriteString("No failed tags.\n")
	}
	builder.WriteString(fmt.Sprintf("\nProgram quits in %s seconds\n", colorFg(strconv.Itoa(m.Ticks), "79")))
	builder.WriteString(helpStyle("\nPress enter to go back to Main Menu\n"))

	return builder.String()
}

// NEW UTILS FILE FOR THIS
const (
	_timeAFK        = 10
	showLastResults = 5
)

var (
	term    = termenv.EnvColorProfile()
	keyword = makeFgStyle("211")
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type (
	tickMsg  struct{}
	frameMsg struct{}
)

type result struct {
	duration time.Duration
	desc     string
}

// Utils Functions

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

func getSequence(i item) flow.Sequence {
	return test.DoNothingSequence
}

func getMetaData(i item) map[string]string {
	metaData := make(map[string]string)
	metaData["title"] = i.title
	metaData["desc"] = i.desc
	return metaData
}

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

// Generate a blend of colors.
func makeRamp(colorA, colorB string, steps float64) (s []string) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, colorToHex(c))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format compatible with termenv.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
