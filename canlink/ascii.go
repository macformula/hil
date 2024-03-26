package canlink

import (
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_asc = ".asc"
)

// Asc stores the information required to create each ASCII file
type Asc struct {
	l    *zap.Logger
	file *os.File
}

// NewAsc returns a new Asc with its own ASCII file
func NewAsc(l *zap.Logger, dir string, busName string) (*Asc, error) {
	f, err := createFile(l, dir, busName, _asc)
	if err != nil {
		l.Error(err.Error())
		return nil, errors.Wrap(err, "error creating ascii file")
	}

	asc := &Asc{
		file: f,
		l:    l.Named("ascii_file"),
	}

	return asc, nil
}

// dumpToFile takes a CAN frame and writes it to an ASCII file
func (a *Asc) dumpToFile(s []TimestampedFrame) error {
	a.l.Info("ASCII: Entered dumpToFile")

	for _, value := range s {
		_, err := a.file.WriteString(a.parseFrame(value) + "\n")
		if err != nil {
			return errors.Wrap(err, "write string to file")
		}
	}

	return nil
}

// parseFrame concatenates the frame components in a standardized format
func (a *Asc) parseFrame(data TimestampedFrame) string {
	var builder strings.Builder

	_, err := builder.WriteString(data.Time.Format(_messageTimeFormat))
	if err != nil {
		a.l.Error(err.Error())
	}

	_, err = builder.WriteString(" " + strconv.FormatUint(uint64(data.Frame.ID), _decimal))
	if err != nil {
		a.l.Error(err.Error())
	}

	_, err = builder.WriteString(" Rx")
	if err != nil {
		a.l.Error(err.Error())
	}

	_, err = builder.WriteString(" " + strconv.FormatUint(uint64(data.Frame.Length), _decimal))
	if err != nil {
		a.l.Error(err.Error())
	}

	for i := uint8(0); i < data.Frame.Length; i++ {
		builder.WriteString(" " + strconv.FormatUint(uint64(data.Frame.Data[i]), _hex))
		if err != nil {
			a.l.Error(err.Error())
		}
	}

	return builder.String()
}
