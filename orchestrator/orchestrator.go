package orchestrator

import (
	"context"
	"time"
	"fmt"

	"go.uber.org/zap"
)

type Ochestrator struct {
	l *zap.Logger
}

func NewOrchestrator(l *zap.Logger) *Ochestrator {
	return &Ochestrator{
		l: l,
	}
}

func (o *Ochestrator) Open(ctx context.Context) error {
	return nil
}

func (o *Ochestrator) StartTests(ctx context.Context) error {
	time.Sleep(10 * time.Second) // mimick tests
	fmt.Println("done orchestrator")
	return nil
}

func (o *Ochestrator) Close(ctx context.Context) error {
	return nil
}
