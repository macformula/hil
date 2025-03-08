package results

import (
	"github.com/google/uuid"
)

type Generator interface {
	Generate(
		testID uuid.UUID,
		sequenceName string,
		tagSubmissions map[string]TagSubmission,
		errorSubmissions []error,
		overallPassFail bool,
		outputDir string,
	) error
}
