package canlink

import (
	"encoding/json"
	"fmt"
	"github.com/foxglove/mcap/go/mcap"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
)

type Mcap struct {
	suffix     string
	dir        string
	busName    string
	cachedData []string
	l          *zap.Logger
	//map[key]value{}
	createdChannels map[string]int
}

func NewMcap(suffix string, dir string, busName string, cachedData []string, l *zap.Logger, createdChannels map[string]int) *Mcap {
	mcap := &Mcap{
		suffix:          suffix,
		dir:             dir,
		busName:         busName,
		cachedData:      cachedData,
		l:               l.Named(_loggerName),
		createdChannels: createdChannels,
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

	//parsing the signalID
	//holderArray := strings.Fields(m.cachedData[0])
	//parser, err := strconv.ParseUint(holderArray[1], 16, 16)
	//signalID := uint16(parser)
	//if err != nil {
	//	fmt.Println("Error parsing signal ID into uint16:", err)
	//}

	//initializing schema (only need one)
	err = w.WriteSchema(&mcap.Schema{
		ID:       1,         //change (cannot be 0)
		Name:     m.busName, //change
		Encoding: "jsonschema",
		Data:     []byte(`{"type":"object"}`),
	})
	if err != nil {
		panic("FAILED")
	}

	for _, value := range m.cachedData {
		tempArray := strings.Fields(value)
		parser, err := strconv.ParseUint(tempArray[1], 16, 16)
		signalID := uint16(parser)

		//creating channels for each signal
		if m.createdChannels[tempArray[1]] != 1 {
			err = w.WriteChannel(&mcap.Channel{
				ID:              signalID,
				Topic:           tempArray[1], //contactorfeedback.positive
				MessageEncoding: "json",
				SchemaID:        1, //change
				Metadata: map[string]string{
					"callerid": "100", // cspell:disable-line
				},
			})
			if err != nil {
				panic("FAILED")
			}
			m.createdChannels[tempArray[1]] = 1
		}

		////creating channels for each message
		//err = w.WriteChannel(&mcap.Channel{
		//	ID:              1,              //change to message ID
		//	Topic:           "message_name", //change to message
		//	MessageEncoding: "json",
		//	SchemaID:        1, //change
		//	//ChannelID: 1,
		//	Metadata: map[string]string{
		//		"callerid": "100", // cspell:disable-line
		//	},
		//})
		//if err != nil {
		//	panic("FAILED")
		//}

		//taking message from the cached data
		message, err := json.Marshal(tempArray[3])
		if err != nil {
			fmt.Println("Error marshalling message data into json format:", err)
			//return (ask about)
		}

		//time should be time of the signal
		//data is in the form: 14:51:53.0772 1574 Rx 1 2
		//t := uint64(time.Now().Nanosecond())

		//parsing time from the cached data
		parsedTime, err := time.Parse("15:04:05.0000", tempArray[0])
		if err != nil {
			fmt.Println("Error parsing string into time format:", err)
			//return (ask about)
		}
		t := uint64(parsedTime.Unix())

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
