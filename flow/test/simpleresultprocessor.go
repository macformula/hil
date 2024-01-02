package test

import (
	"context"
	"github.com/google/uuid"
)

type simpleResultProcessor struct {
	submissions     map[string]any
	overallPassFail bool
}

func (s simpleResultProcessor) Open(ctx context.Context) error {
	s.submissions = make(map[string]any)

	return nil
}

func (s simpleResultProcessor) SubmitTag(ctx context.Context, tagId string, value any) (bool, error) {
	s.submissions[tagId] = value

	if tagId == "FAIL" {
		s.overallPassFail = false

		return false, nil
	}

	return true, nil
}

func (s simpleResultProcessor) CompleteTest(ctx context.Context, testId uuid.UUID) (bool, error) {
	return s.overallPassFail, nil
}

func (s simpleResultProcessor) EncounteredError(ctx context.Context, err error) error {
	s.overallPassFail = false
	return nil
}
