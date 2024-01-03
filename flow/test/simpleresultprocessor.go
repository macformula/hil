package test

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	_loggerName = "simpleresultprocessor"
)

type SimpleResultProcessor struct {
	submissions     map[string]any
	overallPassFail bool
	l               *zap.Logger
}

func NewSimpleResultProcessor(l *zap.Logger) *SimpleResultProcessor {
	return &SimpleResultProcessor{
		l: l.Named(_loggerName),
	}
}

func (s *SimpleResultProcessor) Open(ctx context.Context) error {
	s.l.Info("simple result processor open")
	s.submissions = make(map[string]any)

	return nil
}

func (s *SimpleResultProcessor) SubmitTag(ctx context.Context, tagId string, value any) (bool, error) {
	s.l.Info("simple result processor submit tag", zap.String("tagId", tagId), zap.Any("value", value))
	s.submissions[tagId] = value

	if tagId == "FAIL" {
		s.overallPassFail = false

		return false, nil
	}

	return true, nil
}

func (s *SimpleResultProcessor) CompleteTest(ctx context.Context, testId uuid.UUID) (bool, error) {
	return s.overallPassFail, nil
}

func (s *SimpleResultProcessor) EncounteredError(ctx context.Context, err error) error {
	s.overallPassFail = false
	return nil
}