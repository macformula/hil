package results

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

//go:embed resultstemplate/resultsHtml.go.html
var templateString string

// TagSubmissionDisplay includes a pre-formatted comparison string for display purposes.
type TagSubmissionDisplay struct {
	TagID             string
	Tag               Tag
	Value             any
	IsPassing         bool
	ComparisonDisplay string
}

// TemplateData contains values to fill the results template.
type TemplateData struct {
	TestID           string
	SequenceName     string
	TagSubmissions   []TagSubmissionDisplay
	ErrorSubmissions []error
	OverallPassFail  bool
	Timestamp        string
}

// HtmlReportGenerator generates HTML reports.
type HtmlReportGenerator struct {
	templateString string
}

// NewHtmlReportGenerator creates a new HtmlReportGenerator with the default template.
func NewHtmlReportGenerator() *HtmlReportGenerator {
	return &HtmlReportGenerator{
		templateString: "",
	}
}

// Generate creates an HTML report based on the provided data.
func (g *HtmlReportGenerator) Generate(
	testID uuid.UUID,
	sequenceName string,
	tagSubmissions map[string]TagSubmission,
	errorSubmissions []error,
	overallPassFail bool,
	outputDir string,
) error {
	// Prepare the display-friendly tag submissions.
	displaySubmissions := make([]TagSubmissionDisplay, 0, len(tagSubmissions))
	generated, err := generateTagSubmissionsDisplay(tagSubmissions)
	if err != nil {
		return errors.Wrap(err, "failed to generate submission tags")
	}
	displaySubmissions = append(displaySubmissions, generated...)

	data := TemplateData{
		TestID:           testID.String(),
		SequenceName:     sequenceName,
		TagSubmissions:   displaySubmissions,
		ErrorSubmissions: errorSubmissions,
		OverallPassFail:  overallPassFail,
		Timestamp:        time.Now().Format("2006-01-02 15:04:05"),
	}

	tmpl, err := template.New("report").Parse(templateString)
	if err != nil {
		return errors.Wrap(err, "failed to parse HTML template")
	}

	fileName := fmt.Sprintf("report_%s_%s.html", sequenceName, testID.String())
	filePath := filepath.Join(outputDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(err, "failed to create report file")
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return errors.Wrap(err, "failed to execute template")
	}

	return nil
}

// formatComparison generates a display-friendly comparison string based on the ComparisonOperator.
func formatComparison(tag Tag) (string, error) {
	switch tag.CompOp {
	case Gele:
		return fmt.Sprintf("%v ≤ X ≤ %v", tag.LowerLimit, tag.UpperLimit), nil
	case Gtlt:
		return fmt.Sprintf("%v < X < %v", tag.LowerLimit, tag.UpperLimit), nil
	case Ge:
		return fmt.Sprintf("X ≥ %v", tag.LowerLimit), nil
	case Gt:
		return fmt.Sprintf("X > %v", tag.LowerLimit), nil
	case Le:
		return fmt.Sprintf("X ≤ %v", tag.UpperLimit), nil
	case Lt:
		return fmt.Sprintf("X < %v", tag.UpperLimit), nil
	case Eq:
		return fmt.Sprintf("X == %v", tag.ExpectedValue), nil
	case Log:
		return "LOG", nil // Adjust as needed for the LOG operator
	default:
		return "", fmt.Errorf("unknown ComparisonOperator: %v", tag.CompOp)
	}
}

// generateTagSubmissionsDisplay generates the submission displays for each submission tag.
func generateTagSubmissionsDisplay(tagSubmissions map[string]TagSubmission) ([]TagSubmissionDisplay, error) {
	generated := make([]TagSubmissionDisplay, 0, len(tagSubmissions))
	for tagID, submission := range tagSubmissions {
		comparison, err := formatComparison(submission.Tag)
		if err != nil {
			return generated, errors.Wrapf(err, "failed to format comparison for tag %s", tagID)
		}
		generated = append(generated, TagSubmissionDisplay{
			TagID:             tagID,
			Tag:               submission.Tag,
			Value:             submission.Value,
			IsPassing:         submission.IsPassing,
			ComparisonDisplay: comparison,
		})
	}

	return generated, nil
}
