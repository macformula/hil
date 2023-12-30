package dispatcher

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/macformula/hil/dispatcher/test"
	"github.com/macformula/hil/flow"
	"github.com/macformula/hil/orchestrator"
	"github.com/muesli/reflow/indent"
	"go.uber.org/zap"
	"log"
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
	l             *zap.Logger
	list          list.Model
	Choice        int
	Chosen        bool
	Ticks         int
	Frames        int
	Progress      float64
	Loaded        bool
	Quitting      bool
	statusSignal  orchestrator.StatusSignal
	resultsSignal orchestrator.ResultsSignal
	startChan     chan orchestrator.StartSignal
	resultsChan   chan orchestrator.ResultsSignal
	statusChan    chan orchestrator.StatusSignal
	recoverChan   chan orchestrator.RecoverFromFatalSignal
	cancelChan    chan orchestrator.CancelTestSignal
	currentScreen orchestrator.State
	program       *tea.Program
	spinner       spinner.Model
	results       []result
	testToRun     uuid.UUID
}

func (c model) Init() tea.Cmd {
	//f := runPretendProcess(c)
	return tea.Batch(
		c.spinner.Tick,
		//f,
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
	switch c.currentScreen {
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

	switch c.currentScreen {
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
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.testToRun = uuid.New()
				m.startChan <- orchestrator.StartSignal{
					TestId:   m.testToRun,
					Seq:      m.getSequence(i),
					Metadata: m.getMetaData(i),
				}
			}
			m.currentScreen = orchestrator.Running
			return m, m.spinner.Tick
		}
	default:
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (c model) getSequence(i item) flow.Sequence {
	return test.DoNothingSequence
}

func (c model) getMetaData(i item) map[string]string {
	metaData := make(map[string]string)
	metaData["title"] = i.title
	metaData["desc"] = i.desc
	return metaData
}

func updateRunning(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case frameMsg: // doesnt hit this
		log.Printf("%s Inside updateRunning %s", m.statusSignal, m.currentScreen)
		return m, frame()
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
func runPretendProcess(m model) tea.Cmd {
	startTime := time.Now()
	//status := <-m.statusChan
	//m.status = status
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	return func() tea.Msg {
		return processFinishedMsg(elapsedTime)
	}
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
