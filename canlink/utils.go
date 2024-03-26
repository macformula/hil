package canlink

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_binary  = 2
	_decimal = 10
	_hex     = 16
)

// createFile creates an *os.File given information
func createFile(l *zap.Logger, dir string, busName string, suffix string) (*os.File, error) {
	l.Info("creating file")
	fileName := fmt.Sprintf("%s_%s_%s%s", busName, time.Now().Format(_filenameDateFormat), time.Now().Format(_filenameTimeFormat), suffix)
	filePath := filepath.Join(dir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "create file")
	}

	return file, nil
}
