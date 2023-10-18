package iocontrol

import "go.uber.org/zap"

type IOControl struct {
	l *zap.Logger
}

func NewIOControl(l *zap.Logger) *IOControl {
	return &IOControl{
		l: l,
	}
}
