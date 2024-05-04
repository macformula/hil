package canlink

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	_idNotInDatabaseErrorIndicator = "ID not in database"
)

// createTraceFile creates an *os.File given information
func createTraceFile(dir, busName, suffix string) (*os.File, error) {
	dateStr := time.Now().Format(_filenameDateFormat)
	timeStr := time.Now().Format(_filenameTimeFormat)

	// remove '.' if included in suffix (done in the following line).
	suffix = strings.TrimLeft(suffix, ".")

	fileName := fmt.Sprintf(
		"%s_%s_%s.%s",
		busName,
		dateStr,
		timeStr,
		suffix,
	)

	filePath := filepath.Join(dir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "create file")
	}

	return file, nil
}

func isIdNotInDatabaseError(err error) bool {
	return err != nil && strings.Contains(err.Error(), _idNotInDatabaseErrorIndicator)
}
