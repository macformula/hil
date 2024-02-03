package canlink

import (
	"go.uber.org/zap"
	"os"
)

type Mcap struct {
	suffix string
	dir    string
	l      *zap.Logger
}

func NewMcap(suffix string, dir string, l *zap.Logger, ) *Mcap {
	mcap := &Mcap{
		l:      l.Named(_loggerName),
		suffix: suffix,
		dir:    dir,
	}

	return mcap
}

func (m *Mcap) dumpToFile(file *os.File) error {

}

func (m *Mcap) getFile() (*os.File, error) {

}
