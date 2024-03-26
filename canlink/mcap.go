package canlink

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/foxglove/mcap/go/mcap"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_mcap = ".mcap"
	_json = "json"
)

// Mcap stores the information required to create each MCAP file.
type Mcap struct {
	l              *zap.Logger
	file           *os.File
	signalChannels map[string]uint16
	client         *CANClient
	currIndex      uint16
}

// NewMcap returns a new Mcap with its own MCAP file.
func NewMcap(l *zap.Logger, c *CANClient, dir string, busName string) (*Mcap, error) {
	f, err := createFile(l, dir, busName, _mcap)
	if err != nil {
		e := errors.Wrap(err, "error creating mcap file")
		l.Error(e.Error())
		return nil, e
	}

	m := &Mcap{
		l:              l.Named("mcap_file"),
		client:         c,
		file:           f,
		signalChannels: make(map[string]uint16),
		currIndex:      1,
	}

	return m, nil
}

// dumpToFile takes a CAN frame and writes it to an MCAP file.
func (m *Mcap) dumpToFile(s []TimestampedFrame) error {
	// creating new instance of writer
	w, err := mcap.NewWriter(m.file, &mcap.WriterOptions{
		Chunked: false,
	})
	if err != nil {
		m.l.Info("error creating writer")
		return errors.Wrap(err, "error creating writer")
	}

	defer func() {
		err = m.closeWriter(w)
		if err != nil {
			m.l.Info("error closing writer")
			return
		}
	}()

	err = m.initHeaders(err, w)
	if err != nil {
		m.l.Info("error initializing header")
		return errors.Wrap(err, "error initializing header")
	}

	for _, value := range s {
		err = m.writeMessage(w, value)
		if err != nil {
			m.l.Info("error writing message")
			return errors.Wrap(err, "error writing message")
		}
	}

	return nil
}

// writeMessage loops through all the signals and calls writeSignal.
func (m *Mcap) writeMessage(w *mcap.Writer, value TimestampedFrame) error {
	//checking map to see if channel for the signalID is already created
	frame, err := m.client.UnmarshalFrame(value.Frame)
	if err != nil {
		m.l.Error(err.Error())
		return errors.Wrap(err, "error unmarshalling frame using canclient")
	} else {
		messageName := frame.Descriptor().Name
		for _, signal := range frame.Descriptor().Signals {
			err = m.writeSignal(w, value, messageName, signal.Name, signal.UnmarshalPhysical(value.Frame.Data))
			if err != nil {
				m.l.Info("error writing signal")
				return errors.Wrap(err, "error writing signal")
			}
		}
	}

	return nil
}

// writeSignal creates the Mcap channel if it doesn't exist and then writes the CAN Message to it.
func (m *Mcap) writeSignal(w *mcap.Writer, value TimestampedFrame, messageName string, signalName string, signalValue float64) error {
	signalID := uint16(value.Frame.ID)
	var channelID uint16
	if val, exists := m.signalChannels[strconv.Itoa(int(signalID))+"."+signalName]; exists {
		channelID = val
	} else {
		m.l.Info("channel created")
		channelID = m.currIndex
		m.currIndex++
		m.signalChannels[strconv.Itoa(int(signalID))+"."+signalName] = channelID
		err := w.WriteChannel(&mcap.Channel{
			ID:              channelID,
			Topic:           messageName + "." + signalName,
			MessageEncoding: _json,
			SchemaID:        1,
		})
		if err != nil {
			e := errors.Wrap(err, "error creating channel")
			m.l.Error(e.Error())
			return e
		}
	}

	message, err := json.Marshal(signalValue)
	if err != nil {
		e := errors.Wrap(err, "error marshalling message data into json format")
		m.l.Error(e.Error())
		return e
	}

	t := uint64(value.Time.UnixNano())

	err = w.WriteMessage(&mcap.Message{
		ChannelID:   channelID,
		LogTime:     t, // Time at which the message was recorded.
		PublishTime: t, // Time at which the message was published. If not available, must be set to the log time.
		Data:        message,
	})
	if err != nil {
		m.l.Info(fmt.Sprintf("Failed to write MCAP message: ChannelID: %d, Logtime: %d", signalID, t))
		return errors.Wrap(err, "failed to write MCAP message")
	}

	return nil
}

func (m *Mcap) initHeaders(err error, w *mcap.Writer) error {
	// creating header
	err = w.WriteHeader(&mcap.Header{})
	if err != nil {
		m.l.Info("error creating headers")
		return errors.Wrap(err, "error creating headers")
	}

	// creating the schema
	m.l.Info("schema created")
	err = w.WriteSchema(&mcap.Schema{
		ID:       1,
		Name:     "MacFormula",
		Encoding: "jsonschema",
		Data:     []byte(`{"type":"object"}`),
	})
	if err != nil {
		m.l.Info("error creating schemas")
		return errors.Wrap(err, "error creating schema")
	}

	return err
}

// closeWriter closes the Mcap Writer and write the footer
func (m *Mcap) closeWriter(w *mcap.Writer) error {
	err := w.Close()
	m.l.Info("MCAP: Mcap file is closing and writer turned off ")
	if err != nil {
		e := errors.Wrap(err, "Error closing mcap file")
		m.l.Error(e.Error())
		return e
	}

	return nil
}
