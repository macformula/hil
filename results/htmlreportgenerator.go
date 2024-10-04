package results

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// TagSubmissionDisplay includes a pre-formatted comparison string for display purposes.
type TagSubmissionDisplay struct {
	TagID             string
	Tag               Tag
	Value             any
	IsPassing         bool
	ComparisonDisplay string
}

// HtmlReportGenerator generates HTML reports.
type HtmlReportGenerator struct {
	templateString string
}

// NewHtmlReportGenerator creates a new HtmlReportGenerator with the default template.
func NewHtmlReportGenerator() *HtmlReportGenerator {
	return &HtmlReportGenerator{
		templateString: htmlTemplate,
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
	for tagID, submission := range tagSubmissions {
		comparison, err := formatComparison(submission.Tag)
		if err != nil {
			return errors.Wrapf(err, "failed to format comparison for tag %s", tagID)
		}
		displaySubmissions = append(displaySubmissions, TagSubmissionDisplay{
			TagID:             tagID,
			Tag:               submission.Tag,
			Value:             submission.Value,
			IsPassing:         submission.IsPassing,
			ComparisonDisplay: comparison,
		})
	}

	data := struct {
		TestID           string
		SequenceName     string
		TagSubmissions   []TagSubmissionDisplay
		ErrorSubmissions []error
		OverallPassFail  bool
		Timestamp        string
	}{
		TestID:           testID.String(),
		SequenceName:     sequenceName,
		TagSubmissions:   displaySubmissions,
		ErrorSubmissions: errorSubmissions,
		OverallPassFail:  overallPassFail,
		Timestamp:        time.Now().Format("2006-01-02 15:04:05"),
	}

	tmpl, err := template.New("report").Parse(g.templateString)
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

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Test Report: {{.SequenceName}}</title>
    <!-- DataTables CSS -->
    <link rel="stylesheet" href="https://cdn.datatables.net/1.13.6/css/jquery.dataTables.min.css">
    <style>
        /* Reset and Basic Styles */
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f4f6f8;
            color: #333;
            padding: 20px;
        }
        h1, h2, h3 {
            margin-bottom: 20px;
            color: #2c3e50;
        }
        p {
            margin-bottom: 10px;
            line-height: 1.6;
        }
        /* Container */
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background-color: #ffffff;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        /* Overall Result */
        .overall-result {
            padding: 15px;
            border-radius: 5px;
            text-align: center;
            font-size: 1.2em;
            margin-bottom: 30px;
            color: #fff;
        }
        .overall-result.pass {
            background-color: #27ae60;
        }
        .overall-result.fail {
            background-color: #c0392b;
        }
        /* Tables */
        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 30px;
        }
        th, td {
            padding: 15px;
            text-align: left;
        }
        th {
            background-color: #2980b9;
            color: #fff;
            cursor: pointer;
        }
        tr:nth-child(even) {
            background-color: #ecf0f1;
        }
        tr:hover {
            background-color: #d0d7de;
        }
        /* Tag Details */
        .tag-details {
            font-size: 0.9em;
            color: #7f8c8d;
            margin-top: 5px;
        }
        /* Error List */
        .error-list {
            background-color: #ffecec;
            border-left: 5px solid #e74c3c;
            padding: 20px;
            border-radius: 5px;
        }
        .error-list ul {
            list-style-type: disc;
            padding-left: 20px;
        }
        .error-list li {
            margin-bottom: 10px;
        }
        /* Responsive Design */
        @media (max-width: 768px) {
            th, td {
                padding: 10px;
            }
            .overall-result {
                font-size: 1em;
            }
            .tag-details {
                font-size: 0.8em;
            }
        }
        /* Custom Search Box Styling */
        .dataTables_filter {
            float: right;
            text-align: right;
        }
        .dataTables_filter input {
            padding: 5px;
            border-radius: 4px;
            border: 1px solid #ccc;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Test Report: {{.SequenceName}}</h1>
        <p><strong>Test ID:</strong> {{.TestID}}</p>
        <p><strong>Timestamp:</strong> {{.Timestamp}}</p>
        
        <div class="overall-result {{if .OverallPassFail}}pass{{else}}fail{{end}}">
            Overall Result: {{if .OverallPassFail}}PASS{{else}}FAIL{{end}}
        </div>
    
        <h2>Tag Submissions</h2>
        <table id="tagsTable">
            <thead>
                <tr>
                    <th>Tag ID</th>
                    <th>Description</th>
                    <th>Comparison</th>
                    <th>Submitted Value</th>
                    <th>Result</th>
                </tr>
            </thead>
            <tbody>
                {{range .TagSubmissions}}
                <tr>
                    <td>{{.TagID}}</td>
                    <td>
                        {{.Tag.Description}}
                        <div class="tag-details">
                            Unit: {{.Tag.Unit}}
                        </div>
                    </td>
                    <td>
                        {{.ComparisonDisplay}}
                    </td>
                    <td>{{.Value}}</td>
                    <td>
                        {{if .IsPassing}}
                            <span style="color: #27ae60; font-weight: bold;">PASS</span>
                        {{else}}
                            <span style="color: #c0392b; font-weight: bold;">FAIL</span>
                        {{end}}
                    </td>
                </tr>
                {{end}}
            </tbody>
        </table>
    
        {{if .ErrorSubmissions}}
        <h2>Errors</h2>
        <div class="error-list">
            <ul>
            {{range .ErrorSubmissions}}
                <li>{{.}}</li>
            {{end}}
            </ul>
        </div>
        {{end}}
    </div>
    
    <!-- jQuery -->
    <script src="https://code.jquery.com/jquery-3.7.0.min.js"></script>
    <!-- DataTables JS -->
    <script src="https://cdn.datatables.net/1.13.6/js/jquery.dataTables.min.js"></script>
    <script>
        $(document).ready(function() {
            $('#tagsTable').DataTable({
                "paging": true,
                "searching": true,
                "ordering": true,
                "order": [],
                "columnDefs": [
                    { "orderable": true, "targets": "_all" }
                ],
                "language": {
                    "search": "Filter records:"
                }
            });
        });
    </script>
</body>
</html>
`
