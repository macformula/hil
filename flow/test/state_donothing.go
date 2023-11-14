package test

import (
	"context"
	"errors"
	"time"
)

const (
	_name = "do_nothing_state"
)

type DoNothingState struct {
	setupCalled bool
}

func (d *DoNothingState) Timeout() time.Duration {
	return time.Minute
}

func (d *DoNothingState) Setup(ctx context.Context) error {
	d.setupCalled = true

	return nil
}

func (d *DoNothingState) Name() string {
	return "do_nothing_state"
}

func (d *DoNothingState) Run(_ context.Context) error {
	if !d.setupCalled {
		return errors.New("setup never called")
	}

	return nil
}

func (d *DoNothingState) FatalError() error {
	return nil
}
