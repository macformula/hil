package flow

import "context"

// ResultProcessorIface will be used to get pass/fail statuses on tags.
type ResultProcessorIface interface {
	SubmitTag(ctx context.Context, tagId string, value any) (bool, error)
	CompleteTest(ctx context.Context, generateHtmlReport bool) (string, bool, error)
	EncounteredError(error) error
}
