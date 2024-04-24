package canlink

import (
	"encoding/json"
	"fmt"
	"github.com/macformula/hil/utils"
	"os"
	"time"

	"github.com/foxglove/mcap/go/mcap"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_mcapLoggerName   = "mcap_file"
	_mcap             = "mcap"
	_json             = "json"
	_schemaName       = "MacFormula"
	_schemaEncoding   = "jsonschema"
	_schemaData       = `{"type":"object"}`
	_schemaId         = 1
	_initialChannelId = 1
)

type (
	// Mcap writes records to "channels". Each CAN signal gets mapped to a channelId.
	channelId = uint16
)

// Mcap stores the information required to create each MCAP file.
type Mcap struct {
	l         *zap.Logger
	file      *os.File
	canClient *CanClient

	// Maps (msg name + signal name) to channelId
	signalChannelIds map[string]channelId
	nextChannelId    channelId

	schema *mcap.Schema
	header *mcap.Header
	writer *mcap.Writer
}

// NewMcap returns a new Mcap with its own MCAP file.
func NewMcap(c *CanClient, l *zap.Logger) *Mcap {
	return &Mcap{
		l:                l.Named(_mcapLoggerName),
		canClient:        c,
		signalChannelIds: map[string]channelId{},
		nextChannelId:    _initialChannelId,
		header:           &mcap.Header{},
		schema: &mcap.Schema{
			ID:       _schemaId,
			Name:     _schemaName,
			Encoding: _schemaEncoding,
			Data:     []byte(_schemaData),
		},
	}
}

// dumpToFile takes a CAN frame and writes it to an MCAP file.
func (m *Mcap) dumpToFile(frames []TimestampedFrame, traceDirectory, busName string) error {
	var (
		oneTimeErr        = utils.ResettableError{}
		mcapWriterOptions = &mcap.WriterOptions{
			Chunked: false,
		}
	)

	// Creating trace file.
	file, err := createTraceFile(traceDirectory, busName, _mcap)
	if err != nil {
		m.l.Error("failed to create trace file", zap.Error(err))

		return errors.Wrapf(err, "create trace file (%s)", m.file.Name())
	}

	// Creating new instance of writer. Must do this every time.
	m.writer, err = mcap.NewWriter(file, mcapWriterOptions)
	if err != nil {
		m.l.Error("failed creating mcap writer", zap.Error(err))

		return errors.Wrap(err, "error creating writer")
	}

	// Writing headers.
	err = m.writer.WriteHeader(m.header)
	if err != nil {
		m.l.Error("failed to write header", zap.Error(err))

		oneTimeErr.Set(errors.Wrap(err, "error creating headers"))
	}

	// Writing schema.
	err = m.writer.WriteSchema(m.schema)
	if err != nil {
		m.l.Error("failed to write schemas", zap.Error(err))

		oneTimeErr.Set(errors.Wrap(err, "error creating schema"))
	}

	// Write all timestamped frames to file.
	for i := 0; i < len(frames); i++ {
		err = m.writeCanMessage(&frames[i])
		if err != nil {
			m.l.Info("error writing can message")

			oneTimeErr.Set(errors.Wrap(err, "error writing message"))
		}
	}

	// Closing writer.
	err = m.writer.Close()
	if err != nil {
		m.l.Error("failed to close writer", zap.Error(err))

		return errors.Wrap(err, "close writer")
	}

	return nil
}

// writeCanMessage loops through all the signals and calls writeCanSignal.
func (m *Mcap) writeCanMessage(timestampedFrame *TimestampedFrame) error {
	frame, err := m.canClient.UnmarshalFrame(timestampedFrame.Frame)
	if isIdNotInDatabaseError(err) {
		// This frame id does not exist in the dbc
		return nil
	} else if err != nil {
		m.l.Error("failed to unmarshal frame", zap.Error(err))

		return errors.Wrap(err, "unmarshal frame")
	}

	msg := frame.Descriptor()
	signals := msg.Signals
	msgData := frame.Frame().Data

	for _, sig := range signals {
		sigValue := sig.UnmarshalPhysical(msgData)

		err = m.writeCanSignal(msg.Name, sig.Name, sigValue, timestampedFrame.Time)
		if err != nil {
			m.l.Error("failed to write can signal", zap.Error(err))

			return errors.Wrap(err, "write can signal")
		}
	}

	return nil
}

// writeCanSignal creates the Mcap channel if it doesn't exist and then writes the CAN Message to it.
func (m *Mcap) writeCanSignal(msgName, sigName string, signalValue float64, receivedTime time.Time) error {
	chanId, err := m.getChannelId(msgName, sigName)

	signalValueJson, err := json.Marshal(signalValue)
	if err != nil {
		m.l.Info("failed to marshal message into json", zap.Error(err))

		return errors.Wrap(err, "json marshal")
	}

	timeReceivedUnix := uint64(receivedTime.UnixNano())
	mcapMsg := mcap.Message{
		ChannelID:   chanId,
		LogTime:     timeReceivedUnix, // Time at which the message was recorded.
		PublishTime: timeReceivedUnix, // Time at which the message was published. Not available, must set to log time.
		Data:        signalValueJson,
	}

	err = m.writer.WriteMessage(&mcapMsg)
	if err != nil {
		m.l.Error(
			"failed to write mcap message",
			zap.Uint16("channel_id", chanId),
			zap.Uint64("receive_time", timeReceivedUnix),
			zap.Error(err),
		)

		return errors.Wrap(err, "write message")
	}

	return nil
}

func (m *Mcap) getChannelId(msgName, sigName string) (channelId, error) {
	chanId, ok := m.signalChannelIds[msgName+sigName]
	if !ok {
		chanId = m.nextChannelId
		m.signalChannelIds[sigName] = chanId
		// Increment channel id for next added signal.
		m.nextChannelId++

		// Channel does not yet exist, must assign a new id, then create a channel.
		m.l.Debug("creating channel",
			zap.String("message_name", msgName),
			zap.String("signal_name", sigName),
			zap.Uint16("channel_id", chanId),
		)

		err := m.createChannel(msgName, sigName, chanId)
		if err != nil {
			return 0, errors.Wrap(err, "create channel")
		}
	}

	return chanId, nil
}

func (m *Mcap) createChannel(msgName, sigName string, chanId channelId) error {
	mcapChannel := mcap.Channel{
		ID:              chanId,
		Topic:           fmt.Sprintf("%s.%s", msgName, sigName),
		MessageEncoding: _json,
		SchemaID:        _schemaId,
	}

	err := m.writer.WriteChannel(&mcapChannel)
	if err != nil {
		m.l.Error("failed to write channel", zap.Error(err))

		return errors.Wrap(err, "write channel")
	}

	return nil
}
