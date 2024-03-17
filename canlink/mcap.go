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

// defining MCAP properties
type Mcap struct {
	suffix     string
	dir        string
	busName    string
	cachedData *[]string
	l          *zap.Logger
	//map[key]value{}
	createdChannels map[string]int
	c               *CANClient
}

//defining CAN frame properties
//type CANFrame struct {
//	timestamp uint64
//	signalID  int
//	message   []int
//}

// instantiates a new Mcap object
func NewMcap(suffix string, dir string, busName string, cachedData *[]string, l *zap.Logger, c *CANClient) *Mcap {
	mcap := &Mcap{
		suffix:          suffix,
		dir:             dir,
		busName:         busName,
		cachedData:      cachedData,
		l:               l.Named("mcap_logger"),
		createdChannels: make(map[string]int),
		c:               c,
	}

	return mcap
}

//instantiates new CANFrame object
//func NewCANFrame(timestamp uint64, signalID int, message []int) *CANFrame {
//	canframe := &CANFrame{
//		timestamp: timestamp,
//		signalID: signalID,
//		message: message,
//	}
//}

// dumps all contents from cached data to Mcap file
func (m *Mcap) dumpToFile(file *os.File) error {

	//creating new instance of writer
	w, _ := mcap.NewWriter(file, &mcap.WriterOptions{
		Chunked: false,
	})

	//declaring footer
	defer func() {
		err := w.Close()
		m.l.Info("MCAP: Mcap file is closing and writer turned off ")
		if err != nil {
			fmt.Println("Error closing mcap file:", err)
			return
		}
	}()

	//creating header
	err := w.WriteHeader(&mcap.Header{})
	if err != nil {
		m.l.Info("error creating headers")
		panic("FAILED creating header")
	}

	//parsing the signalID
	//holderArray := strings.Fields(m.cachedData[0])
	//parser, err := strconv.ParseUint(holderArray[1], 16, 16)
	//signalID := uint16(parser)
	//if err != nil {
	//	fmt.Println("Error parsing signal ID into uint16:", err)
	//}

	//creating the schema (only need one)
	m.l.Info("schema created")
	err = w.WriteSchema(&mcap.Schema{
		ID:       1,
		Name:     m.busName, //change
		Encoding: "jsonschema",
		Data:     []byte(`{"type":"object"}`),
	})
	if err != nil {
		m.l.Info("error creating schemas")
		panic("FAILED creating schema")
	}

	//dataSlice stores all the content of cachedData
	dataSlice := *m.cachedData
	//iterates through cachedData (every value is of format: 14:51:53.0772 1574 Rx 1 2)
	for _, value := range dataSlice {
		//separating the string by blank spaces
		tempArray := strings.Fields(value)
		//converting the signalID into uint16
		parser, err := strconv.ParseUint(tempArray[1], 16, 16) //signalID is stored in the second index of array
		signalID := uint16(parser)

		//checking map to see if channel for the signalID is already created
		if m.createdChannels[tempArray[1]] != 1 { //tempArray[1]
			//if map doesn't return 1, channel doesn't exist and must be created
			m.l.Info("channel created")
			//creating channel
			err = w.WriteChannel(&mcap.Channel{
				ID:    signalID,
				Topic: "frame-" + tempArray[1],
				//use canclient to find signal name and message name
				MessageEncoding: "json",
				SchemaID:        1,
				Metadata: map[string]string{
					"callerid": "100", // cspell:disable-line
				},
			})
			if err != nil {
				m.l.Info("error creating channels")
				panic("FAILED creating channels")
			}
			//adding created channel into the map using signalID as key and setting value as 1
			m.createdChannels[tempArray[1]] = 1
		}

		//taking message from cachedData and marshaling (will be replaced with canclient implementation)
		message, err := json.Marshal(tempArray[4]) // Might be an issue with how this is being processed
		if err != nil {
			fmt.Println("Error marshalling message data into json format:", err)
			//return (ask about)
		}

		//parsing time from the cached data to set as message timestamp
		parsedTime, err := time.Parse("15:04:05.0000", tempArray[0])
		tSeconds := uint64(parsedTime.Second()) * 1000000000
		t := uint64(parsedTime.Nanosecond()) + tSeconds
		m.l.Info("Parsed time: " + strconv.FormatUint(t, 10) + "\n")
		//m.l.Info("Parsed time in seconds" + strconv.FormatUint(tSeconds, 10))

		//creating messages

		err = w.WriteMessage(&mcap.Message{
			ChannelID: signalID,
			//Sequence:    0,	(removed since it wasn't included in one of the message loops)
			LogTime:     t,
			PublishTime: t, //should publishtime be the same as logtime? is it even needed?
			Data:        message,
		})
		if err != nil {
			logmsg := fmt.Sprintf("Failed to write MCAP message: ChannelID: %d, Logtime: %d", signalID, t)
			m.l.Info(logmsg)
			panic("FAILED writing message")

		} else {
			m.l.Info("message created")
		}
	}

	return nil
}

// getting MCAP file that will be used to dump contents of cachedData to
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

	_, err = builder.WriteString(time.Now().Format("2006.01.02") + "_")
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add time to file name")
	}

	_, err = builder.WriteString(time.Now().Format("15.04.05"))
	if err != nil {
		return &os.File{}, errors.Wrap(err, "add date to file name")
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

//func (m *Mcap) checkChannel(value string) (int) {
//	tempArray := strings.Fields(value)
//	parser, err := strconv.ParseUint(tempArray[1], 16, 16)
//	signalID := uint16(parser)
//
//	if m.createdChannels[tempArray[1]] != 1 { //tempArray[1]
//		err = w.WriteChannel(&mcap.Channel{
//			ID:              signalID,
//			Topic:           tempArray[1], //contactorfeedback.positive
//			MessageEncoding: "json",
//			SchemaID:        1, //change
//			Metadata: map[string]string{
//				"callerid": "100", // cspell:disable-line
//			},
//		})
//		if err != nil {
//			panic("FAILED")
//		}
//		m.createdChannels[tempArray[1]] = 1
//	}
//}
