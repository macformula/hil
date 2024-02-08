package canlink

import (
	"encoding/json"
	"fmt"
	"github.com/foxglove/mcap/go/mcap"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

type Mcap struct {
	suffix     string
	dir        string
	busName    string
	cachedData []string
	l          *zap.Logger
}

func NewMcap(suffix string, dir string, busName string, cachedData []string, l *zap.Logger) *Mcap {
	mcap := &Mcap{
		suffix:     suffix,
		dir:        dir,
		busName:    busName,
		cachedData: cachedData,
		l:          l.Named(_loggerName),
	}

	return mcap
}

func (m *Mcap) dumpToFile(file *os.File) error {
	//for _, value := range m.cachedData {
	//	_, err := file.WriteString(value + "\n")
	//	if err != nil {
	//		return errors.Wrap(err, "write string to file")
	//	}
	//}

	w, _ := mcap.NewWriter(file, &mcap.WriterOptions{
		Chunked: false,
	})

	err := w.WriteHeader(&mcap.Header{})
	if err != nil {
		panic("FAILED")
	}

	err = w.WriteSchema(&mcap.Schema{
		ID:       1,        //change
		Name:     "schema", //change
		Encoding: "jsonschema",
		Data:     []byte(`{"type":"object"}`),
	})
	if err != nil {
		panic("FAILED")
	}

	err = w.WriteChannel(&mcap.Channel{
		ID:              1, //change
		Topic:           m.busName,
		MessageEncoding: "json",
		SchemaID:        0, //change
		Metadata: map[string]string{
			"callerid": "100", // cspell:disable-line
		},
	})
	if err != nil {
		panic("FAILED")
	}

	for _, value := range m.cachedData {

		message, err := json.Marshal(value)
		if err != nil {
			fmt.Println("Error marshalling message data into json format:", err)
			//return (ask about)
		}

		t := uint64(time.Now().Nanosecond())

		err = w.WriteMessage(&mcap.Message{
			ChannelID: 1, //change
			//Sequence:    0,	(removed since it wasn't included in one of the message loops)
			LogTime:     t,
			PublishTime: t, //should publishtime be the same as logtime? is it even needed?
			Data:        message,
		})
		if err != nil {
			panic("FAILED")
		}
	}

	return nil
}

func (m *Mcap) getFile() (*os.File, error) {
	var file *os.File
	var builder strings.Builder

	//appends the contents of m.dir to b's buffer
	//a buffer is a region of memory to temporarily store data while it is being moved
	_, err := builder.WriteString(m.dir + "/")
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add directory to filepath")
	}

	_, err = builder.WriteString(m.busName + "_")
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add bus name to file name")
	}

	_, err = builder.WriteString(time.Now().Format("15.04.05") + "_")
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add date to file name")
	}

	_, err = builder.WriteString(time.Now().Format("2006.01.02"))
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add time to file name")
	}

	_, err = builder.WriteString(m.suffix)
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add file type to file name")
	}

	//if file does not exist, creates the named file
	//if file exists, it is truncated (removes content from file)
	file, err = os.Create(builder.String())
	if err != nil {
		return &os.File{}, errors.Wrap(err, "create file")
	}

	return file, nil

}
