package canlink

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

type Asc struct {
	suffix     string
	dir        string
	busName    string
	cachedData []string
	l          *zap.Logger
}

func NewAsc(suffix string, dir string, busName string, cachedData []string, l *zap.Logger) *Asc {
	asc := &Asc{
		suffix:     suffix,
		dir:        dir,
		busName:    busName,
		cachedData: cachedData,
		l:          l.Named("ascii_logger"),
	}

	return asc
}

func (a *Asc) dumpToFile(file *os.File) error {
	for _, value := range a.cachedData {
		_, err := file.WriteString(value + "\n")
		if err != nil {
			return errors.Wrap(err, "write string to file")
		}
	}

	return nil
}

func (a *Asc) getFile() (*os.File, error) {
	var file *os.File
	var builder strings.Builder

	_, err := builder.WriteString(a.dir + "/")
	if err != nil {
		a.l.Error("ASCII: Error adding directory to filepath")
		return &os.File{}, errors.Wrap(err, "add directory to filepath")
	}

	_, err = builder.WriteString(a.busName + "_")
	if err != nil {
		a.l.Error("ASCII: Error adding bus name to filename")

		return &os.File{}, errors.Wrap(err, "add bus name to file name")
	}

	_, err = builder.WriteString(time.Now().Format("2006.01.02") + "_")
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add date to file name")
	}

	_, err = builder.WriteString(time.Now().Format("15.04.05"))
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add time to file name")
	}

	_, err = builder.WriteString(".asc")
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add file type to file name")
	}

	file, err = os.Create(builder.String())
	if err != nil {
		return &os.File{}, errors.Wrap(err, "create file")
	}

	return file, nil
}
