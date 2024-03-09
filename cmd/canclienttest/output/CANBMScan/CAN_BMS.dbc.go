// Package CANBMScan provides primitives for encoding and decoding CAN_BMS CAN messages.
//
// Source: CAN_BMS.dbc
package CANBMScan

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
// Contactor_StatesReader provides read access to a Contactor_States message.
type Contactor_StatesReader interface {
	// Pack_Positive returns the value of the Pack_Positive signal.
	Pack_Positive() uint8
	// Pack_Precharge returns the value of the Pack_Precharge signal.
	Pack_Precharge() uint8
	// Pack_Negative returns the value of the Pack_Negative signal.
	Pack_Negative() uint8
}

// Contactor_StatesWriter provides write access to a Contactor_States message.
type Contactor_StatesWriter interface {
	// CopyFrom copies all values from Contactor_StatesReader.
	CopyFrom(Contactor_StatesReader) *Contactor_States
	// SetPack_Positive sets the value of the Pack_Positive signal.
	SetPack_Positive(uint8) *Contactor_States
	// SetPack_Precharge sets the value of the Pack_Precharge signal.
	SetPack_Precharge(uint8) *Contactor_States
	// SetPack_Negative sets the value of the Pack_Negative signal.
	SetPack_Negative(uint8) *Contactor_States
}

type Contactor_States struct {
	xxx_Pack_Positive  uint8
	xxx_Pack_Precharge uint8
	xxx_Pack_Negative  uint8
}

func NewContactor_States() *Contactor_States {
	m := &Contactor_States{}
	m.Reset()
	return m
}

func (m *Contactor_States) Reset() {
	m.xxx_Pack_Positive = 0
	m.xxx_Pack_Precharge = 0
	m.xxx_Pack_Negative = 0
}

func (m *Contactor_States) CopyFrom(o Contactor_StatesReader) *Contactor_States {
	m.xxx_Pack_Positive = o.Pack_Positive()
	m.xxx_Pack_Precharge = o.Pack_Precharge()
	m.xxx_Pack_Negative = o.Pack_Negative()
	return m
}

// Descriptor returns the Contactor_States descriptor.
func (m *Contactor_States) Descriptor() *descriptor.Message {
	return Messages().Contactor_States.Message
}

// String returns a compact string representation of the message.
func (m *Contactor_States) String() string {
	return cantext.MessageString(m)
}

func (m *Contactor_States) Pack_Positive() uint8 {
	return m.xxx_Pack_Positive
}

func (m *Contactor_States) SetPack_Positive(v uint8) *Contactor_States {
	m.xxx_Pack_Positive = uint8(Messages().Contactor_States.Pack_Positive.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *Contactor_States) Pack_Precharge() uint8 {
	return m.xxx_Pack_Precharge
}

func (m *Contactor_States) SetPack_Precharge(v uint8) *Contactor_States {
	m.xxx_Pack_Precharge = uint8(Messages().Contactor_States.Pack_Precharge.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *Contactor_States) Pack_Negative() uint8 {
	return m.xxx_Pack_Negative
}

func (m *Contactor_States) SetPack_Negative(v uint8) *Contactor_States {
	m.xxx_Pack_Negative = uint8(Messages().Contactor_States.Pack_Negative.SaturatedCastUnsigned(uint64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *Contactor_States) Frame() can.Frame {
	md := Messages().Contactor_States
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.Pack_Positive.MarshalUnsigned(&f.Data, uint64(m.xxx_Pack_Positive))
	md.Pack_Precharge.MarshalUnsigned(&f.Data, uint64(m.xxx_Pack_Precharge))
	md.Pack_Negative.MarshalUnsigned(&f.Data, uint64(m.xxx_Pack_Negative))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *Contactor_States) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *Contactor_States) UnmarshalFrame(f can.Frame) error {
	md := Messages().Contactor_States
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal Contactor_States: expects ID 1570 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal Contactor_States: expects length 3 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal Contactor_States: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal Contactor_States: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_Pack_Positive = uint8(md.Pack_Positive.UnmarshalUnsigned(f.Data))
	m.xxx_Pack_Precharge = uint8(md.Pack_Precharge.UnmarshalUnsigned(f.Data))
	m.xxx_Pack_Negative = uint8(md.Pack_Negative.UnmarshalUnsigned(f.Data))
	return nil
}

// Pack_Current_LimitsReader provides read access to a Pack_Current_Limits message.
type Pack_Current_LimitsReader interface {
	// Pack_CCL returns the value of the Pack_CCL signal.
	Pack_CCL() uint16
	// Pack_DCL returns the value of the Pack_DCL signal.
	Pack_DCL() uint16
}

// Pack_Current_LimitsWriter provides write access to a Pack_Current_Limits message.
type Pack_Current_LimitsWriter interface {
	// CopyFrom copies all values from Pack_Current_LimitsReader.
	CopyFrom(Pack_Current_LimitsReader) *Pack_Current_Limits
	// SetPack_CCL sets the value of the Pack_CCL signal.
	SetPack_CCL(uint16) *Pack_Current_Limits
	// SetPack_DCL sets the value of the Pack_DCL signal.
	SetPack_DCL(uint16) *Pack_Current_Limits
}

type Pack_Current_Limits struct {
	xxx_Pack_CCL uint16
	xxx_Pack_DCL uint16
}

func NewPack_Current_Limits() *Pack_Current_Limits {
	m := &Pack_Current_Limits{}
	m.Reset()
	return m
}

func (m *Pack_Current_Limits) Reset() {
	m.xxx_Pack_CCL = 0
	m.xxx_Pack_DCL = 0
}

func (m *Pack_Current_Limits) CopyFrom(o Pack_Current_LimitsReader) *Pack_Current_Limits {
	m.xxx_Pack_CCL = o.Pack_CCL()
	m.xxx_Pack_DCL = o.Pack_DCL()
	return m
}

// Descriptor returns the Pack_Current_Limits descriptor.
func (m *Pack_Current_Limits) Descriptor() *descriptor.Message {
	return Messages().Pack_Current_Limits.Message
}

// String returns a compact string representation of the message.
func (m *Pack_Current_Limits) String() string {
	return cantext.MessageString(m)
}

func (m *Pack_Current_Limits) Pack_CCL() uint16 {
	return m.xxx_Pack_CCL
}

func (m *Pack_Current_Limits) SetPack_CCL(v uint16) *Pack_Current_Limits {
	m.xxx_Pack_CCL = uint16(Messages().Pack_Current_Limits.Pack_CCL.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *Pack_Current_Limits) Pack_DCL() uint16 {
	return m.xxx_Pack_DCL
}

func (m *Pack_Current_Limits) SetPack_DCL(v uint16) *Pack_Current_Limits {
	m.xxx_Pack_DCL = uint16(Messages().Pack_Current_Limits.Pack_DCL.SaturatedCastUnsigned(uint64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *Pack_Current_Limits) Frame() can.Frame {
	md := Messages().Pack_Current_Limits
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.Pack_CCL.MarshalUnsigned(&f.Data, uint64(m.xxx_Pack_CCL))
	md.Pack_DCL.MarshalUnsigned(&f.Data, uint64(m.xxx_Pack_DCL))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *Pack_Current_Limits) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *Pack_Current_Limits) UnmarshalFrame(f can.Frame) error {
	md := Messages().Pack_Current_Limits
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal Pack_Current_Limits: expects ID 1571 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal Pack_Current_Limits: expects length 4 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal Pack_Current_Limits: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal Pack_Current_Limits: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_Pack_CCL = uint16(md.Pack_CCL.UnmarshalUnsigned(f.Data))
	m.xxx_Pack_DCL = uint16(md.Pack_DCL.UnmarshalUnsigned(f.Data))
	return nil
}

// Pack_StateReader provides read access to a Pack_State message.
type Pack_StateReader interface {
	// Pack_Current returns the physical value of the Pack_Current signal.
	Pack_Current() float64
	// Pack_Inst_Voltage returns the physical value of the Pack_Inst_Voltage signal.
	Pack_Inst_Voltage() float64
	// Avg_Cell_Voltage returns the physical value of the Avg_Cell_Voltage signal.
	Avg_Cell_Voltage() float64
	// Populated_Cells returns the value of the Populated_Cells signal.
	Populated_Cells() uint8
}

// Pack_StateWriter provides write access to a Pack_State message.
type Pack_StateWriter interface {
	// CopyFrom copies all values from Pack_StateReader.
	CopyFrom(Pack_StateReader) *Pack_State
	// SetPack_Current sets the physical value of the Pack_Current signal.
	SetPack_Current(float64) *Pack_State
	// SetPack_Inst_Voltage sets the physical value of the Pack_Inst_Voltage signal.
	SetPack_Inst_Voltage(float64) *Pack_State
	// SetAvg_Cell_Voltage sets the physical value of the Avg_Cell_Voltage signal.
	SetAvg_Cell_Voltage(float64) *Pack_State
	// SetPopulated_Cells sets the value of the Populated_Cells signal.
	SetPopulated_Cells(uint8) *Pack_State
}

type Pack_State struct {
	xxx_Pack_Current      uint16
	xxx_Pack_Inst_Voltage uint16
	xxx_Avg_Cell_Voltage  uint16
	xxx_Populated_Cells   uint8
}

func NewPack_State() *Pack_State {
	m := &Pack_State{}
	m.Reset()
	return m
}

func (m *Pack_State) Reset() {
	m.xxx_Pack_Current = 0
	m.xxx_Pack_Inst_Voltage = 0
	m.xxx_Avg_Cell_Voltage = 0
	m.xxx_Populated_Cells = 0
}

func (m *Pack_State) CopyFrom(o Pack_StateReader) *Pack_State {
	m.SetPack_Current(o.Pack_Current())
	m.SetPack_Inst_Voltage(o.Pack_Inst_Voltage())
	m.SetAvg_Cell_Voltage(o.Avg_Cell_Voltage())
	m.xxx_Populated_Cells = o.Populated_Cells()
	return m
}

// Descriptor returns the Pack_State descriptor.
func (m *Pack_State) Descriptor() *descriptor.Message {
	return Messages().Pack_State.Message
}

// String returns a compact string representation of the message.
func (m *Pack_State) String() string {
	return cantext.MessageString(m)
}

func (m *Pack_State) Pack_Current() float64 {
	return Messages().Pack_State.Pack_Current.ToPhysical(float64(m.xxx_Pack_Current))
}

func (m *Pack_State) SetPack_Current(v float64) *Pack_State {
	m.xxx_Pack_Current = uint16(Messages().Pack_State.Pack_Current.FromPhysical(v))
	return m
}

func (m *Pack_State) Pack_Inst_Voltage() float64 {
	return Messages().Pack_State.Pack_Inst_Voltage.ToPhysical(float64(m.xxx_Pack_Inst_Voltage))
}

func (m *Pack_State) SetPack_Inst_Voltage(v float64) *Pack_State {
	m.xxx_Pack_Inst_Voltage = uint16(Messages().Pack_State.Pack_Inst_Voltage.FromPhysical(v))
	return m
}

func (m *Pack_State) Avg_Cell_Voltage() float64 {
	return Messages().Pack_State.Avg_Cell_Voltage.ToPhysical(float64(m.xxx_Avg_Cell_Voltage))
}

func (m *Pack_State) SetAvg_Cell_Voltage(v float64) *Pack_State {
	m.xxx_Avg_Cell_Voltage = uint16(Messages().Pack_State.Avg_Cell_Voltage.FromPhysical(v))
	return m
}

func (m *Pack_State) Populated_Cells() uint8 {
	return m.xxx_Populated_Cells
}

func (m *Pack_State) SetPopulated_Cells(v uint8) *Pack_State {
	m.xxx_Populated_Cells = uint8(Messages().Pack_State.Populated_Cells.SaturatedCastUnsigned(uint64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *Pack_State) Frame() can.Frame {
	md := Messages().Pack_State
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.Pack_Current.MarshalUnsigned(&f.Data, uint64(m.xxx_Pack_Current))
	md.Pack_Inst_Voltage.MarshalUnsigned(&f.Data, uint64(m.xxx_Pack_Inst_Voltage))
	md.Avg_Cell_Voltage.MarshalUnsigned(&f.Data, uint64(m.xxx_Avg_Cell_Voltage))
	md.Populated_Cells.MarshalUnsigned(&f.Data, uint64(m.xxx_Populated_Cells))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *Pack_State) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *Pack_State) UnmarshalFrame(f can.Frame) error {
	md := Messages().Pack_State
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal Pack_State: expects ID 1572 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal Pack_State: expects length 7 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal Pack_State: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal Pack_State: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_Pack_Current = uint16(md.Pack_Current.UnmarshalUnsigned(f.Data))
	m.xxx_Pack_Inst_Voltage = uint16(md.Pack_Inst_Voltage.UnmarshalUnsigned(f.Data))
	m.xxx_Avg_Cell_Voltage = uint16(md.Avg_Cell_Voltage.UnmarshalUnsigned(f.Data))
	m.xxx_Populated_Cells = uint8(md.Populated_Cells.UnmarshalUnsigned(f.Data))
	return nil
}

// Pack_SOCReader provides read access to a Pack_SOC message.
type Pack_SOCReader interface {
	// Pack_SOC returns the physical value of the Pack_SOC signal.
	Pack_SOC() float64
	// Maximum_Pack_Voltage returns the physical value of the Maximum_Pack_Voltage signal.
	Maximum_Pack_Voltage() float64
}

// Pack_SOCWriter provides write access to a Pack_SOC message.
type Pack_SOCWriter interface {
	// CopyFrom copies all values from Pack_SOCReader.
	CopyFrom(Pack_SOCReader) *Pack_SOC
	// SetPack_SOC sets the physical value of the Pack_SOC signal.
	SetPack_SOC(float64) *Pack_SOC
	// SetMaximum_Pack_Voltage sets the physical value of the Maximum_Pack_Voltage signal.
	SetMaximum_Pack_Voltage(float64) *Pack_SOC
}

type Pack_SOC struct {
	xxx_Pack_SOC             uint8
	xxx_Maximum_Pack_Voltage uint16
}

func NewPack_SOC() *Pack_SOC {
	m := &Pack_SOC{}
	m.Reset()
	return m
}

func (m *Pack_SOC) Reset() {
	m.xxx_Pack_SOC = 0
	m.xxx_Maximum_Pack_Voltage = 0
}

func (m *Pack_SOC) CopyFrom(o Pack_SOCReader) *Pack_SOC {
	m.SetPack_SOC(o.Pack_SOC())
	m.SetMaximum_Pack_Voltage(o.Maximum_Pack_Voltage())
	return m
}

// Descriptor returns the Pack_SOC descriptor.
func (m *Pack_SOC) Descriptor() *descriptor.Message {
	return Messages().Pack_SOC.Message
}

// String returns a compact string representation of the message.
func (m *Pack_SOC) String() string {
	return cantext.MessageString(m)
}

func (m *Pack_SOC) Pack_SOC() float64 {
	return Messages().Pack_SOC.Pack_SOC.ToPhysical(float64(m.xxx_Pack_SOC))
}

func (m *Pack_SOC) SetPack_SOC(v float64) *Pack_SOC {
	m.xxx_Pack_SOC = uint8(Messages().Pack_SOC.Pack_SOC.FromPhysical(v))
	return m
}

func (m *Pack_SOC) Maximum_Pack_Voltage() float64 {
	return Messages().Pack_SOC.Maximum_Pack_Voltage.ToPhysical(float64(m.xxx_Maximum_Pack_Voltage))
}

func (m *Pack_SOC) SetMaximum_Pack_Voltage(v float64) *Pack_SOC {
	m.xxx_Maximum_Pack_Voltage = uint16(Messages().Pack_SOC.Maximum_Pack_Voltage.FromPhysical(v))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *Pack_SOC) Frame() can.Frame {
	md := Messages().Pack_SOC
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.Pack_SOC.MarshalUnsigned(&f.Data, uint64(m.xxx_Pack_SOC))
	md.Maximum_Pack_Voltage.MarshalUnsigned(&f.Data, uint64(m.xxx_Maximum_Pack_Voltage))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *Pack_SOC) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *Pack_SOC) UnmarshalFrame(f can.Frame) error {
	md := Messages().Pack_SOC
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal Pack_SOC: expects ID 1573 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal Pack_SOC: expects length 3 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal Pack_SOC: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal Pack_SOC: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_Pack_SOC = uint8(md.Pack_SOC.UnmarshalUnsigned(f.Data))
	m.xxx_Maximum_Pack_Voltage = uint16(md.Maximum_Pack_Voltage.UnmarshalUnsigned(f.Data))
	return nil
}

// Contactor_FeedbackReader provides read access to a Contactor_Feedback message.
type Contactor_FeedbackReader interface {
	// Pack_Precharge_Feedback returns the value of the Pack_Precharge_Feedback signal.
	Pack_Precharge_Feedback() bool
	// Pack_Negative_Feedback returns the value of the Pack_Negative_Feedback signal.
	Pack_Negative_Feedback() bool
	// Pack_Positive_Feedback returns the value of the Pack_Positive_Feedback signal.
	Pack_Positive_Feedback() bool
}

// Contactor_FeedbackWriter provides write access to a Contactor_Feedback message.
type Contactor_FeedbackWriter interface {
	// CopyFrom copies all values from Contactor_FeedbackReader.
	CopyFrom(Contactor_FeedbackReader) *Contactor_Feedback
	// SetPack_Precharge_Feedback sets the value of the Pack_Precharge_Feedback signal.
	SetPack_Precharge_Feedback(bool) *Contactor_Feedback
	// SetPack_Negative_Feedback sets the value of the Pack_Negative_Feedback signal.
	SetPack_Negative_Feedback(bool) *Contactor_Feedback
	// SetPack_Positive_Feedback sets the value of the Pack_Positive_Feedback signal.
	SetPack_Positive_Feedback(bool) *Contactor_Feedback
}

type Contactor_Feedback struct {
	xxx_Pack_Precharge_Feedback bool
	xxx_Pack_Negative_Feedback  bool
	xxx_Pack_Positive_Feedback  bool
}

func NewContactor_Feedback() *Contactor_Feedback {
	m := &Contactor_Feedback{}
	m.Reset()
	return m
}

func (m *Contactor_Feedback) Reset() {
	m.xxx_Pack_Precharge_Feedback = false
	m.xxx_Pack_Negative_Feedback = false
	m.xxx_Pack_Positive_Feedback = false
}

func (m *Contactor_Feedback) CopyFrom(o Contactor_FeedbackReader) *Contactor_Feedback {
	m.xxx_Pack_Precharge_Feedback = o.Pack_Precharge_Feedback()
	m.xxx_Pack_Negative_Feedback = o.Pack_Negative_Feedback()
	m.xxx_Pack_Positive_Feedback = o.Pack_Positive_Feedback()
	return m
}

// Descriptor returns the Contactor_Feedback descriptor.
func (m *Contactor_Feedback) Descriptor() *descriptor.Message {
	return Messages().Contactor_Feedback.Message
}

// String returns a compact string representation of the message.
func (m *Contactor_Feedback) String() string {
	return cantext.MessageString(m)
}

func (m *Contactor_Feedback) Pack_Precharge_Feedback() bool {
	return m.xxx_Pack_Precharge_Feedback
}

func (m *Contactor_Feedback) SetPack_Precharge_Feedback(v bool) *Contactor_Feedback {
	m.xxx_Pack_Precharge_Feedback = v
	return m
}

func (m *Contactor_Feedback) Pack_Negative_Feedback() bool {
	return m.xxx_Pack_Negative_Feedback
}

func (m *Contactor_Feedback) SetPack_Negative_Feedback(v bool) *Contactor_Feedback {
	m.xxx_Pack_Negative_Feedback = v
	return m
}

func (m *Contactor_Feedback) Pack_Positive_Feedback() bool {
	return m.xxx_Pack_Positive_Feedback
}

func (m *Contactor_Feedback) SetPack_Positive_Feedback(v bool) *Contactor_Feedback {
	m.xxx_Pack_Positive_Feedback = v
	return m
}

// Frame returns a CAN frame representing the message.
func (m *Contactor_Feedback) Frame() can.Frame {
	md := Messages().Contactor_Feedback
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.Pack_Precharge_Feedback.MarshalBool(&f.Data, bool(m.xxx_Pack_Precharge_Feedback))
	md.Pack_Negative_Feedback.MarshalBool(&f.Data, bool(m.xxx_Pack_Negative_Feedback))
	md.Pack_Positive_Feedback.MarshalBool(&f.Data, bool(m.xxx_Pack_Positive_Feedback))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *Contactor_Feedback) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *Contactor_Feedback) UnmarshalFrame(f can.Frame) error {
	md := Messages().Contactor_Feedback
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal Contactor_Feedback: expects ID 1574 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal Contactor_Feedback: expects length 1 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal Contactor_Feedback: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal Contactor_Feedback: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_Pack_Precharge_Feedback = bool(md.Pack_Precharge_Feedback.UnmarshalBool(f.Data))
	m.xxx_Pack_Negative_Feedback = bool(md.Pack_Negative_Feedback.UnmarshalBool(f.Data))
	m.xxx_Pack_Positive_Feedback = bool(md.Pack_Positive_Feedback.UnmarshalBool(f.Data))
	return nil
}

// Nodes returns the CAN_BMS node descriptors.
func Nodes() *NodesDescriptor {
	return nd
}

// NodesDescriptor contains all CAN_BMS node descriptors.
type NodesDescriptor struct {
	BMS *descriptor.Node
	FC  *descriptor.Node
}

// Messages returns the CAN_BMS message descriptors.
func Messages() *MessagesDescriptor {
	return md
}

// MessagesDescriptor contains all CAN_BMS message descriptors.
type MessagesDescriptor struct {
	Contactor_States    *Contactor_StatesDescriptor
	Pack_Current_Limits *Pack_Current_LimitsDescriptor
	Pack_State          *Pack_StateDescriptor
	Pack_SOC            *Pack_SOCDescriptor
	Contactor_Feedback  *Contactor_FeedbackDescriptor
}

// UnmarshalFrame unmarshals the provided CAN_BMS CAN frame.
func (md *MessagesDescriptor) UnmarshalFrame(f can.Frame) (generated.Message, error) {
	switch f.ID {
	case md.Contactor_States.ID:
		var msg Contactor_States
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal CAN_BMS frame: %w", err)
		}
		return &msg, nil
	case md.Pack_Current_Limits.ID:
		var msg Pack_Current_Limits
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal CAN_BMS frame: %w", err)
		}
		return &msg, nil
	case md.Pack_State.ID:
		var msg Pack_State
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal CAN_BMS frame: %w", err)
		}
		return &msg, nil
	case md.Pack_SOC.ID:
		var msg Pack_SOC
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal CAN_BMS frame: %w", err)
		}
		return &msg, nil
	case md.Contactor_Feedback.ID:
		var msg Contactor_Feedback
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal CAN_BMS frame: %w", err)
		}
		return &msg, nil
	default:
		return nil, fmt.Errorf("unmarshal CAN_BMS frame: ID not in database: %d", f.ID)
	}
}

type Contactor_StatesDescriptor struct {
	*descriptor.Message
	Pack_Positive  *descriptor.Signal
	Pack_Precharge *descriptor.Signal
	Pack_Negative  *descriptor.Signal
}

type Pack_Current_LimitsDescriptor struct {
	*descriptor.Message
	Pack_CCL *descriptor.Signal
	Pack_DCL *descriptor.Signal
}

type Pack_StateDescriptor struct {
	*descriptor.Message
	Pack_Current      *descriptor.Signal
	Pack_Inst_Voltage *descriptor.Signal
	Avg_Cell_Voltage  *descriptor.Signal
	Populated_Cells   *descriptor.Signal
}

type Pack_SOCDescriptor struct {
	*descriptor.Message
	Pack_SOC             *descriptor.Signal
	Maximum_Pack_Voltage *descriptor.Signal
}

type Contactor_FeedbackDescriptor struct {
	*descriptor.Message
	Pack_Precharge_Feedback *descriptor.Signal
	Pack_Negative_Feedback  *descriptor.Signal
	Pack_Positive_Feedback  *descriptor.Signal
}

// Database returns the CAN_BMS database descriptor.
func (md *MessagesDescriptor) Database() *descriptor.Database {
	return d
}

var nd = &NodesDescriptor{
	BMS: d.Nodes[0],
	FC:  d.Nodes[1],
}

var md = &MessagesDescriptor{
	Contactor_States: &Contactor_StatesDescriptor{
		Message:        d.Messages[0],
		Pack_Positive:  d.Messages[0].Signals[0],
		Pack_Precharge: d.Messages[0].Signals[1],
		Pack_Negative:  d.Messages[0].Signals[2],
	},
	Pack_Current_Limits: &Pack_Current_LimitsDescriptor{
		Message:  d.Messages[1],
		Pack_CCL: d.Messages[1].Signals[0],
		Pack_DCL: d.Messages[1].Signals[1],
	},
	Pack_State: &Pack_StateDescriptor{
		Message:           d.Messages[2],
		Pack_Current:      d.Messages[2].Signals[0],
		Pack_Inst_Voltage: d.Messages[2].Signals[1],
		Avg_Cell_Voltage:  d.Messages[2].Signals[2],
		Populated_Cells:   d.Messages[2].Signals[3],
	},
	Pack_SOC: &Pack_SOCDescriptor{
		Message:              d.Messages[3],
		Pack_SOC:             d.Messages[3].Signals[0],
		Maximum_Pack_Voltage: d.Messages[3].Signals[1],
	},
	Contactor_Feedback: &Contactor_FeedbackDescriptor{
		Message:                 d.Messages[4],
		Pack_Precharge_Feedback: d.Messages[4].Signals[0],
		Pack_Negative_Feedback:  d.Messages[4].Signals[1],
		Pack_Positive_Feedback:  d.Messages[4].Signals[2],
	},
}

var d = (*descriptor.Database)(&descriptor.Database{
	SourceFile: (string)("CAN_BMS.dbc"),
	Version:    (string)(""),
	Messages: ([]*descriptor.Message)([]*descriptor.Message{
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("Contactor_States"),
			ID:          (uint32)(1570),
			IsExtended:  (bool)(false),
			Length:      (uint8)(3),
			SendType:    (descriptor.SendType)(0),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Pack_Positive"),
					Start:             (uint8)(0),
					Length:            (uint8)(8),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Pack_Precharge"),
					Start:             (uint8)(8),
					Length:            (uint8)(8),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Pack_Negative"),
					Start:             (uint8)(16),
					Length:            (uint8)(8),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("FC"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("Pack_Current_Limits"),
			ID:          (uint32)(1571),
			IsExtended:  (bool)(false),
			Length:      (uint8)(4),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("This ID Transmits at 8 ms."),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Pack_CCL"),
					Start:             (uint8)(0),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("Amps"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Pack_DCL"),
					Start:             (uint8)(16),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("Amps"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("BMS"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("Pack_State"),
			ID:          (uint32)(1572),
			IsExtended:  (bool)(false),
			Length:      (uint8)(7),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("This ID Transmits at 8 ms."),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Pack_Current"),
					Start:             (uint8)(0),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("Amps"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Pack_Inst_Voltage"),
					Start:             (uint8)(16),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("Volts"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Avg_Cell_Voltage"),
					Start:             (uint8)(32),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.0001),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("Volts"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Populated_Cells"),
					Start:             (uint8)(48),
					Length:            (uint8)(8),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("Num"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("BMS"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("Pack_SOC"),
			ID:          (uint32)(1573),
			IsExtended:  (bool)(false),
			Length:      (uint8)(3),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("This ID Transmits at 8 ms."),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Pack_SOC"),
					Start:             (uint8)(0),
					Length:            (uint8)(8),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.5),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("Percent"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Maximum_Pack_Voltage"),
					Start:             (uint8)(8),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("Volts"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("BMS"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("Contactor_Feedback"),
			ID:          (uint32)(1574),
			IsExtended:  (bool)(false),
			Length:      (uint8)(1),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("This ID Transmits at 8 ms."),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Pack_Precharge_Feedback"),
					Start:             (uint8)(0),
					Length:            (uint8)(1),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(1),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Pack_Negative_Feedback"),
					Start:             (uint8)(1),
					Length:            (uint8)(1),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(1),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Pack_Positive_Feedback"),
					Start:             (uint8)(2),
					Length:            (uint8)(1),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(1),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("BMS"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
	}),
	Nodes: ([]*descriptor.Node)([]*descriptor.Node{
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("BMS"),
			Description: (string)(""),
		}),
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("FC"),
			Description: (string)(""),
		}),
	}),
})
