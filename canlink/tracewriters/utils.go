package tracewriters

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.einride.tech/can"
	"github.com/pkg/errors"
)

// createEmptyTraceFile creates an *os.File given information
func createEmptyTraceFile(dir string, busName string, fileSuffix string) (*os.File, error) {
	dateStr := time.Now().Format(_filenameDateFormat)
	timeStr := time.Now().Format(_filenameTimeFormat)

	fileName := fmt.Sprintf(
		"%s_%s_%s.%s",
		busName,
		dateStr,
		timeStr,
		fileSuffix,
	)

	filePath := filepath.Join(dir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "create file")
	}

	return file, nil
}

type TimestampedFrame struct {
	Frame can.Frame
	Time  time.Time
}
