package results

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Set this to false to keep the test folders for inspection
// WARNING: setting this to false may cause failing tests. This is because from test to test,
// it is assumed that certain files do not exist anymore.
var deleteTestFolders = true

type testSetup struct {
	tempDir         string
	tagsFile        string
	historicTestsFP string
	reportsDir      string
	ra              *ResultAccumulator
}

func setupTest(t *testing.T) testSetup {
	currentDir, err := os.Getwd()
	require.NoError(t, err)
	tempDir := filepath.Join(currentDir, "resultaccumulator_test_"+time.Now().Format("2006-01-02_15-04-05"))
	err = os.MkdirAll(tempDir, 0755)
	require.NoError(t, err)

	t.Logf("Test files created in: %s", tempDir)

	if deleteTestFolders {
		t.Cleanup(func() { os.RemoveAll(tempDir) })
	}

	tagsFile := filepath.Join(currentDir, "testtags.yaml")
	historicTestsFP := filepath.Join(tempDir, "historic_tests.json")
	reportsDir := filepath.Join(tempDir, "reports")

	// Create the reports directory
	err = os.MkdirAll(reportsDir, 0755)
	require.NoError(t, err)

	htmlGenerator := NewHtmlReportGenerator()
	ra := NewResultAccumulator(zap.NewNop(), tagsFile, historicTestsFP, reportsDir, htmlGenerator)

	return testSetup{
		tempDir:         tempDir,
		tagsFile:        tagsFile,
		historicTestsFP: historicTestsFP,
		reportsDir:      reportsDir,
		ra:              ra,
	}
}

func TestResultAccumulatorOpen(t *testing.T) {
	setup := setupTest(t)
	err := setup.ra.Open(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, setup.ra.tagDB)
}

func TestResultAccumulatorSubmitTag(t *testing.T) {
	setup := setupTest(t)
	ctx := context.Background()

	err := setup.ra.Open(ctx)
	require.NoError(t, err)

	testCases := []struct {
		tagID    string
		value    any
		expected bool
	}{
		{"numericGt", 15, true},
		{"numericGt", 5, false},
		{"numericLt", 5, true},
		{"numericLt", 15, false},
		{"numericGe", 20, true},
		{"numericGe", 25, true},
		{"numericGe", 15, false},
		{"numericLe", 20, true},
		{"numericLe", 15, true},
		{"numericLe", 25, false},
		{"numericEq", 15, true},
		{"numericEq", 14, false},
		{"numericGele", 15, true},
		{"numericGele", 10, true},
		{"numericGele", 20, true},
		{"numericGele", 5, false},
		{"numericGele", 25, false},
		{"numericGtlt", 15, true},
		{"numericGtlt", 10, false},
		{"numericGtlt", 20, false},
		{"stringEq", "expected", true},
		{"stringEq", "unexpected", false},
		{"boolEq", true, true},
		{"boolEq", false, false},
		{"logTag", "any value", true},
		{"logNumber", 10, true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%v", tc.tagID, tc.value), func(t *testing.T) {
			result, err := setup.ra.SubmitTag(ctx, tc.tagID, tc.value)
			if err != nil {
				t.Log("error for tag submission (tagId, value):", tc.tagID, tc.value)
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}

	// Test submitting an unknown tag
	_, err = setup.ra.SubmitTag(ctx, "unknownTag", 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tag not found")

	// Test submitting an invalid value type
	_, err = setup.ra.SubmitTag(ctx, "numericGt", "not a number")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to validate tag")
}

func TestResultAccumulatorSubmitError(t *testing.T) {
	setup := setupTest(t)
	err := setup.ra.Open(context.Background())
	require.NoError(t, err)

	err = setup.ra.SubmitError(context.Background(), assert.AnError)
	assert.NoError(t, err)
	assert.False(t, setup.ra.allTagsPassing)
	assert.Len(t, setup.ra.errorSubmissions, 1)
}

func TestResultAccumulatorCompleteTest(t *testing.T) {
	setup := setupTest(t)
	err := setup.ra.Open(context.Background())
	require.NoError(t, err)

	ctx := context.Background()
	testID := uuid.New()
	sequenceName := "TestSequence"

	_, err = setup.ra.SubmitTag(ctx, "numericGt", 15)
	require.NoError(t, err)

	err = setup.ra.SubmitError(ctx, assert.AnError)
	require.NoError(t, err)

	overallPassFail, err := setup.ra.CompleteTest(ctx, testID, sequenceName)
	assert.NoError(t, err)
	assert.False(t, overallPassFail)

	t.Run("HistoricTestsFile", func(t *testing.T) {
		historicTests, err := setup.ra.loadExistingTests()
		assert.NoError(t, err)
		assert.Len(t, historicTests, 1)
		assert.Equal(t, testID.String(), historicTests[0]["testId"])
		assert.Equal(t, sequenceName, historicTests[0]["sequenceName"])
		assert.Equal(t, false, historicTests[0]["testPassed"])
	})

	t.Run("HtmlReportGeneration", func(t *testing.T) {
		htmlReportPath := filepath.Join(setup.reportsDir, fmt.Sprintf("report_%s_%s.html", sequenceName, testID.String()))
		_, err := os.Stat(htmlReportPath)
		assert.NoError(t, err, "HTML report should exist")

		htmlContent, err := os.ReadFile(htmlReportPath)
		assert.NoError(t, err)
		assert.Contains(t, string(htmlContent), "Test Report: TestSequence")
		assert.Contains(t, string(htmlContent), testID.String())
		assert.Contains(t, string(htmlContent), "Overall Result: FAIL")
		assert.Contains(t, string(htmlContent), "Numeric greater than test")
		assert.Contains(t, string(htmlContent), "15")
		assert.Contains(t, string(htmlContent), "PASS")
	})

	t.Run("ResetSubmissions", func(t *testing.T) {
		assert.Empty(t, setup.ra.tagSubmissions)
		assert.Empty(t, setup.ra.errorSubmissions)
		assert.True(t, setup.ra.allTagsPassing)
	})
}

func TestResultAccumulatorEdgeCases(t *testing.T) {
	t.Run("NonExistentTagsFile", func(t *testing.T) {
		setup := setupTest(t)
		ra := NewResultAccumulator(zap.NewNop(), "non_existent_file.yaml", setup.historicTestsFP, setup.reportsDir)
		err := ra.Open(context.Background())
		assert.Error(t, err)
	})

	t.Run("InvalidTagsFile", func(t *testing.T) {
		setup := setupTest(t)

		// Create a new invalid tags file
		invalidTagsFile := filepath.Join(setup.tempDir, "invalid_tags.yaml")
		err := os.WriteFile(invalidTagsFile, []byte("invalid yaml"), 0644)
		require.NoError(t, err)

		ra := NewResultAccumulator(zap.NewNop(), invalidTagsFile, setup.historicTestsFP, setup.reportsDir)
		err = ra.Open(context.Background())
		assert.Error(t, err)
	})

	t.Run("CompleteTestWithExistingHistoricTests", func(t *testing.T) {
		setup := setupTest(t)
		existingTests := []map[string]interface{}{
			{
				"testId":       "existing-test-id",
				"sequenceName": "ExistingSequence",
				"testPassed":   true,
				"date":         "2023-04-20",
				"time":         "12-00-00",
			},
		}
		existingTestsData, err := json.Marshal(existingTests)
		require.NoError(t, err)
		err = os.WriteFile(setup.historicTestsFP, existingTestsData, 0644)
		require.NoError(t, err)

		err = setup.ra.Open(context.Background())
		require.NoError(t, err)

		ctx := context.Background()
		testID := uuid.New()
		_, err = setup.ra.SubmitTag(ctx, "numericGt", 15)
		require.NoError(t, err)

		overallPassFail, err := setup.ra.CompleteTest(ctx, testID, "NewSequence")
		assert.NoError(t, err)
		assert.True(t, overallPassFail)

		historicTests, err := setup.ra.loadExistingTests()
		assert.NoError(t, err)
		assert.Len(t, historicTests, 2)
		assert.Equal(t, testID.String(), historicTests[0]["testId"])
		assert.Equal(t, "existing-test-id", historicTests[1]["testId"])
	})
}

func TestResultAccumulatorSubmitTagAndCompleteTestPass(t *testing.T) {
	setup := setupTest(t)
	ctx := context.Background()

	err := setup.ra.Open(ctx)
	require.NoError(t, err)

	testCases := []struct {
		tagID    string
		value    any
		expected bool
	}{
		{"numericGt", 15, true},
		{"numericLt", 5, true},
		{"numericGe", 20, true},
		{"numericLe", 15, true},
		{"numericEq", 15, true},
		{"numericGele", 15, true},
		{"numericGtlt", 15, true},
		{"stringEq", "expected", true},
		{"boolEq", true, true},
		{"logTag", "any value", true},
		{"logNumber", 10, true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%v", tc.tagID, tc.value), func(t *testing.T) {
			result, err := setup.ra.SubmitTag(ctx, tc.tagID, tc.value)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}

	// Complete the test
	testID := uuid.New()
	sequenceName := "TestSequencePass"

	overallPassFail, err := setup.ra.CompleteTest(ctx, testID, sequenceName)
	assert.NoError(t, err)
	assert.True(t, overallPassFail)

	t.Run("HistoricTestsFile", func(t *testing.T) {
		historicTests, err := setup.ra.loadExistingTests()
		assert.NoError(t, err)
		assert.Len(t, historicTests, 1)
		assert.Equal(t, testID.String(), historicTests[0]["testId"])
		assert.Equal(t, sequenceName, historicTests[0]["sequenceName"])
		assert.Equal(t, true, historicTests[0]["testPassed"])
	})

	t.Run("HtmlReportGeneration", func(t *testing.T) {
		htmlReportPath := filepath.Join(setup.reportsDir, fmt.Sprintf("report_%s_%s.html", sequenceName, testID.String()))
		_, err := os.Stat(htmlReportPath)
		assert.NoError(t, err, "HTML report should exist")

		htmlContent, err := os.ReadFile(htmlReportPath)
		assert.NoError(t, err)
		assert.Contains(t, string(htmlContent), "Test Report: TestSequencePass")
		assert.Contains(t, string(htmlContent), testID.String())
		assert.Contains(t, string(htmlContent), "Overall Result: PASS")
	})

	t.Run("ResetSubmissions", func(t *testing.T) {
		assert.Empty(t, setup.ra.tagSubmissions)
		assert.Empty(t, setup.ra.errorSubmissions)
		assert.True(t, setup.ra.allTagsPassing)
	})
}

func TestResultAccumulatorSubmitTagAndCompleteTestFail(t *testing.T) {
	setup := setupTest(t)
	ctx := context.Background()

	err := setup.ra.Open(ctx)
	require.NoError(t, err)

	testCases := []struct {
		tagID    string
		value    any
		expected bool
	}{
		{"numericGt", 15, true},
		{"numericLt", 15, false},
		{"numericGe", 15, false},
		{"numericLe", 15, true},
		{"numericEq", 14, false},
		{"numericGele", 25, false},
		{"numericGtlt", 20, false},
		{"stringEq", "unexpected", false},
		{"boolEq", false, false},
		{"logTag", "any value", true},
		{"logNumber", 10, true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%v", tc.tagID, tc.value), func(t *testing.T) {
			result, err := setup.ra.SubmitTag(ctx, tc.tagID, tc.value)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}

	// Submit an error
	err = setup.ra.SubmitError(ctx, assert.AnError)
	assert.NoError(t, err)

	// Complete the test
	testID := uuid.New()
	sequenceName := "TestSequenceFail"

	overallPassFail, err := setup.ra.CompleteTest(ctx, testID, sequenceName)
	assert.NoError(t, err)
	assert.False(t, overallPassFail)

	t.Run("HistoricTestsFile", func(t *testing.T) {
		historicTests, err := setup.ra.loadExistingTests()
		assert.NoError(t, err)
		assert.Len(t, historicTests, 1)
		assert.Equal(t, testID.String(), historicTests[0]["testId"])
		assert.Equal(t, sequenceName, historicTests[0]["sequenceName"])
		assert.Equal(t, false, historicTests[0]["testPassed"])
	})

	t.Run("HtmlReportGeneration", func(t *testing.T) {
		htmlReportPath := filepath.Join(setup.reportsDir, fmt.Sprintf("report_%s_%s.html", sequenceName, testID.String()))
		_, err := os.Stat(htmlReportPath)
		assert.NoError(t, err, "HTML report should exist")

		htmlContent, err := os.ReadFile(htmlReportPath)
		assert.NoError(t, err)
		assert.Contains(t, string(htmlContent), "Test Report: TestSequenceFail")
		assert.Contains(t, string(htmlContent), testID.String())
		assert.Contains(t, string(htmlContent), "Overall Result: FAIL")
	})

	t.Run("ResetSubmissions", func(t *testing.T) {
		assert.Empty(t, setup.ra.tagSubmissions)
		assert.Empty(t, setup.ra.errorSubmissions)
		assert.True(t, setup.ra.allTagsPassing)
	})
}
