package dispatcher

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

type Progress struct {
	CurrentState int
	StateIndex   int
	TotalStates  int
}

func NewProgress(CurrentState int, StateIndex int, TotalStates int) *Progress {
	return &Progress{
		CurrentState: CurrentState,
		StateIndex:   StateIndex,
		TotalStates:  TotalStates,
	}
}

type FakeOrchestrator struct {
	l            *zap.Logger
	progressChan chan Progress
	progress     Progress
}

func NewFakeOrchestrator(l *zap.Logger) *FakeOrchestrator {
	return &FakeOrchestrator{
		l:            l,
		progressChan: make(chan Progress, 100),
		progress: Progress{
			CurrentState: 1,
			StateIndex:   0,
			TotalStates:  8,
		},
	}
}

func (o *FakeOrchestrator) GetProgressChannel() chan Progress {
	return o.progressChan
}

func (o *FakeOrchestrator) GetNumberStates() int {
	return o.progress.TotalStates
}

func (o *FakeOrchestrator) Start() {
	ticker := time.NewTicker(1 * time.Second)
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	// Run the background task
	for {
		select {
		case <-ticker.C:
			log.Println("here")
			o.progress.StateIndex += 25

			if o.progress.StateIndex >= 100 {
				o.progress.StateIndex = 0

				if o.progress.CurrentState < o.progress.TotalStates {
					o.progress.CurrentState++
				}
			}

			if o.progress.CurrentState == o.progress.TotalStates && o.progress.StateIndex == 0 {
				log.Println("Reached the end of progress. Exiting.")
				o.progressChan <- o.progress
				return
			}
			o.progressChan <- o.progress
		}
	}

}
