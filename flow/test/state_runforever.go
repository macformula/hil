package test

import (
	"context"
)

type ForeverState struct{}

func (d *ForeverState) Name() string {
	return "forever_state"
}

func (d *ForeverState) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (d *ForeverState) FatalError() error {
	return nil
}
