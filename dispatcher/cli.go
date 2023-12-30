package dispatcher

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/macformula/hil/orchestrator"
	"log"
	"math/rand"

	"github.com/muesli/reflow/indent"
	"time"
)

type (
	tickMsg  struct{}
	frameMsg struct{}
)

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

type result struct {
	duration time.Duration
	desc     string
}

// toss Cli in here with its interface
// rename model to cli and add signals here
type model struct {
	list        list.Model
	Choice      int
	Chosen      bool
	Ticks       int
	Frames      int
	Progress    float64
	Loaded      bool
	Quitting    bool
	status      orchestrator.StatusSignal
	startChan   chan orchestrator.StartSignal
	resultsChan chan orchestrator.ResultSignal
	statusChan  chan orchestrator.StatusSignal
	recoverChan chan struct{}
	cancelChan  chan orchestrator.TestId
	program     *tea.Program
	spinner     spinner.Model
	results     []result
}

func (c model) Init() tea.Cmd {
	return tea.Batch(
		c.spinner.Tick,
		runPretendProcess,
	)
}

// Main update function.
func (c model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			c.Quitting = true
			return c, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	switch c.status.OrchestratorState {
	case orchestrator.Idle:
		return updateIdle(msg, c)
	case orchestrator.Running:
		return updateRunning(msg, c)
	case orchestrator.FatalError:
		return updateFatal(msg, c)
	default:
		return updateIdle(msg, c)
	}
}

// The main view, which just calls the appropriate sub-view
func (c model) View() string {
	var s string
	if c.Quitting {
		return "\n  See you later!\n\n"
	}

	switch c.status.OrchestratorState {
	case orchestrator.Idle:
		//fmt.Println("Idle state")
		s = idleView(c)
	case orchestrator.Running:
		//fmt.Println("Running state")
		s = runningView(c)
	case orchestrator.FatalError:
		//fmt.Println("FatalError state")
		s = fatalView(c)
	default:
		idleView(c)
	}
	return indent.String("\n"+s+"\n\n", 2)
}

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func updateIdle(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func updateRunning(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case processFinishedMsg:
		d := time.Duration(msg)
		res := result{desc: "hi", duration: d}
		log.Printf("%s Job finished in %s", res.desc, res.duration)
		m.results = append(m.results[1:], res)
		return m, runPretendProcess
	default:
		return m, nil
	}
}

func updateFatal(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":

		}
	}

	return m, nil
}

// Sub-views

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func idleView(m model) string {
	return docStyle.Render(m.list.View())
}

// processFinishedMsg is sent when a pretend process completes.
type processFinishedMsg time.Duration

// pretendProcess simulates a long-running process.
func runPretendProcess() tea.Msg {
	pause := time.Duration(rand.Int63n(899)+100) * time.Millisecond
	time.Sleep(pause)
	return processFinishedMsg(pause)
}

func runningView(m model) string {
	s := "\n" +
		m.spinner.View() + " Doing some work...\n\n"

	for _, res := range m.results {
		if res.duration == 0 {
			s += "........................\n"
		} else {
			s += fmt.Sprintf("%s Job finished in %s\n", res.desc, res.duration)
		}
	}

	if m.Quitting {
		s += "\n"
	}

	return indent.String(s, 1)
}

func fatalView(m model) string {
	return fmt.Sprintf("fatalView")
}
