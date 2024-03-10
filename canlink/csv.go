package canlink

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

type Csv struct {
	suffix     string
	dir        string
	busName    string
	cachedData *[]string
	l          *zap.Logger
	c          *CANClient
}

func NewCsv(suffix string, dir string, busName string, cachedData *[]string, l *zap.Logger, c *CANClient) *Csv {
	csv := &Csv{
		suffix:     suffix,
		dir:        dir,
		busName:    busName,
		cachedData: cachedData,
		l:          l.Named("CSV_logger"),
		c:          c,
	}

	return csv
}

func (c *Csv) dumpToFile(file *os.File) error {
	//<<<<<<< HEAD
	dataSlice := *c.cachedData

	for _, value := range dataSlice {
		holderArray := strings.Fields(value)

		for _, word := range holderArray {
			_, err := file.WriteString(word + ",")
			if err != nil {
				return errors.Wrap(err, "write string to file")
			}
		}
		_, err := file.WriteString("\n")
		//=======
		//	dataSlice := *c.cachedData
		//	for _, value := range dataSlice {
		//		_, err := file.WriteString(value + ",")
		//>>>>>>> 9c4d69a2299b72a506edb6124e042135b622b671
		if err != nil {
			return errors.Wrap(err, "starting new line in file")
		}
	}

	return nil
}

func (c *Csv) getFile() (*os.File, error) {
	c.l.Info("CSV: entered getFile")
	var file *os.File
	var builder strings.Builder

	_, err := builder.WriteString(c.dir + "/")
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add directory to filepath")
	}

	_, err = builder.WriteString(c.busName + "_")
	if err != nil {
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

	_, err = builder.WriteString(c.suffix)
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add file type to file name")
	}

	file, err = os.Create(builder.String())
	if err != nil {
		return &os.File{}, errors.Wrap(err, "create file")
	}

	return file, nil
}
