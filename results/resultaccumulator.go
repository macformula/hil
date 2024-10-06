package results

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	_loggerName = "result_accumulator"
)

type ResultAccumulator struct {
	l                *zap.Logger
	tagDB            map[string]Tag
	tagSubmissions   map[string]TagSubmission
	errorSubmissions []error
	historicTestsFP  string
	tagsFP           string
	reportsDir       string
	allTagsPassing   bool
	generators       []Generator
}

type TagSubmission struct {
	Tag       Tag
	Value     any
	IsPassing bool
}

func NewResultAccumulator(l *zap.Logger, tagsFP, historicTestsFP, reportsDir string, generators ...Generator) *ResultAccumulator {
	return &ResultAccumulator{
		l:                l.Named(_loggerName),
		tagSubmissions:   make(map[string]TagSubmission),
		errorSubmissions: []error{},
		historicTestsFP:  historicTestsFP,
		tagsFP:           tagsFP,
		reportsDir:       reportsDir,
		allTagsPassing:   true,
		generators:       generators,
	}
}

func (r *ResultAccumulator) Open(ctx context.Context) error {
	if err := r.loadTags(ctx); err != nil {
		return errors.Wrap(err, "failed to load tags")
	}
	return nil
}

func (r *ResultAccumulator) Close() error {
	return nil
}

func (r *ResultAccumulator) loadTags(_ context.Context) error {
	r.tagDB = make(map[string]Tag)

	data, err := os.ReadFile(r.tagsFP)
	if err != nil {
		return errors.Wrapf(err, "failed to read tags file (%v)", r.tagsFP)
	}

	var tagsMap map[string]Tag
	err = yaml.Unmarshal(data, &tagsMap)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal tags YAML")
	}

	for key, tag := range tagsMap {
		// Convert CompOpString to ComparisonOperator
		compOp, err := ComparisonOperatorString(strings.ToLower(tag.CompOpString))
		if err != nil {
			return errors.Wrapf(err, "invalid comparison operator for tag %s", key)
		}
		tag.CompOp = compOp
		r.tagDB[key] = tag
	}

	return nil
}

func (r *ResultAccumulator) SubmitTag(_ context.Context, tagID string, value any) (bool, error) {
	tag, ok := r.tagDB[tagID]
	if !ok {
		return false, errors.Errorf("tag not found: %s", tagID)
	}

	isPassing, err := tag.IsPassing(value)
	if err != nil {
		return false, errors.Wrapf(err, "failed to validate tag %s", tagID)
	}

	r.tagSubmissions[tagID] = TagSubmission{
		Tag:       tag,
		Value:     value,
		IsPassing: isPassing,
	}

	if !isPassing {
		r.allTagsPassing = false
	}

	return isPassing, nil
}

func (r *ResultAccumulator) SubmitError(_ context.Context, err error) error {
	r.errorSubmissions = append(r.errorSubmissions, err)
	r.allTagsPassing = false
	return nil
}

func (r *ResultAccumulator) CompleteTest(_ context.Context, testID uuid.UUID, sequenceName string) (bool, error) {
	dateTime := time.Now().Format("2006-01-02_15-04-05")

	overallPassFail := r.allTagsPassing && len(r.errorSubmissions) == 0

	if err := r.updateHistoricTests(testID.String(), sequenceName, dateTime, overallPassFail); err != nil {
		return false, errors.Wrap(err, "failed to update historic tests")
	}

	for _, generator := range r.generators {
		err := generator.Generate(
			testID,
			sequenceName,
			r.tagSubmissions,
			r.errorSubmissions,
			overallPassFail,
			r.reportsDir,
		)
		if err != nil {
			return false, errors.Wrap(err, "failed to generate report")
		}
	}

	// Reset cached submissions
	r.tagSubmissions = make(map[string]TagSubmission)
	r.errorSubmissions = []error{}
	r.allTagsPassing = true

	return overallPassFail, nil
}

func (r *ResultAccumulator) updateHistoricTests(testID, sequenceName, dt string, testPassed bool) error {
	ymd, hms, _ := strings.Cut(dt, "_")

	newTest := map[string]any{
		"testId":       testID,
		"sequenceName": sequenceName,
		"testPassed":   testPassed,
		"date":         ymd,
		"time":         hms,
	}

	existingTests, err := r.loadExistingTests()
	if err != nil {
		return err
	}

	updatedTests := append([]map[string]any{newTest}, existingTests...)

	return r.saveUpdatedTests(updatedTests)
}

func (r *ResultAccumulator) loadExistingTests() ([]map[string]any, error) {
	data, err := os.ReadFile(r.historicTestsFP)
	if err != nil {
		if os.IsNotExist(err) {
			return []map[string]any{}, nil
		}
		return nil, errors.Wrap(err, "failed to read historic tests file")
	}

	var existingTests []map[string]any
	err = json.Unmarshal(data, &existingTests)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal historic tests")
	}

	return existingTests, nil
}

func (r *ResultAccumulator) saveUpdatedTests(updatedTests []map[string]any) error {
	data, err := json.MarshalIndent(updatedTests, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal updated tests")
	}

	err = os.WriteFile(r.historicTestsFP, data, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write updated tests to file")
	}

	return nil
}
