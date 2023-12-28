package flow

import "context"

type ResultProcessor interface {
	SubmitTag(ctx context.Context, tagId string, value any) (bool, error)
	CompleteTest(ctx context.Context) (bool, error)
	EncounteredError(error) error
}
