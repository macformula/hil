package dispatcher

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"os"
	"sync"
	"time"
)

// TODO: Create a CLI struct and have options to add to it,
// one of those being adding being able to give it a channel for progress states

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list         list.Model
	selected     item
	progress     progress.Model
	progressList []progress.Model
	progressChan chan Progress
}

func (m model) Init() tea.Cmd {
	return checkProgressTick(m)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case progressMsg:
		log.Println("In progressMsg")
		p := Progress(msg)
		if p.CurrentState > p.TotalStates {
			return m, nil
		}
		if p.TotalStates == 0 {
			return m, checkProgressTick(m)
		}
		//for i := 0; i < p.CurrentState-1; i++ {
		//	m.progressList[i].SetPercent(1.0)
		//}
		//cmd1 := m.progressList[p.CurrentState-1].SetPercent(float64(p.CurrentState / p.TotalStates))
		cmd1 := m.progressList[p.CurrentState-1].SetPercent(float64(p.StateIndex) / 100) // float64(p.StateIndex / 100) doesnt work tho....

		var currentState int

		if p.StateIndex == 100 {
			currentState = p.CurrentState
		} else {
			currentState = p.CurrentState - 1
		}

		cmd2 := m.progress.SetPercent(float64(currentState) / float64(p.TotalStates))

		for i := 0; i < p.TotalStates; i++ {
			log.Printf("progress list: state=%d, percent=%f", i, m.progressList[i].Percent())
		}
		//log.Println(m.progressList)
		//cmd := m.progress.IncrPercent(0.25)
		return m, tea.Batch(cmd1, cmd2, checkProgressTick(m))
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			//m.quitting = true
			return m, tea.Quit

		case "enter":
			//if m.progress.Percent() == 1.0 {
			//	cmd := m.progress.SetPercent(0)
			//	return m, cmd
			//}
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selected = i
			}

			// Note that you can also use progress.Model.SetPercent to set the
			// percentage value explicitly, too.
			//cmd := m.progress.IncrPercent(0.25)
			return m, nil
			//return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		m.list.SetSize(msg.Width-h, msg.Height-v)
	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		cmdList := make([]tea.Cmd, len(m.progressList))
		progressModel, cmdTemp := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		cmdList = append(cmdList, cmdTemp)

		for i := 0; i < len(m.progressList); i++ {
			newProgressModel, cmd := m.progressList[i].Update(msg)
			cmdList = append(cmdList, cmd)
			m.progressList[i] = newProgressModel.(progress.Model)
		}

		return m, tea.Batch(cmdList...)
		//default:
		//	if msg == m.progressList[0].Fra
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	//if m.selected.title == "" {
	//	return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, m.list.View()))
	//}

	//return docStyle.Render(lipgloss.JoinVertical(lipgloss.Left, m.list.View(), m.selected.title, m.progress.View()))
	newArray := ""

	for _, progressbar := range m.progressList {
		newArray += progressbar.View() + "\n\n"
	}
	outer := lipgloss.JoinVertical(lipgloss.Left, m.list.View(), m.selected.title, m.progress.View())
	return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, outer, newArray))
}

func (m model) checkProgress() tea.Msg {
	log.Println("at checkProgress")
	p := <-m.progressChan
	return progressMsg(p)
}

type progressMsg Progress

func checkProgressTick(m model) tea.Cmd {
	return tea.Tick(time.Millisecond*1, func(t time.Time) tea.Msg {
		select {
		case value := <-m.progressChan:
			return progressMsg(value)
		default:
			return progressMsg{}
		}
	})
}

func Start() {
	items := []list.Item{
		item{title: "All Tests", desc: "Run all Test Suites"},
		item{title: "AMK Test Suite", desc: "Runs all tests regarding the motor"},
		item{title: "BMS Test Suite", desc: "Runs all tests regarding the battery"},
	}

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	o := NewFakeOrchestrator(nil)

	m := model{
		list:         list.New(items, list.NewDefaultDelegate(), 0, 0),
		progress:     progress.New(progress.WithDefaultGradient()),
		progressList: []progress.Model{},
		progressChan: o.GetProgressChannel(),
	}

	for i := 0; i < o.GetNumberStates(); i++ {
		newProgress := progress.New(progress.WithDefaultGradient())
		m.progressList = append(m.progressList, newProgress)
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		go o.Start()
	}()

	m.list.Title = "HIL"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	wg.Wait()
}
