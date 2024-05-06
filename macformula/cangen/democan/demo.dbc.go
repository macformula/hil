// Package democan provides primitives for encoding and decoding demo CAN messages.
//
// Source: temp/democan/demo.dbc
package democan

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"go.einride.tech/can"
	"go.einride.tech/can/pkg/candebug"
	"go.einride.tech/can/pkg/canrunner"
	"go.einride.tech/can/pkg/cantext"
	"go.einride.tech/can/pkg/descriptor"
	"go.einride.tech/can/pkg/generated"
	"go.einride.tech/can/pkg/socketcan"
)

// prevent unused imports
var (
	_ = context.Background
	_ = fmt.Print
	_ = net.Dial
	_ = http.Error
	_ = sync.Mutex{}
	_ = time.Now
	_ = socketcan.Dial
	_ = candebug.ServeMessagesHTTP
	_ = canrunner.Run
)

// Generated code. DO NOT EDIT.
// TempSensorsReader provides read access to a TempSensors message.
type TempSensorsReader interface {
	// Sensor6 returns the physical value of the Sensor6 signal.
	Sensor6() float64
	// Sensor5 returns the physical value of the Sensor5 signal.
	Sensor5() float64
	// Sensor4 returns the physical value of the Sensor4 signal.
	Sensor4() float64
	// Sensor3 returns the physical value of the Sensor3 signal.
	Sensor3() float64
	// Sensor2 returns the physical value of the Sensor2 signal.
	Sensor2() float64
	// Sensor1 returns the physical value of the Sensor1 signal.
	Sensor1() float64
}

// TempSensorsWriter provides write access to a TempSensors message.
type TempSensorsWriter interface {
	// CopyFrom copies all values from TempSensorsReader.
	CopyFrom(TempSensorsReader) *TempSensors
	// SetSensor6 sets the physical value of the Sensor6 signal.
	SetSensor6(float64) *TempSensors
	// SetSensor5 sets the physical value of the Sensor5 signal.
	SetSensor5(float64) *TempSensors
	// SetSensor4 sets the physical value of the Sensor4 signal.
	SetSensor4(float64) *TempSensors
	// SetSensor3 sets the physical value of the Sensor3 signal.
	SetSensor3(float64) *TempSensors
	// SetSensor2 sets the physical value of the Sensor2 signal.
	SetSensor2(float64) *TempSensors
	// SetSensor1 sets the physical value of the Sensor1 signal.
	SetSensor1(float64) *TempSensors
}

type TempSensors struct {
	xxx_Sensor6 int16
	xxx_Sensor5 int16
	xxx_Sensor4 int16
	xxx_Sensor3 int16
	xxx_Sensor2 int16
	xxx_Sensor1 int16
}

func NewTempSensors() *TempSensors {
	m := &TempSensors{}
	m.Reset()
	return m
}

func (m *TempSensors) Reset() {
	m.xxx_Sensor6 = 0
	m.xxx_Sensor5 = 0
	m.xxx_Sensor4 = 0
	m.xxx_Sensor3 = 0
	m.xxx_Sensor2 = 0
	m.xxx_Sensor1 = 0
}

func (m *TempSensors) CopyFrom(o TempSensorsReader) *TempSensors {
	m.SetSensor6(o.Sensor6())
	m.SetSensor5(o.Sensor5())
	m.SetSensor4(o.Sensor4())
	m.SetSensor3(o.Sensor3())
	m.SetSensor2(o.Sensor2())
	m.SetSensor1(o.Sensor1())
	return m
}

// Descriptor returns the TempSensors descriptor.
func (m *TempSensors) Descriptor() *descriptor.Message {
	return Messages().TempSensors.Message
}

// String returns a compact string representation of the message.
func (m *TempSensors) String() string {
	return cantext.MessageString(m)
}

func (m *TempSensors) Sensor6() float64 {
	return Messages().TempSensors.Sensor6.ToPhysical(float64(m.xxx_Sensor6))
}

func (m *TempSensors) SetSensor6(v float64) *TempSensors {
	m.xxx_Sensor6 = int16(Messages().TempSensors.Sensor6.FromPhysical(v))
	return m
}

func (m *TempSensors) Sensor5() float64 {
	return Messages().TempSensors.Sensor5.ToPhysical(float64(m.xxx_Sensor5))
}

func (m *TempSensors) SetSensor5(v float64) *TempSensors {
	m.xxx_Sensor5 = int16(Messages().TempSensors.Sensor5.FromPhysical(v))
	return m
}

func (m *TempSensors) Sensor4() float64 {
	return Messages().TempSensors.Sensor4.ToPhysical(float64(m.xxx_Sensor4))
}

func (m *TempSensors) SetSensor4(v float64) *TempSensors {
	m.xxx_Sensor4 = int16(Messages().TempSensors.Sensor4.FromPhysical(v))
	return m
}

func (m *TempSensors) Sensor3() float64 {
	return Messages().TempSensors.Sensor3.ToPhysical(float64(m.xxx_Sensor3))
}

func (m *TempSensors) SetSensor3(v float64) *TempSensors {
	m.xxx_Sensor3 = int16(Messages().TempSensors.Sensor3.FromPhysical(v))
	return m
}

func (m *TempSensors) Sensor2() float64 {
	return Messages().TempSensors.Sensor2.ToPhysical(float64(m.xxx_Sensor2))
}

func (m *TempSensors) SetSensor2(v float64) *TempSensors {
	m.xxx_Sensor2 = int16(Messages().TempSensors.Sensor2.FromPhysical(v))
	return m
}

func (m *TempSensors) Sensor1() float64 {
	return Messages().TempSensors.Sensor1.ToPhysical(float64(m.xxx_Sensor1))
}

func (m *TempSensors) SetSensor1(v float64) *TempSensors {
	m.xxx_Sensor1 = int16(Messages().TempSensors.Sensor1.FromPhysical(v))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *TempSensors) Frame() can.Frame {
	md := Messages().TempSensors
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.Sensor6.MarshalSigned(&f.Data, int64(m.xxx_Sensor6))
	md.Sensor5.MarshalSigned(&f.Data, int64(m.xxx_Sensor5))
	md.Sensor4.MarshalSigned(&f.Data, int64(m.xxx_Sensor4))
	md.Sensor3.MarshalSigned(&f.Data, int64(m.xxx_Sensor3))
	md.Sensor2.MarshalSigned(&f.Data, int64(m.xxx_Sensor2))
	md.Sensor1.MarshalSigned(&f.Data, int64(m.xxx_Sensor1))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *TempSensors) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *TempSensors) UnmarshalFrame(f can.Frame) error {
	md := Messages().TempSensors
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal TempSensors: expects ID 940 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal TempSensors: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal TempSensors: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal TempSensors: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_Sensor6 = int16(md.Sensor6.UnmarshalSigned(f.Data))
	m.xxx_Sensor5 = int16(md.Sensor5.UnmarshalSigned(f.Data))
	m.xxx_Sensor4 = int16(md.Sensor4.UnmarshalSigned(f.Data))
	m.xxx_Sensor3 = int16(md.Sensor3.UnmarshalSigned(f.Data))
	m.xxx_Sensor2 = int16(md.Sensor2.UnmarshalSigned(f.Data))
	m.xxx_Sensor1 = int16(md.Sensor1.UnmarshalSigned(f.Data))
	return nil
}

// TempSensorsReplyReader provides read access to a TempSensorsReply message.
type TempSensorsReplyReader interface {
	// Sensor6 returns the physical value of the Sensor6 signal.
	Sensor6() float64
	// Sensor5 returns the physical value of the Sensor5 signal.
	Sensor5() float64
	// Sensor4 returns the physical value of the Sensor4 signal.
	Sensor4() float64
	// Sensor3 returns the physical value of the Sensor3 signal.
	Sensor3() float64
	// Sensor2 returns the physical value of the Sensor2 signal.
	Sensor2() float64
	// Sensor1 returns the physical value of the Sensor1 signal.
	Sensor1() float64
}

// TempSensorsReplyWriter provides write access to a TempSensorsReply message.
type TempSensorsReplyWriter interface {
	// CopyFrom copies all values from TempSensorsReplyReader.
	CopyFrom(TempSensorsReplyReader) *TempSensorsReply
	// SetSensor6 sets the physical value of the Sensor6 signal.
	SetSensor6(float64) *TempSensorsReply
	// SetSensor5 sets the physical value of the Sensor5 signal.
	SetSensor5(float64) *TempSensorsReply
	// SetSensor4 sets the physical value of the Sensor4 signal.
	SetSensor4(float64) *TempSensorsReply
	// SetSensor3 sets the physical value of the Sensor3 signal.
	SetSensor3(float64) *TempSensorsReply
	// SetSensor2 sets the physical value of the Sensor2 signal.
	SetSensor2(float64) *TempSensorsReply
	// SetSensor1 sets the physical value of the Sensor1 signal.
	SetSensor1(float64) *TempSensorsReply
}

type TempSensorsReply struct {
	xxx_Sensor6 int16
	xxx_Sensor5 int16
	xxx_Sensor4 int16
	xxx_Sensor3 int16
	xxx_Sensor2 int16
	xxx_Sensor1 int16
}

func NewTempSensorsReply() *TempSensorsReply {
	m := &TempSensorsReply{}
	m.Reset()
	return m
}

func (m *TempSensorsReply) Reset() {
	m.xxx_Sensor6 = 0
	m.xxx_Sensor5 = 0
	m.xxx_Sensor4 = 0
	m.xxx_Sensor3 = 0
	m.xxx_Sensor2 = 0
	m.xxx_Sensor1 = 0
}

func (m *TempSensorsReply) CopyFrom(o TempSensorsReplyReader) *TempSensorsReply {
	m.SetSensor6(o.Sensor6())
	m.SetSensor5(o.Sensor5())
	m.SetSensor4(o.Sensor4())
	m.SetSensor3(o.Sensor3())
	m.SetSensor2(o.Sensor2())
	m.SetSensor1(o.Sensor1())
	return m
}

// Descriptor returns the TempSensorsReply descriptor.
func (m *TempSensorsReply) Descriptor() *descriptor.Message {
	return Messages().TempSensorsReply.Message
}

// String returns a compact string representation of the message.
func (m *TempSensorsReply) String() string {
	return cantext.MessageString(m)
}

func (m *TempSensorsReply) Sensor6() float64 {
	return Messages().TempSensorsReply.Sensor6.ToPhysical(float64(m.xxx_Sensor6))
}

func (m *TempSensorsReply) SetSensor6(v float64) *TempSensorsReply {
	m.xxx_Sensor6 = int16(Messages().TempSensorsReply.Sensor6.FromPhysical(v))
	return m
}

func (m *TempSensorsReply) Sensor5() float64 {
	return Messages().TempSensorsReply.Sensor5.ToPhysical(float64(m.xxx_Sensor5))
}

func (m *TempSensorsReply) SetSensor5(v float64) *TempSensorsReply {
	m.xxx_Sensor5 = int16(Messages().TempSensorsReply.Sensor5.FromPhysical(v))
	return m
}

func (m *TempSensorsReply) Sensor4() float64 {
	return Messages().TempSensorsReply.Sensor4.ToPhysical(float64(m.xxx_Sensor4))
}

func (m *TempSensorsReply) SetSensor4(v float64) *TempSensorsReply {
	m.xxx_Sensor4 = int16(Messages().TempSensorsReply.Sensor4.FromPhysical(v))
	return m
}

func (m *TempSensorsReply) Sensor3() float64 {
	return Messages().TempSensorsReply.Sensor3.ToPhysical(float64(m.xxx_Sensor3))
}

func (m *TempSensorsReply) SetSensor3(v float64) *TempSensorsReply {
	m.xxx_Sensor3 = int16(Messages().TempSensorsReply.Sensor3.FromPhysical(v))
	return m
}

func (m *TempSensorsReply) Sensor2() float64 {
	return Messages().TempSensorsReply.Sensor2.ToPhysical(float64(m.xxx_Sensor2))
}

func (m *TempSensorsReply) SetSensor2(v float64) *TempSensorsReply {
	m.xxx_Sensor2 = int16(Messages().TempSensorsReply.Sensor2.FromPhysical(v))
	return m
}

func (m *TempSensorsReply) Sensor1() float64 {
	return Messages().TempSensorsReply.Sensor1.ToPhysical(float64(m.xxx_Sensor1))
}

func (m *TempSensorsReply) SetSensor1(v float64) *TempSensorsReply {
	m.xxx_Sensor1 = int16(Messages().TempSensorsReply.Sensor1.FromPhysical(v))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *TempSensorsReply) Frame() can.Frame {
	md := Messages().TempSensorsReply
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.Sensor6.MarshalSigned(&f.Data, int64(m.xxx_Sensor6))
	md.Sensor5.MarshalSigned(&f.Data, int64(m.xxx_Sensor5))
	md.Sensor4.MarshalSigned(&f.Data, int64(m.xxx_Sensor4))
	md.Sensor3.MarshalSigned(&f.Data, int64(m.xxx_Sensor3))
	md.Sensor2.MarshalSigned(&f.Data, int64(m.xxx_Sensor2))
	md.Sensor1.MarshalSigned(&f.Data, int64(m.xxx_Sensor1))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *TempSensorsReply) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *TempSensorsReply) UnmarshalFrame(f can.Frame) error {
	md := Messages().TempSensorsReply
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal TempSensorsReply: expects ID 941 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal TempSensorsReply: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal TempSensorsReply: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal TempSensorsReply: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_Sensor6 = int16(md.Sensor6.UnmarshalSigned(f.Data))
	m.xxx_Sensor5 = int16(md.Sensor5.UnmarshalSigned(f.Data))
	m.xxx_Sensor4 = int16(md.Sensor4.UnmarshalSigned(f.Data))
	m.xxx_Sensor3 = int16(md.Sensor3.UnmarshalSigned(f.Data))
	m.xxx_Sensor2 = int16(md.Sensor2.UnmarshalSigned(f.Data))
	m.xxx_Sensor1 = int16(md.Sensor1.UnmarshalSigned(f.Data))
	return nil
}

// Nodes returns the demo node descriptors.
func Nodes() *NodesDescriptor {
	return nd
}

// NodesDescriptor contains all demo node descriptors.
type NodesDescriptor struct {
	BAR *descriptor.Node
	FOO *descriptor.Node
}

// Messages returns the demo message descriptors.
func Messages() *MessagesDescriptor {
	return md
}

// MessagesDescriptor contains all demo message descriptors.
type MessagesDescriptor struct {
	TempSensors      *TempSensorsDescriptor
	TempSensorsReply *TempSensorsReplyDescriptor
}

// UnmarshalFrame unmarshals the provided demo CAN frame.
func (md *MessagesDescriptor) UnmarshalFrame(f can.Frame) (generated.Message, error) {
	switch f.ID {
	case md.TempSensors.ID:
		var msg TempSensors
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal demo frame: %w", err)
		}
		return &msg, nil
	case md.TempSensorsReply.ID:
		var msg TempSensorsReply
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal demo frame: %w", err)
		}
		return &msg, nil
	default:
		return nil, fmt.Errorf("unmarshal demo frame: ID not in database: %d", f.ID)
	}
}

type TempSensorsDescriptor struct {
	*descriptor.Message
	Sensor6 *descriptor.Signal
	Sensor5 *descriptor.Signal
	Sensor4 *descriptor.Signal
	Sensor3 *descriptor.Signal
	Sensor2 *descriptor.Signal
	Sensor1 *descriptor.Signal
}

type TempSensorsReplyDescriptor struct {
	*descriptor.Message
	Sensor6 *descriptor.Signal
	Sensor5 *descriptor.Signal
	Sensor4 *descriptor.Signal
	Sensor3 *descriptor.Signal
	Sensor2 *descriptor.Signal
	Sensor1 *descriptor.Signal
}

// Database returns the demo database descriptor.
func (md *MessagesDescriptor) Database() *descriptor.Database {
	return d
}

var nd = &NodesDescriptor{
	BAR: d.Nodes[0],
	FOO: d.Nodes[1],
}

var md = &MessagesDescriptor{
	TempSensors: &TempSensorsDescriptor{
		Message: d.Messages[0],
		Sensor6: d.Messages[0].Signals[0],
		Sensor5: d.Messages[0].Signals[1],
		Sensor4: d.Messages[0].Signals[2],
		Sensor3: d.Messages[0].Signals[3],
		Sensor2: d.Messages[0].Signals[4],
		Sensor1: d.Messages[0].Signals[5],
	},
	TempSensorsReply: &TempSensorsReplyDescriptor{
		Message: d.Messages[1],
		Sensor6: d.Messages[1].Signals[0],
		Sensor5: d.Messages[1].Signals[1],
		Sensor4: d.Messages[1].Signals[2],
		Sensor3: d.Messages[1].Signals[3],
		Sensor2: d.Messages[1].Signals[4],
		Sensor1: d.Messages[1].Signals[5],
	},
}

var d = (*descriptor.Database)(&descriptor.Database{
	SourceFile: (string)("temp/democan/demo.dbc"),
	Version:    (string)(""),
	Messages: ([]*descriptor.Message)([]*descriptor.Message{
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("TempSensors"),
			ID:          (uint32)(940),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Sensor6"),
					Start:             (uint8)(4),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.2),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("degC"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BAR"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Sensor5"),
					Start:             (uint8)(14),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.2),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("degC"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BAR"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Sensor4"),
					Start:             (uint8)(24),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.2),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("degC"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BAR"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Sensor3"),
					Start:             (uint8)(34),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.2),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("degC"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BAR"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Sensor2"),
					Start:             (uint8)(44),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.2),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("degC"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BAR"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Sensor1"),
					Start:             (uint8)(54),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.2),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("degC"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BAR"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("FOO"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("TempSensorsReply"),
			ID:          (uint32)(941),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Sensor6"),
					Start:             (uint8)(4),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.2),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("degC"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FOO"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Sensor5"),
					Start:             (uint8)(14),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.2),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("degC"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FOO"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Sensor4"),
					Start:             (uint8)(24),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.2),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("degC"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FOO"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Sensor3"),
					Start:             (uint8)(34),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.2),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("degC"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FOO"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Sensor2"),
					Start:             (uint8)(44),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.2),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("degC"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FOO"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Sensor1"),
					Start:             (uint8)(54),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.2),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("degC"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FOO"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("BAR"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
	}),
	Nodes: ([]*descriptor.Node)([]*descriptor.Node{
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("BAR"),
			Description: (string)(""),
		}),
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("FOO"),
			Description: (string)(""),
		}),
	}),
})
