package canlink

import (
	"go.uber.org/zap"
	"os"
)

type Asc struct {
	suffix string
	dir    string
	l      *zap.Logger
}

func NewAsc(suffix string, dir string, l *zap.Logger) *Asc {
	asc := &Asc{
		l:      l.Named(_loggerName),
		suffix: suffix,
		dir:    dir,
	}

	return asc
}

func (a *Asc) dumpToFile(file *os.File) error {

}

func (a *Asc) getFile() (*os.File, error) {

}
