// Package ptcan provides primitives for encoding and decoding pt CAN messages.
//
// Source: temp/ptcan/pt.dbc
package ptcan

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
// AMK0_SetPoints1Reader provides read access to a AMK0_SetPoints1 message.
type AMK0_SetPoints1Reader interface {
	// AMK_bInverterOn returns the value of the AMK_bInverterOn signal.
	AMK_bInverterOn() bool
	// AMK_bDcOn returns the value of the AMK_bDcOn signal.
	AMK_bDcOn() bool
	// AMK_bEnable returns the value of the AMK_bEnable signal.
	AMK_bEnable() bool
	// AMK_bErrorReset returns the value of the AMK_bErrorReset signal.
	AMK_bErrorReset() bool
	// AMK_TargetVelocity returns the value of the AMK_TargetVelocity signal.
	AMK_TargetVelocity() int16
	// AMK_TorqueLimitPositiv returns the value of the AMK_TorqueLimitPositiv signal.
	AMK_TorqueLimitPositiv() int16
	// AMK_TorqueLimitNegativ returns the value of the AMK_TorqueLimitNegativ signal.
	AMK_TorqueLimitNegativ() int16
}

// AMK0_SetPoints1Writer provides write access to a AMK0_SetPoints1 message.
type AMK0_SetPoints1Writer interface {
	// CopyFrom copies all values from AMK0_SetPoints1Reader.
	CopyFrom(AMK0_SetPoints1Reader) *AMK0_SetPoints1
	// SetAMK_bInverterOn sets the value of the AMK_bInverterOn signal.
	SetAMK_bInverterOn(bool) *AMK0_SetPoints1
	// SetAMK_bDcOn sets the value of the AMK_bDcOn signal.
	SetAMK_bDcOn(bool) *AMK0_SetPoints1
	// SetAMK_bEnable sets the value of the AMK_bEnable signal.
	SetAMK_bEnable(bool) *AMK0_SetPoints1
	// SetAMK_bErrorReset sets the value of the AMK_bErrorReset signal.
	SetAMK_bErrorReset(bool) *AMK0_SetPoints1
	// SetAMK_TargetVelocity sets the value of the AMK_TargetVelocity signal.
	SetAMK_TargetVelocity(int16) *AMK0_SetPoints1
	// SetAMK_TorqueLimitPositiv sets the value of the AMK_TorqueLimitPositiv signal.
	SetAMK_TorqueLimitPositiv(int16) *AMK0_SetPoints1
	// SetAMK_TorqueLimitNegativ sets the value of the AMK_TorqueLimitNegativ signal.
	SetAMK_TorqueLimitNegativ(int16) *AMK0_SetPoints1
}

type AMK0_SetPoints1 struct {
	xxx_AMK_bInverterOn        bool
	xxx_AMK_bDcOn              bool
	xxx_AMK_bEnable            bool
	xxx_AMK_bErrorReset        bool
	xxx_AMK_TargetVelocity     int16
	xxx_AMK_TorqueLimitPositiv int16
	xxx_AMK_TorqueLimitNegativ int16
}

func NewAMK0_SetPoints1() *AMK0_SetPoints1 {
	m := &AMK0_SetPoints1{}
	m.Reset()
	return m
}

func (m *AMK0_SetPoints1) Reset() {
	m.xxx_AMK_bInverterOn = false
	m.xxx_AMK_bDcOn = false
	m.xxx_AMK_bEnable = false
	m.xxx_AMK_bErrorReset = false
	m.xxx_AMK_TargetVelocity = 0
	m.xxx_AMK_TorqueLimitPositiv = 0
	m.xxx_AMK_TorqueLimitNegativ = 0
}

func (m *AMK0_SetPoints1) CopyFrom(o AMK0_SetPoints1Reader) *AMK0_SetPoints1 {
	m.xxx_AMK_bInverterOn = o.AMK_bInverterOn()
	m.xxx_AMK_bDcOn = o.AMK_bDcOn()
	m.xxx_AMK_bEnable = o.AMK_bEnable()
	m.xxx_AMK_bErrorReset = o.AMK_bErrorReset()
	m.xxx_AMK_TargetVelocity = o.AMK_TargetVelocity()
	m.xxx_AMK_TorqueLimitPositiv = o.AMK_TorqueLimitPositiv()
	m.xxx_AMK_TorqueLimitNegativ = o.AMK_TorqueLimitNegativ()
	return m
}

// Descriptor returns the AMK0_SetPoints1 descriptor.
func (m *AMK0_SetPoints1) Descriptor() *descriptor.Message {
	return Messages().AMK0_SetPoints1.Message
}

// String returns a compact string representation of the message.
func (m *AMK0_SetPoints1) String() string {
	return cantext.MessageString(m)
}

func (m *AMK0_SetPoints1) AMK_bInverterOn() bool {
	return m.xxx_AMK_bInverterOn
}

func (m *AMK0_SetPoints1) SetAMK_bInverterOn(v bool) *AMK0_SetPoints1 {
	m.xxx_AMK_bInverterOn = v
	return m
}

func (m *AMK0_SetPoints1) AMK_bDcOn() bool {
	return m.xxx_AMK_bDcOn
}

func (m *AMK0_SetPoints1) SetAMK_bDcOn(v bool) *AMK0_SetPoints1 {
	m.xxx_AMK_bDcOn = v
	return m
}

func (m *AMK0_SetPoints1) AMK_bEnable() bool {
	return m.xxx_AMK_bEnable
}

func (m *AMK0_SetPoints1) SetAMK_bEnable(v bool) *AMK0_SetPoints1 {
	m.xxx_AMK_bEnable = v
	return m
}

func (m *AMK0_SetPoints1) AMK_bErrorReset() bool {
	return m.xxx_AMK_bErrorReset
}

func (m *AMK0_SetPoints1) SetAMK_bErrorReset(v bool) *AMK0_SetPoints1 {
	m.xxx_AMK_bErrorReset = v
	return m
}

func (m *AMK0_SetPoints1) AMK_TargetVelocity() int16 {
	return m.xxx_AMK_TargetVelocity
}

func (m *AMK0_SetPoints1) SetAMK_TargetVelocity(v int16) *AMK0_SetPoints1 {
	m.xxx_AMK_TargetVelocity = int16(Messages().AMK0_SetPoints1.AMK_TargetVelocity.SaturatedCastSigned(int64(v)))
	return m
}

func (m *AMK0_SetPoints1) AMK_TorqueLimitPositiv() int16 {
	return m.xxx_AMK_TorqueLimitPositiv
}

func (m *AMK0_SetPoints1) SetAMK_TorqueLimitPositiv(v int16) *AMK0_SetPoints1 {
	m.xxx_AMK_TorqueLimitPositiv = int16(Messages().AMK0_SetPoints1.AMK_TorqueLimitPositiv.SaturatedCastSigned(int64(v)))
	return m
}

func (m *AMK0_SetPoints1) AMK_TorqueLimitNegativ() int16 {
	return m.xxx_AMK_TorqueLimitNegativ
}

func (m *AMK0_SetPoints1) SetAMK_TorqueLimitNegativ(v int16) *AMK0_SetPoints1 {
	m.xxx_AMK_TorqueLimitNegativ = int16(Messages().AMK0_SetPoints1.AMK_TorqueLimitNegativ.SaturatedCastSigned(int64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *AMK0_SetPoints1) Frame() can.Frame {
	md := Messages().AMK0_SetPoints1
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.AMK_bInverterOn.MarshalBool(&f.Data, bool(m.xxx_AMK_bInverterOn))
	md.AMK_bDcOn.MarshalBool(&f.Data, bool(m.xxx_AMK_bDcOn))
	md.AMK_bEnable.MarshalBool(&f.Data, bool(m.xxx_AMK_bEnable))
	md.AMK_bErrorReset.MarshalBool(&f.Data, bool(m.xxx_AMK_bErrorReset))
	md.AMK_TargetVelocity.MarshalSigned(&f.Data, int64(m.xxx_AMK_TargetVelocity))
	md.AMK_TorqueLimitPositiv.MarshalSigned(&f.Data, int64(m.xxx_AMK_TorqueLimitPositiv))
	md.AMK_TorqueLimitNegativ.MarshalSigned(&f.Data, int64(m.xxx_AMK_TorqueLimitNegativ))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *AMK0_SetPoints1) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *AMK0_SetPoints1) UnmarshalFrame(f can.Frame) error {
	md := Messages().AMK0_SetPoints1
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal AMK0_SetPoints1: expects ID 389 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal AMK0_SetPoints1: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal AMK0_SetPoints1: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal AMK0_SetPoints1: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_AMK_bInverterOn = bool(md.AMK_bInverterOn.UnmarshalBool(f.Data))
	m.xxx_AMK_bDcOn = bool(md.AMK_bDcOn.UnmarshalBool(f.Data))
	m.xxx_AMK_bEnable = bool(md.AMK_bEnable.UnmarshalBool(f.Data))
	m.xxx_AMK_bErrorReset = bool(md.AMK_bErrorReset.UnmarshalBool(f.Data))
	m.xxx_AMK_TargetVelocity = int16(md.AMK_TargetVelocity.UnmarshalSigned(f.Data))
	m.xxx_AMK_TorqueLimitPositiv = int16(md.AMK_TorqueLimitPositiv.UnmarshalSigned(f.Data))
	m.xxx_AMK_TorqueLimitNegativ = int16(md.AMK_TorqueLimitNegativ.UnmarshalSigned(f.Data))
	return nil
}

// AMK1_SetPoints1Reader provides read access to a AMK1_SetPoints1 message.
type AMK1_SetPoints1Reader interface {
	// AMK_bInverterOn returns the value of the AMK_bInverterOn signal.
	AMK_bInverterOn() bool
	// AMK_bDcOn returns the value of the AMK_bDcOn signal.
	AMK_bDcOn() bool
	// AMK_bEnable returns the value of the AMK_bEnable signal.
	AMK_bEnable() bool
	// AMK_bErrorReset returns the value of the AMK_bErrorReset signal.
	AMK_bErrorReset() bool
	// AMK_TargetVelocity returns the value of the AMK_TargetVelocity signal.
	AMK_TargetVelocity() int16
	// AMK_TorqueLimitPositiv returns the value of the AMK_TorqueLimitPositiv signal.
	AMK_TorqueLimitPositiv() int16
	// AMK_TorqueLimitNegativ returns the value of the AMK_TorqueLimitNegativ signal.
	AMK_TorqueLimitNegativ() int16
}

// AMK1_SetPoints1Writer provides write access to a AMK1_SetPoints1 message.
type AMK1_SetPoints1Writer interface {
	// CopyFrom copies all values from AMK1_SetPoints1Reader.
	CopyFrom(AMK1_SetPoints1Reader) *AMK1_SetPoints1
	// SetAMK_bInverterOn sets the value of the AMK_bInverterOn signal.
	SetAMK_bInverterOn(bool) *AMK1_SetPoints1
	// SetAMK_bDcOn sets the value of the AMK_bDcOn signal.
	SetAMK_bDcOn(bool) *AMK1_SetPoints1
	// SetAMK_bEnable sets the value of the AMK_bEnable signal.
	SetAMK_bEnable(bool) *AMK1_SetPoints1
	// SetAMK_bErrorReset sets the value of the AMK_bErrorReset signal.
	SetAMK_bErrorReset(bool) *AMK1_SetPoints1
	// SetAMK_TargetVelocity sets the value of the AMK_TargetVelocity signal.
	SetAMK_TargetVelocity(int16) *AMK1_SetPoints1
	// SetAMK_TorqueLimitPositiv sets the value of the AMK_TorqueLimitPositiv signal.
	SetAMK_TorqueLimitPositiv(int16) *AMK1_SetPoints1
	// SetAMK_TorqueLimitNegativ sets the value of the AMK_TorqueLimitNegativ signal.
	SetAMK_TorqueLimitNegativ(int16) *AMK1_SetPoints1
}

type AMK1_SetPoints1 struct {
	xxx_AMK_bInverterOn        bool
	xxx_AMK_bDcOn              bool
	xxx_AMK_bEnable            bool
	xxx_AMK_bErrorReset        bool
	xxx_AMK_TargetVelocity     int16
	xxx_AMK_TorqueLimitPositiv int16
	xxx_AMK_TorqueLimitNegativ int16
}

func NewAMK1_SetPoints1() *AMK1_SetPoints1 {
	m := &AMK1_SetPoints1{}
	m.Reset()
	return m
}

func (m *AMK1_SetPoints1) Reset() {
	m.xxx_AMK_bInverterOn = false
	m.xxx_AMK_bDcOn = false
	m.xxx_AMK_bEnable = false
	m.xxx_AMK_bErrorReset = false
	m.xxx_AMK_TargetVelocity = 0
	m.xxx_AMK_TorqueLimitPositiv = 0
	m.xxx_AMK_TorqueLimitNegativ = 0
}

func (m *AMK1_SetPoints1) CopyFrom(o AMK1_SetPoints1Reader) *AMK1_SetPoints1 {
	m.xxx_AMK_bInverterOn = o.AMK_bInverterOn()
	m.xxx_AMK_bDcOn = o.AMK_bDcOn()
	m.xxx_AMK_bEnable = o.AMK_bEnable()
	m.xxx_AMK_bErrorReset = o.AMK_bErrorReset()
	m.xxx_AMK_TargetVelocity = o.AMK_TargetVelocity()
	m.xxx_AMK_TorqueLimitPositiv = o.AMK_TorqueLimitPositiv()
	m.xxx_AMK_TorqueLimitNegativ = o.AMK_TorqueLimitNegativ()
	return m
}

// Descriptor returns the AMK1_SetPoints1 descriptor.
func (m *AMK1_SetPoints1) Descriptor() *descriptor.Message {
	return Messages().AMK1_SetPoints1.Message
}

// String returns a compact string representation of the message.
func (m *AMK1_SetPoints1) String() string {
	return cantext.MessageString(m)
}

func (m *AMK1_SetPoints1) AMK_bInverterOn() bool {
	return m.xxx_AMK_bInverterOn
}

func (m *AMK1_SetPoints1) SetAMK_bInverterOn(v bool) *AMK1_SetPoints1 {
	m.xxx_AMK_bInverterOn = v
	return m
}

func (m *AMK1_SetPoints1) AMK_bDcOn() bool {
	return m.xxx_AMK_bDcOn
}

func (m *AMK1_SetPoints1) SetAMK_bDcOn(v bool) *AMK1_SetPoints1 {
	m.xxx_AMK_bDcOn = v
	return m
}

func (m *AMK1_SetPoints1) AMK_bEnable() bool {
	return m.xxx_AMK_bEnable
}

func (m *AMK1_SetPoints1) SetAMK_bEnable(v bool) *AMK1_SetPoints1 {
	m.xxx_AMK_bEnable = v
	return m
}

func (m *AMK1_SetPoints1) AMK_bErrorReset() bool {
	return m.xxx_AMK_bErrorReset
}

func (m *AMK1_SetPoints1) SetAMK_bErrorReset(v bool) *AMK1_SetPoints1 {
	m.xxx_AMK_bErrorReset = v
	return m
}

func (m *AMK1_SetPoints1) AMK_TargetVelocity() int16 {
	return m.xxx_AMK_TargetVelocity
}

func (m *AMK1_SetPoints1) SetAMK_TargetVelocity(v int16) *AMK1_SetPoints1 {
	m.xxx_AMK_TargetVelocity = int16(Messages().AMK1_SetPoints1.AMK_TargetVelocity.SaturatedCastSigned(int64(v)))
	return m
}

func (m *AMK1_SetPoints1) AMK_TorqueLimitPositiv() int16 {
	return m.xxx_AMK_TorqueLimitPositiv
}

func (m *AMK1_SetPoints1) SetAMK_TorqueLimitPositiv(v int16) *AMK1_SetPoints1 {
	m.xxx_AMK_TorqueLimitPositiv = int16(Messages().AMK1_SetPoints1.AMK_TorqueLimitPositiv.SaturatedCastSigned(int64(v)))
	return m
}

func (m *AMK1_SetPoints1) AMK_TorqueLimitNegativ() int16 {
	return m.xxx_AMK_TorqueLimitNegativ
}

func (m *AMK1_SetPoints1) SetAMK_TorqueLimitNegativ(v int16) *AMK1_SetPoints1 {
	m.xxx_AMK_TorqueLimitNegativ = int16(Messages().AMK1_SetPoints1.AMK_TorqueLimitNegativ.SaturatedCastSigned(int64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *AMK1_SetPoints1) Frame() can.Frame {
	md := Messages().AMK1_SetPoints1
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.AMK_bInverterOn.MarshalBool(&f.Data, bool(m.xxx_AMK_bInverterOn))
	md.AMK_bDcOn.MarshalBool(&f.Data, bool(m.xxx_AMK_bDcOn))
	md.AMK_bEnable.MarshalBool(&f.Data, bool(m.xxx_AMK_bEnable))
	md.AMK_bErrorReset.MarshalBool(&f.Data, bool(m.xxx_AMK_bErrorReset))
	md.AMK_TargetVelocity.MarshalSigned(&f.Data, int64(m.xxx_AMK_TargetVelocity))
	md.AMK_TorqueLimitPositiv.MarshalSigned(&f.Data, int64(m.xxx_AMK_TorqueLimitPositiv))
	md.AMK_TorqueLimitNegativ.MarshalSigned(&f.Data, int64(m.xxx_AMK_TorqueLimitNegativ))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *AMK1_SetPoints1) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *AMK1_SetPoints1) UnmarshalFrame(f can.Frame) error {
	md := Messages().AMK1_SetPoints1
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal AMK1_SetPoints1: expects ID 390 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal AMK1_SetPoints1: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal AMK1_SetPoints1: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal AMK1_SetPoints1: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_AMK_bInverterOn = bool(md.AMK_bInverterOn.UnmarshalBool(f.Data))
	m.xxx_AMK_bDcOn = bool(md.AMK_bDcOn.UnmarshalBool(f.Data))
	m.xxx_AMK_bEnable = bool(md.AMK_bEnable.UnmarshalBool(f.Data))
	m.xxx_AMK_bErrorReset = bool(md.AMK_bErrorReset.UnmarshalBool(f.Data))
	m.xxx_AMK_TargetVelocity = int16(md.AMK_TargetVelocity.UnmarshalSigned(f.Data))
	m.xxx_AMK_TorqueLimitPositiv = int16(md.AMK_TorqueLimitPositiv.UnmarshalSigned(f.Data))
	m.xxx_AMK_TorqueLimitNegativ = int16(md.AMK_TorqueLimitNegativ.UnmarshalSigned(f.Data))
	return nil
}

// AMK0_ActualValues1Reader provides read access to a AMK0_ActualValues1 message.
type AMK0_ActualValues1Reader interface {
	// AMK_bSystemReady returns the value of the AMK_bSystemReady signal.
	AMK_bSystemReady() bool
	// AMK_bError returns the value of the AMK_bError signal.
	AMK_bError() bool
	// AMK_bWarn returns the value of the AMK_bWarn signal.
	AMK_bWarn() bool
	// AMK_bQuitDcOn returns the value of the AMK_bQuitDcOn signal.
	AMK_bQuitDcOn() bool
	// AMK_bDcOn returns the value of the AMK_bDcOn signal.
	AMK_bDcOn() bool
	// AMK_bQuitInverterOn returns the value of the AMK_bQuitInverterOn signal.
	AMK_bQuitInverterOn() bool
	// AMK_bInverterOn returns the value of the AMK_bInverterOn signal.
	AMK_bInverterOn() bool
	// AMK_bDerating returns the value of the AMK_bDerating signal.
	AMK_bDerating() bool
	// AMK_ActualVelocity returns the value of the AMK_ActualVelocity signal.
	AMK_ActualVelocity() int16
	// AMK_TorqueCurrent returns the value of the AMK_TorqueCurrent signal.
	AMK_TorqueCurrent() int16
	// AMK_MagnetizingCurrent returns the value of the AMK_MagnetizingCurrent signal.
	AMK_MagnetizingCurrent() int16
}

// AMK0_ActualValues1Writer provides write access to a AMK0_ActualValues1 message.
type AMK0_ActualValues1Writer interface {
	// CopyFrom copies all values from AMK0_ActualValues1Reader.
	CopyFrom(AMK0_ActualValues1Reader) *AMK0_ActualValues1
	// SetAMK_bSystemReady sets the value of the AMK_bSystemReady signal.
	SetAMK_bSystemReady(bool) *AMK0_ActualValues1
	// SetAMK_bError sets the value of the AMK_bError signal.
	SetAMK_bError(bool) *AMK0_ActualValues1
	// SetAMK_bWarn sets the value of the AMK_bWarn signal.
	SetAMK_bWarn(bool) *AMK0_ActualValues1
	// SetAMK_bQuitDcOn sets the value of the AMK_bQuitDcOn signal.
	SetAMK_bQuitDcOn(bool) *AMK0_ActualValues1
	// SetAMK_bDcOn sets the value of the AMK_bDcOn signal.
	SetAMK_bDcOn(bool) *AMK0_ActualValues1
	// SetAMK_bQuitInverterOn sets the value of the AMK_bQuitInverterOn signal.
	SetAMK_bQuitInverterOn(bool) *AMK0_ActualValues1
	// SetAMK_bInverterOn sets the value of the AMK_bInverterOn signal.
	SetAMK_bInverterOn(bool) *AMK0_ActualValues1
	// SetAMK_bDerating sets the value of the AMK_bDerating signal.
	SetAMK_bDerating(bool) *AMK0_ActualValues1
	// SetAMK_ActualVelocity sets the value of the AMK_ActualVelocity signal.
	SetAMK_ActualVelocity(int16) *AMK0_ActualValues1
	// SetAMK_TorqueCurrent sets the value of the AMK_TorqueCurrent signal.
	SetAMK_TorqueCurrent(int16) *AMK0_ActualValues1
	// SetAMK_MagnetizingCurrent sets the value of the AMK_MagnetizingCurrent signal.
	SetAMK_MagnetizingCurrent(int16) *AMK0_ActualValues1
}

type AMK0_ActualValues1 struct {
	xxx_AMK_bSystemReady       bool
	xxx_AMK_bError             bool
	xxx_AMK_bWarn              bool
	xxx_AMK_bQuitDcOn          bool
	xxx_AMK_bDcOn              bool
	xxx_AMK_bQuitInverterOn    bool
	xxx_AMK_bInverterOn        bool
	xxx_AMK_bDerating          bool
	xxx_AMK_ActualVelocity     int16
	xxx_AMK_TorqueCurrent      int16
	xxx_AMK_MagnetizingCurrent int16
}

func NewAMK0_ActualValues1() *AMK0_ActualValues1 {
	m := &AMK0_ActualValues1{}
	m.Reset()
	return m
}

func (m *AMK0_ActualValues1) Reset() {
	m.xxx_AMK_bSystemReady = false
	m.xxx_AMK_bError = false
	m.xxx_AMK_bWarn = false
	m.xxx_AMK_bQuitDcOn = false
	m.xxx_AMK_bDcOn = false
	m.xxx_AMK_bQuitInverterOn = false
	m.xxx_AMK_bInverterOn = false
	m.xxx_AMK_bDerating = false
	m.xxx_AMK_ActualVelocity = 0
	m.xxx_AMK_TorqueCurrent = 0
	m.xxx_AMK_MagnetizingCurrent = 0
}

func (m *AMK0_ActualValues1) CopyFrom(o AMK0_ActualValues1Reader) *AMK0_ActualValues1 {
	m.xxx_AMK_bSystemReady = o.AMK_bSystemReady()
	m.xxx_AMK_bError = o.AMK_bError()
	m.xxx_AMK_bWarn = o.AMK_bWarn()
	m.xxx_AMK_bQuitDcOn = o.AMK_bQuitDcOn()
	m.xxx_AMK_bDcOn = o.AMK_bDcOn()
	m.xxx_AMK_bQuitInverterOn = o.AMK_bQuitInverterOn()
	m.xxx_AMK_bInverterOn = o.AMK_bInverterOn()
	m.xxx_AMK_bDerating = o.AMK_bDerating()
	m.xxx_AMK_ActualVelocity = o.AMK_ActualVelocity()
	m.xxx_AMK_TorqueCurrent = o.AMK_TorqueCurrent()
	m.xxx_AMK_MagnetizingCurrent = o.AMK_MagnetizingCurrent()
	return m
}

// Descriptor returns the AMK0_ActualValues1 descriptor.
func (m *AMK0_ActualValues1) Descriptor() *descriptor.Message {
	return Messages().AMK0_ActualValues1.Message
}

// String returns a compact string representation of the message.
func (m *AMK0_ActualValues1) String() string {
	return cantext.MessageString(m)
}

func (m *AMK0_ActualValues1) AMK_bSystemReady() bool {
	return m.xxx_AMK_bSystemReady
}

func (m *AMK0_ActualValues1) SetAMK_bSystemReady(v bool) *AMK0_ActualValues1 {
	m.xxx_AMK_bSystemReady = v
	return m
}

func (m *AMK0_ActualValues1) AMK_bError() bool {
	return m.xxx_AMK_bError
}

func (m *AMK0_ActualValues1) SetAMK_bError(v bool) *AMK0_ActualValues1 {
	m.xxx_AMK_bError = v
	return m
}

func (m *AMK0_ActualValues1) AMK_bWarn() bool {
	return m.xxx_AMK_bWarn
}

func (m *AMK0_ActualValues1) SetAMK_bWarn(v bool) *AMK0_ActualValues1 {
	m.xxx_AMK_bWarn = v
	return m
}

func (m *AMK0_ActualValues1) AMK_bQuitDcOn() bool {
	return m.xxx_AMK_bQuitDcOn
}

func (m *AMK0_ActualValues1) SetAMK_bQuitDcOn(v bool) *AMK0_ActualValues1 {
	m.xxx_AMK_bQuitDcOn = v
	return m
}

func (m *AMK0_ActualValues1) AMK_bDcOn() bool {
	return m.xxx_AMK_bDcOn
}

func (m *AMK0_ActualValues1) SetAMK_bDcOn(v bool) *AMK0_ActualValues1 {
	m.xxx_AMK_bDcOn = v
	return m
}

func (m *AMK0_ActualValues1) AMK_bQuitInverterOn() bool {
	return m.xxx_AMK_bQuitInverterOn
}

func (m *AMK0_ActualValues1) SetAMK_bQuitInverterOn(v bool) *AMK0_ActualValues1 {
	m.xxx_AMK_bQuitInverterOn = v
	return m
}

func (m *AMK0_ActualValues1) AMK_bInverterOn() bool {
	return m.xxx_AMK_bInverterOn
}

func (m *AMK0_ActualValues1) SetAMK_bInverterOn(v bool) *AMK0_ActualValues1 {
	m.xxx_AMK_bInverterOn = v
	return m
}

func (m *AMK0_ActualValues1) AMK_bDerating() bool {
	return m.xxx_AMK_bDerating
}

func (m *AMK0_ActualValues1) SetAMK_bDerating(v bool) *AMK0_ActualValues1 {
	m.xxx_AMK_bDerating = v
	return m
}

func (m *AMK0_ActualValues1) AMK_ActualVelocity() int16 {
	return m.xxx_AMK_ActualVelocity
}

func (m *AMK0_ActualValues1) SetAMK_ActualVelocity(v int16) *AMK0_ActualValues1 {
	m.xxx_AMK_ActualVelocity = int16(Messages().AMK0_ActualValues1.AMK_ActualVelocity.SaturatedCastSigned(int64(v)))
	return m
}

func (m *AMK0_ActualValues1) AMK_TorqueCurrent() int16 {
	return m.xxx_AMK_TorqueCurrent
}

func (m *AMK0_ActualValues1) SetAMK_TorqueCurrent(v int16) *AMK0_ActualValues1 {
	m.xxx_AMK_TorqueCurrent = int16(Messages().AMK0_ActualValues1.AMK_TorqueCurrent.SaturatedCastSigned(int64(v)))
	return m
}

func (m *AMK0_ActualValues1) AMK_MagnetizingCurrent() int16 {
	return m.xxx_AMK_MagnetizingCurrent
}

func (m *AMK0_ActualValues1) SetAMK_MagnetizingCurrent(v int16) *AMK0_ActualValues1 {
	m.xxx_AMK_MagnetizingCurrent = int16(Messages().AMK0_ActualValues1.AMK_MagnetizingCurrent.SaturatedCastSigned(int64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *AMK0_ActualValues1) Frame() can.Frame {
	md := Messages().AMK0_ActualValues1
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.AMK_bSystemReady.MarshalBool(&f.Data, bool(m.xxx_AMK_bSystemReady))
	md.AMK_bError.MarshalBool(&f.Data, bool(m.xxx_AMK_bError))
	md.AMK_bWarn.MarshalBool(&f.Data, bool(m.xxx_AMK_bWarn))
	md.AMK_bQuitDcOn.MarshalBool(&f.Data, bool(m.xxx_AMK_bQuitDcOn))
	md.AMK_bDcOn.MarshalBool(&f.Data, bool(m.xxx_AMK_bDcOn))
	md.AMK_bQuitInverterOn.MarshalBool(&f.Data, bool(m.xxx_AMK_bQuitInverterOn))
	md.AMK_bInverterOn.MarshalBool(&f.Data, bool(m.xxx_AMK_bInverterOn))
	md.AMK_bDerating.MarshalBool(&f.Data, bool(m.xxx_AMK_bDerating))
	md.AMK_ActualVelocity.MarshalSigned(&f.Data, int64(m.xxx_AMK_ActualVelocity))
	md.AMK_TorqueCurrent.MarshalSigned(&f.Data, int64(m.xxx_AMK_TorqueCurrent))
	md.AMK_MagnetizingCurrent.MarshalSigned(&f.Data, int64(m.xxx_AMK_MagnetizingCurrent))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *AMK0_ActualValues1) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *AMK0_ActualValues1) UnmarshalFrame(f can.Frame) error {
	md := Messages().AMK0_ActualValues1
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal AMK0_ActualValues1: expects ID 644 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal AMK0_ActualValues1: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal AMK0_ActualValues1: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal AMK0_ActualValues1: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_AMK_bSystemReady = bool(md.AMK_bSystemReady.UnmarshalBool(f.Data))
	m.xxx_AMK_bError = bool(md.AMK_bError.UnmarshalBool(f.Data))
	m.xxx_AMK_bWarn = bool(md.AMK_bWarn.UnmarshalBool(f.Data))
	m.xxx_AMK_bQuitDcOn = bool(md.AMK_bQuitDcOn.UnmarshalBool(f.Data))
	m.xxx_AMK_bDcOn = bool(md.AMK_bDcOn.UnmarshalBool(f.Data))
	m.xxx_AMK_bQuitInverterOn = bool(md.AMK_bQuitInverterOn.UnmarshalBool(f.Data))
	m.xxx_AMK_bInverterOn = bool(md.AMK_bInverterOn.UnmarshalBool(f.Data))
	m.xxx_AMK_bDerating = bool(md.AMK_bDerating.UnmarshalBool(f.Data))
	m.xxx_AMK_ActualVelocity = int16(md.AMK_ActualVelocity.UnmarshalSigned(f.Data))
	m.xxx_AMK_TorqueCurrent = int16(md.AMK_TorqueCurrent.UnmarshalSigned(f.Data))
	m.xxx_AMK_MagnetizingCurrent = int16(md.AMK_MagnetizingCurrent.UnmarshalSigned(f.Data))
	return nil
}

// AMK1_ActualValues1Reader provides read access to a AMK1_ActualValues1 message.
type AMK1_ActualValues1Reader interface {
	// AMK_bSystemReady returns the value of the AMK_bSystemReady signal.
	AMK_bSystemReady() bool
	// AMK_bError returns the value of the AMK_bError signal.
	AMK_bError() bool
	// AMK_bWarn returns the value of the AMK_bWarn signal.
	AMK_bWarn() bool
	// AMK_bQuitDcOn returns the value of the AMK_bQuitDcOn signal.
	AMK_bQuitDcOn() bool
	// AMK_bDcOn returns the value of the AMK_bDcOn signal.
	AMK_bDcOn() bool
	// AMK_bQuitInverterOn returns the value of the AMK_bQuitInverterOn signal.
	AMK_bQuitInverterOn() bool
	// AMK_bInverterOn returns the value of the AMK_bInverterOn signal.
	AMK_bInverterOn() bool
	// AMK_bDerating returns the value of the AMK_bDerating signal.
	AMK_bDerating() bool
	// AMK_ActualVelocity returns the value of the AMK_ActualVelocity signal.
	AMK_ActualVelocity() int16
	// AMK_TorqueCurrent returns the value of the AMK_TorqueCurrent signal.
	AMK_TorqueCurrent() int16
	// AMK_MagnetizingCurrent returns the value of the AMK_MagnetizingCurrent signal.
	AMK_MagnetizingCurrent() int16
}

// AMK1_ActualValues1Writer provides write access to a AMK1_ActualValues1 message.
type AMK1_ActualValues1Writer interface {
	// CopyFrom copies all values from AMK1_ActualValues1Reader.
	CopyFrom(AMK1_ActualValues1Reader) *AMK1_ActualValues1
	// SetAMK_bSystemReady sets the value of the AMK_bSystemReady signal.
	SetAMK_bSystemReady(bool) *AMK1_ActualValues1
	// SetAMK_bError sets the value of the AMK_bError signal.
	SetAMK_bError(bool) *AMK1_ActualValues1
	// SetAMK_bWarn sets the value of the AMK_bWarn signal.
	SetAMK_bWarn(bool) *AMK1_ActualValues1
	// SetAMK_bQuitDcOn sets the value of the AMK_bQuitDcOn signal.
	SetAMK_bQuitDcOn(bool) *AMK1_ActualValues1
	// SetAMK_bDcOn sets the value of the AMK_bDcOn signal.
	SetAMK_bDcOn(bool) *AMK1_ActualValues1
	// SetAMK_bQuitInverterOn sets the value of the AMK_bQuitInverterOn signal.
	SetAMK_bQuitInverterOn(bool) *AMK1_ActualValues1
	// SetAMK_bInverterOn sets the value of the AMK_bInverterOn signal.
	SetAMK_bInverterOn(bool) *AMK1_ActualValues1
	// SetAMK_bDerating sets the value of the AMK_bDerating signal.
	SetAMK_bDerating(bool) *AMK1_ActualValues1
	// SetAMK_ActualVelocity sets the value of the AMK_ActualVelocity signal.
	SetAMK_ActualVelocity(int16) *AMK1_ActualValues1
	// SetAMK_TorqueCurrent sets the value of the AMK_TorqueCurrent signal.
	SetAMK_TorqueCurrent(int16) *AMK1_ActualValues1
	// SetAMK_MagnetizingCurrent sets the value of the AMK_MagnetizingCurrent signal.
	SetAMK_MagnetizingCurrent(int16) *AMK1_ActualValues1
}

type AMK1_ActualValues1 struct {
	xxx_AMK_bSystemReady       bool
	xxx_AMK_bError             bool
	xxx_AMK_bWarn              bool
	xxx_AMK_bQuitDcOn          bool
	xxx_AMK_bDcOn              bool
	xxx_AMK_bQuitInverterOn    bool
	xxx_AMK_bInverterOn        bool
	xxx_AMK_bDerating          bool
	xxx_AMK_ActualVelocity     int16
	xxx_AMK_TorqueCurrent      int16
	xxx_AMK_MagnetizingCurrent int16
}

func NewAMK1_ActualValues1() *AMK1_ActualValues1 {
	m := &AMK1_ActualValues1{}
	m.Reset()
	return m
}

func (m *AMK1_ActualValues1) Reset() {
	m.xxx_AMK_bSystemReady = false
	m.xxx_AMK_bError = false
	m.xxx_AMK_bWarn = false
	m.xxx_AMK_bQuitDcOn = false
	m.xxx_AMK_bDcOn = false
	m.xxx_AMK_bQuitInverterOn = false
	m.xxx_AMK_bInverterOn = false
	m.xxx_AMK_bDerating = false
	m.xxx_AMK_ActualVelocity = 0
	m.xxx_AMK_TorqueCurrent = 0
	m.xxx_AMK_MagnetizingCurrent = 0
}

func (m *AMK1_ActualValues1) CopyFrom(o AMK1_ActualValues1Reader) *AMK1_ActualValues1 {
	m.xxx_AMK_bSystemReady = o.AMK_bSystemReady()
	m.xxx_AMK_bError = o.AMK_bError()
	m.xxx_AMK_bWarn = o.AMK_bWarn()
	m.xxx_AMK_bQuitDcOn = o.AMK_bQuitDcOn()
	m.xxx_AMK_bDcOn = o.AMK_bDcOn()
	m.xxx_AMK_bQuitInverterOn = o.AMK_bQuitInverterOn()
	m.xxx_AMK_bInverterOn = o.AMK_bInverterOn()
	m.xxx_AMK_bDerating = o.AMK_bDerating()
	m.xxx_AMK_ActualVelocity = o.AMK_ActualVelocity()
	m.xxx_AMK_TorqueCurrent = o.AMK_TorqueCurrent()
	m.xxx_AMK_MagnetizingCurrent = o.AMK_MagnetizingCurrent()
	return m
}

// Descriptor returns the AMK1_ActualValues1 descriptor.
func (m *AMK1_ActualValues1) Descriptor() *descriptor.Message {
	return Messages().AMK1_ActualValues1.Message
}

// String returns a compact string representation of the message.
func (m *AMK1_ActualValues1) String() string {
	return cantext.MessageString(m)
}

func (m *AMK1_ActualValues1) AMK_bSystemReady() bool {
	return m.xxx_AMK_bSystemReady
}

func (m *AMK1_ActualValues1) SetAMK_bSystemReady(v bool) *AMK1_ActualValues1 {
	m.xxx_AMK_bSystemReady = v
	return m
}

func (m *AMK1_ActualValues1) AMK_bError() bool {
	return m.xxx_AMK_bError
}

func (m *AMK1_ActualValues1) SetAMK_bError(v bool) *AMK1_ActualValues1 {
	m.xxx_AMK_bError = v
	return m
}

func (m *AMK1_ActualValues1) AMK_bWarn() bool {
	return m.xxx_AMK_bWarn
}

func (m *AMK1_ActualValues1) SetAMK_bWarn(v bool) *AMK1_ActualValues1 {
	m.xxx_AMK_bWarn = v
	return m
}

func (m *AMK1_ActualValues1) AMK_bQuitDcOn() bool {
	return m.xxx_AMK_bQuitDcOn
}

func (m *AMK1_ActualValues1) SetAMK_bQuitDcOn(v bool) *AMK1_ActualValues1 {
	m.xxx_AMK_bQuitDcOn = v
	return m
}

func (m *AMK1_ActualValues1) AMK_bDcOn() bool {
	return m.xxx_AMK_bDcOn
}

func (m *AMK1_ActualValues1) SetAMK_bDcOn(v bool) *AMK1_ActualValues1 {
	m.xxx_AMK_bDcOn = v
	return m
}

func (m *AMK1_ActualValues1) AMK_bQuitInverterOn() bool {
	return m.xxx_AMK_bQuitInverterOn
}

func (m *AMK1_ActualValues1) SetAMK_bQuitInverterOn(v bool) *AMK1_ActualValues1 {
	m.xxx_AMK_bQuitInverterOn = v
	return m
}

func (m *AMK1_ActualValues1) AMK_bInverterOn() bool {
	return m.xxx_AMK_bInverterOn
}

func (m *AMK1_ActualValues1) SetAMK_bInverterOn(v bool) *AMK1_ActualValues1 {
	m.xxx_AMK_bInverterOn = v
	return m
}

func (m *AMK1_ActualValues1) AMK_bDerating() bool {
	return m.xxx_AMK_bDerating
}

func (m *AMK1_ActualValues1) SetAMK_bDerating(v bool) *AMK1_ActualValues1 {
	m.xxx_AMK_bDerating = v
	return m
}

func (m *AMK1_ActualValues1) AMK_ActualVelocity() int16 {
	return m.xxx_AMK_ActualVelocity
}

func (m *AMK1_ActualValues1) SetAMK_ActualVelocity(v int16) *AMK1_ActualValues1 {
	m.xxx_AMK_ActualVelocity = int16(Messages().AMK1_ActualValues1.AMK_ActualVelocity.SaturatedCastSigned(int64(v)))
	return m
}

func (m *AMK1_ActualValues1) AMK_TorqueCurrent() int16 {
	return m.xxx_AMK_TorqueCurrent
}

func (m *AMK1_ActualValues1) SetAMK_TorqueCurrent(v int16) *AMK1_ActualValues1 {
	m.xxx_AMK_TorqueCurrent = int16(Messages().AMK1_ActualValues1.AMK_TorqueCurrent.SaturatedCastSigned(int64(v)))
	return m
}

func (m *AMK1_ActualValues1) AMK_MagnetizingCurrent() int16 {
	return m.xxx_AMK_MagnetizingCurrent
}

func (m *AMK1_ActualValues1) SetAMK_MagnetizingCurrent(v int16) *AMK1_ActualValues1 {
	m.xxx_AMK_MagnetizingCurrent = int16(Messages().AMK1_ActualValues1.AMK_MagnetizingCurrent.SaturatedCastSigned(int64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *AMK1_ActualValues1) Frame() can.Frame {
	md := Messages().AMK1_ActualValues1
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.AMK_bSystemReady.MarshalBool(&f.Data, bool(m.xxx_AMK_bSystemReady))
	md.AMK_bError.MarshalBool(&f.Data, bool(m.xxx_AMK_bError))
	md.AMK_bWarn.MarshalBool(&f.Data, bool(m.xxx_AMK_bWarn))
	md.AMK_bQuitDcOn.MarshalBool(&f.Data, bool(m.xxx_AMK_bQuitDcOn))
	md.AMK_bDcOn.MarshalBool(&f.Data, bool(m.xxx_AMK_bDcOn))
	md.AMK_bQuitInverterOn.MarshalBool(&f.Data, bool(m.xxx_AMK_bQuitInverterOn))
	md.AMK_bInverterOn.MarshalBool(&f.Data, bool(m.xxx_AMK_bInverterOn))
	md.AMK_bDerating.MarshalBool(&f.Data, bool(m.xxx_AMK_bDerating))
	md.AMK_ActualVelocity.MarshalSigned(&f.Data, int64(m.xxx_AMK_ActualVelocity))
	md.AMK_TorqueCurrent.MarshalSigned(&f.Data, int64(m.xxx_AMK_TorqueCurrent))
	md.AMK_MagnetizingCurrent.MarshalSigned(&f.Data, int64(m.xxx_AMK_MagnetizingCurrent))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *AMK1_ActualValues1) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *AMK1_ActualValues1) UnmarshalFrame(f can.Frame) error {
	md := Messages().AMK1_ActualValues1
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal AMK1_ActualValues1: expects ID 645 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal AMK1_ActualValues1: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal AMK1_ActualValues1: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal AMK1_ActualValues1: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_AMK_bSystemReady = bool(md.AMK_bSystemReady.UnmarshalBool(f.Data))
	m.xxx_AMK_bError = bool(md.AMK_bError.UnmarshalBool(f.Data))
	m.xxx_AMK_bWarn = bool(md.AMK_bWarn.UnmarshalBool(f.Data))
	m.xxx_AMK_bQuitDcOn = bool(md.AMK_bQuitDcOn.UnmarshalBool(f.Data))
	m.xxx_AMK_bDcOn = bool(md.AMK_bDcOn.UnmarshalBool(f.Data))
	m.xxx_AMK_bQuitInverterOn = bool(md.AMK_bQuitInverterOn.UnmarshalBool(f.Data))
	m.xxx_AMK_bInverterOn = bool(md.AMK_bInverterOn.UnmarshalBool(f.Data))
	m.xxx_AMK_bDerating = bool(md.AMK_bDerating.UnmarshalBool(f.Data))
	m.xxx_AMK_ActualVelocity = int16(md.AMK_ActualVelocity.UnmarshalSigned(f.Data))
	m.xxx_AMK_TorqueCurrent = int16(md.AMK_TorqueCurrent.UnmarshalSigned(f.Data))
	m.xxx_AMK_MagnetizingCurrent = int16(md.AMK_MagnetizingCurrent.UnmarshalSigned(f.Data))
	return nil
}

// AMK0_ActualValues2Reader provides read access to a AMK0_ActualValues2 message.
type AMK0_ActualValues2Reader interface {
	// AMK_TempMotor returns the physical value of the AMK_TempMotor signal.
	AMK_TempMotor() float64
	// AMK_TempInverter returns the physical value of the AMK_TempInverter signal.
	AMK_TempInverter() float64
	// AMK_ErrorInfo returns the value of the AMK_ErrorInfo signal.
	AMK_ErrorInfo() uint16
	// AMK_TempIGBT returns the physical value of the AMK_TempIGBT signal.
	AMK_TempIGBT() float64
}

// AMK0_ActualValues2Writer provides write access to a AMK0_ActualValues2 message.
type AMK0_ActualValues2Writer interface {
	// CopyFrom copies all values from AMK0_ActualValues2Reader.
	CopyFrom(AMK0_ActualValues2Reader) *AMK0_ActualValues2
	// SetAMK_TempMotor sets the physical value of the AMK_TempMotor signal.
	SetAMK_TempMotor(float64) *AMK0_ActualValues2
	// SetAMK_TempInverter sets the physical value of the AMK_TempInverter signal.
	SetAMK_TempInverter(float64) *AMK0_ActualValues2
	// SetAMK_ErrorInfo sets the value of the AMK_ErrorInfo signal.
	SetAMK_ErrorInfo(uint16) *AMK0_ActualValues2
	// SetAMK_TempIGBT sets the physical value of the AMK_TempIGBT signal.
	SetAMK_TempIGBT(float64) *AMK0_ActualValues2
}

type AMK0_ActualValues2 struct {
	xxx_AMK_TempMotor    int16
	xxx_AMK_TempInverter int16
	xxx_AMK_ErrorInfo    uint16
	xxx_AMK_TempIGBT     int16
}

func NewAMK0_ActualValues2() *AMK0_ActualValues2 {
	m := &AMK0_ActualValues2{}
	m.Reset()
	return m
}

func (m *AMK0_ActualValues2) Reset() {
	m.xxx_AMK_TempMotor = 0
	m.xxx_AMK_TempInverter = 0
	m.xxx_AMK_ErrorInfo = 0
	m.xxx_AMK_TempIGBT = 0
}

func (m *AMK0_ActualValues2) CopyFrom(o AMK0_ActualValues2Reader) *AMK0_ActualValues2 {
	m.SetAMK_TempMotor(o.AMK_TempMotor())
	m.SetAMK_TempInverter(o.AMK_TempInverter())
	m.xxx_AMK_ErrorInfo = o.AMK_ErrorInfo()
	m.SetAMK_TempIGBT(o.AMK_TempIGBT())
	return m
}

// Descriptor returns the AMK0_ActualValues2 descriptor.
func (m *AMK0_ActualValues2) Descriptor() *descriptor.Message {
	return Messages().AMK0_ActualValues2.Message
}

// String returns a compact string representation of the message.
func (m *AMK0_ActualValues2) String() string {
	return cantext.MessageString(m)
}

func (m *AMK0_ActualValues2) AMK_TempMotor() float64 {
	return Messages().AMK0_ActualValues2.AMK_TempMotor.ToPhysical(float64(m.xxx_AMK_TempMotor))
}

func (m *AMK0_ActualValues2) SetAMK_TempMotor(v float64) *AMK0_ActualValues2 {
	m.xxx_AMK_TempMotor = int16(Messages().AMK0_ActualValues2.AMK_TempMotor.FromPhysical(v))
	return m
}

func (m *AMK0_ActualValues2) AMK_TempInverter() float64 {
	return Messages().AMK0_ActualValues2.AMK_TempInverter.ToPhysical(float64(m.xxx_AMK_TempInverter))
}

func (m *AMK0_ActualValues2) SetAMK_TempInverter(v float64) *AMK0_ActualValues2 {
	m.xxx_AMK_TempInverter = int16(Messages().AMK0_ActualValues2.AMK_TempInverter.FromPhysical(v))
	return m
}

func (m *AMK0_ActualValues2) AMK_ErrorInfo() uint16 {
	return m.xxx_AMK_ErrorInfo
}

func (m *AMK0_ActualValues2) SetAMK_ErrorInfo(v uint16) *AMK0_ActualValues2 {
	m.xxx_AMK_ErrorInfo = uint16(Messages().AMK0_ActualValues2.AMK_ErrorInfo.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *AMK0_ActualValues2) AMK_TempIGBT() float64 {
	return Messages().AMK0_ActualValues2.AMK_TempIGBT.ToPhysical(float64(m.xxx_AMK_TempIGBT))
}

func (m *AMK0_ActualValues2) SetAMK_TempIGBT(v float64) *AMK0_ActualValues2 {
	m.xxx_AMK_TempIGBT = int16(Messages().AMK0_ActualValues2.AMK_TempIGBT.FromPhysical(v))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *AMK0_ActualValues2) Frame() can.Frame {
	md := Messages().AMK0_ActualValues2
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.AMK_TempMotor.MarshalSigned(&f.Data, int64(m.xxx_AMK_TempMotor))
	md.AMK_TempInverter.MarshalSigned(&f.Data, int64(m.xxx_AMK_TempInverter))
	md.AMK_ErrorInfo.MarshalUnsigned(&f.Data, uint64(m.xxx_AMK_ErrorInfo))
	md.AMK_TempIGBT.MarshalSigned(&f.Data, int64(m.xxx_AMK_TempIGBT))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *AMK0_ActualValues2) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *AMK0_ActualValues2) UnmarshalFrame(f can.Frame) error {
	md := Messages().AMK0_ActualValues2
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal AMK0_ActualValues2: expects ID 646 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal AMK0_ActualValues2: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal AMK0_ActualValues2: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal AMK0_ActualValues2: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_AMK_TempMotor = int16(md.AMK_TempMotor.UnmarshalSigned(f.Data))
	m.xxx_AMK_TempInverter = int16(md.AMK_TempInverter.UnmarshalSigned(f.Data))
	m.xxx_AMK_ErrorInfo = uint16(md.AMK_ErrorInfo.UnmarshalUnsigned(f.Data))
	m.xxx_AMK_TempIGBT = int16(md.AMK_TempIGBT.UnmarshalSigned(f.Data))
	return nil
}

// AMK1_ActualValues2Reader provides read access to a AMK1_ActualValues2 message.
type AMK1_ActualValues2Reader interface {
	// AMK_TempMotor returns the physical value of the AMK_TempMotor signal.
	AMK_TempMotor() float64
	// AMK_TempInverter returns the physical value of the AMK_TempInverter signal.
	AMK_TempInverter() float64
	// AMK_ErrorInfo returns the value of the AMK_ErrorInfo signal.
	AMK_ErrorInfo() uint16
	// AMK_TempIGBT returns the physical value of the AMK_TempIGBT signal.
	AMK_TempIGBT() float64
}

// AMK1_ActualValues2Writer provides write access to a AMK1_ActualValues2 message.
type AMK1_ActualValues2Writer interface {
	// CopyFrom copies all values from AMK1_ActualValues2Reader.
	CopyFrom(AMK1_ActualValues2Reader) *AMK1_ActualValues2
	// SetAMK_TempMotor sets the physical value of the AMK_TempMotor signal.
	SetAMK_TempMotor(float64) *AMK1_ActualValues2
	// SetAMK_TempInverter sets the physical value of the AMK_TempInverter signal.
	SetAMK_TempInverter(float64) *AMK1_ActualValues2
	// SetAMK_ErrorInfo sets the value of the AMK_ErrorInfo signal.
	SetAMK_ErrorInfo(uint16) *AMK1_ActualValues2
	// SetAMK_TempIGBT sets the physical value of the AMK_TempIGBT signal.
	SetAMK_TempIGBT(float64) *AMK1_ActualValues2
}

type AMK1_ActualValues2 struct {
	xxx_AMK_TempMotor    int16
	xxx_AMK_TempInverter int16
	xxx_AMK_ErrorInfo    uint16
	xxx_AMK_TempIGBT     int16
}

func NewAMK1_ActualValues2() *AMK1_ActualValues2 {
	m := &AMK1_ActualValues2{}
	m.Reset()
	return m
}

func (m *AMK1_ActualValues2) Reset() {
	m.xxx_AMK_TempMotor = 0
	m.xxx_AMK_TempInverter = 0
	m.xxx_AMK_ErrorInfo = 0
	m.xxx_AMK_TempIGBT = 0
}

func (m *AMK1_ActualValues2) CopyFrom(o AMK1_ActualValues2Reader) *AMK1_ActualValues2 {
	m.SetAMK_TempMotor(o.AMK_TempMotor())
	m.SetAMK_TempInverter(o.AMK_TempInverter())
	m.xxx_AMK_ErrorInfo = o.AMK_ErrorInfo()
	m.SetAMK_TempIGBT(o.AMK_TempIGBT())
	return m
}

// Descriptor returns the AMK1_ActualValues2 descriptor.
func (m *AMK1_ActualValues2) Descriptor() *descriptor.Message {
	return Messages().AMK1_ActualValues2.Message
}

// String returns a compact string representation of the message.
func (m *AMK1_ActualValues2) String() string {
	return cantext.MessageString(m)
}

func (m *AMK1_ActualValues2) AMK_TempMotor() float64 {
	return Messages().AMK1_ActualValues2.AMK_TempMotor.ToPhysical(float64(m.xxx_AMK_TempMotor))
}

func (m *AMK1_ActualValues2) SetAMK_TempMotor(v float64) *AMK1_ActualValues2 {
	m.xxx_AMK_TempMotor = int16(Messages().AMK1_ActualValues2.AMK_TempMotor.FromPhysical(v))
	return m
}

func (m *AMK1_ActualValues2) AMK_TempInverter() float64 {
	return Messages().AMK1_ActualValues2.AMK_TempInverter.ToPhysical(float64(m.xxx_AMK_TempInverter))
}

func (m *AMK1_ActualValues2) SetAMK_TempInverter(v float64) *AMK1_ActualValues2 {
	m.xxx_AMK_TempInverter = int16(Messages().AMK1_ActualValues2.AMK_TempInverter.FromPhysical(v))
	return m
}

func (m *AMK1_ActualValues2) AMK_ErrorInfo() uint16 {
	return m.xxx_AMK_ErrorInfo
}

func (m *AMK1_ActualValues2) SetAMK_ErrorInfo(v uint16) *AMK1_ActualValues2 {
	m.xxx_AMK_ErrorInfo = uint16(Messages().AMK1_ActualValues2.AMK_ErrorInfo.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *AMK1_ActualValues2) AMK_TempIGBT() float64 {
	return Messages().AMK1_ActualValues2.AMK_TempIGBT.ToPhysical(float64(m.xxx_AMK_TempIGBT))
}

func (m *AMK1_ActualValues2) SetAMK_TempIGBT(v float64) *AMK1_ActualValues2 {
	m.xxx_AMK_TempIGBT = int16(Messages().AMK1_ActualValues2.AMK_TempIGBT.FromPhysical(v))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *AMK1_ActualValues2) Frame() can.Frame {
	md := Messages().AMK1_ActualValues2
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.AMK_TempMotor.MarshalSigned(&f.Data, int64(m.xxx_AMK_TempMotor))
	md.AMK_TempInverter.MarshalSigned(&f.Data, int64(m.xxx_AMK_TempInverter))
	md.AMK_ErrorInfo.MarshalUnsigned(&f.Data, uint64(m.xxx_AMK_ErrorInfo))
	md.AMK_TempIGBT.MarshalSigned(&f.Data, int64(m.xxx_AMK_TempIGBT))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *AMK1_ActualValues2) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *AMK1_ActualValues2) UnmarshalFrame(f can.Frame) error {
	md := Messages().AMK1_ActualValues2
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal AMK1_ActualValues2: expects ID 647 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal AMK1_ActualValues2: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal AMK1_ActualValues2: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal AMK1_ActualValues2: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_AMK_TempMotor = int16(md.AMK_TempMotor.UnmarshalSigned(f.Data))
	m.xxx_AMK_TempInverter = int16(md.AMK_TempInverter.UnmarshalSigned(f.Data))
	m.xxx_AMK_ErrorInfo = uint16(md.AMK_ErrorInfo.UnmarshalUnsigned(f.Data))
	m.xxx_AMK_TempIGBT = int16(md.AMK_TempIGBT.UnmarshalSigned(f.Data))
	return nil
}

// Nodes returns the pt node descriptors.
func Nodes() *NodesDescriptor {
	return nd
}

// NodesDescriptor contains all pt node descriptors.
type NodesDescriptor struct {
	AMK0 *descriptor.Node
	AMK1 *descriptor.Node
	FC   *descriptor.Node
}

// Messages returns the pt message descriptors.
func Messages() *MessagesDescriptor {
	return md
}

// MessagesDescriptor contains all pt message descriptors.
type MessagesDescriptor struct {
	AMK0_SetPoints1    *AMK0_SetPoints1Descriptor
	AMK1_SetPoints1    *AMK1_SetPoints1Descriptor
	AMK0_ActualValues1 *AMK0_ActualValues1Descriptor
	AMK1_ActualValues1 *AMK1_ActualValues1Descriptor
	AMK0_ActualValues2 *AMK0_ActualValues2Descriptor
	AMK1_ActualValues2 *AMK1_ActualValues2Descriptor
}

// UnmarshalFrame unmarshals the provided pt CAN frame.
func (md *MessagesDescriptor) UnmarshalFrame(f can.Frame) (generated.Message, error) {
	switch f.ID {
	case md.AMK0_SetPoints1.ID:
		var msg AMK0_SetPoints1
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal pt frame: %w", err)
		}
		return &msg, nil
	case md.AMK1_SetPoints1.ID:
		var msg AMK1_SetPoints1
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal pt frame: %w", err)
		}
		return &msg, nil
	case md.AMK0_ActualValues1.ID:
		var msg AMK0_ActualValues1
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal pt frame: %w", err)
		}
		return &msg, nil
	case md.AMK1_ActualValues1.ID:
		var msg AMK1_ActualValues1
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal pt frame: %w", err)
		}
		return &msg, nil
	case md.AMK0_ActualValues2.ID:
		var msg AMK0_ActualValues2
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal pt frame: %w", err)
		}
		return &msg, nil
	case md.AMK1_ActualValues2.ID:
		var msg AMK1_ActualValues2
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal pt frame: %w", err)
		}
		return &msg, nil
	default:
		return nil, fmt.Errorf("unmarshal pt frame: ID not in database: %d", f.ID)
	}
}

type AMK0_SetPoints1Descriptor struct {
	*descriptor.Message
	AMK_bInverterOn        *descriptor.Signal
	AMK_bDcOn              *descriptor.Signal
	AMK_bEnable            *descriptor.Signal
	AMK_bErrorReset        *descriptor.Signal
	AMK_TargetVelocity     *descriptor.Signal
	AMK_TorqueLimitPositiv *descriptor.Signal
	AMK_TorqueLimitNegativ *descriptor.Signal
}

type AMK1_SetPoints1Descriptor struct {
	*descriptor.Message
	AMK_bInverterOn        *descriptor.Signal
	AMK_bDcOn              *descriptor.Signal
	AMK_bEnable            *descriptor.Signal
	AMK_bErrorReset        *descriptor.Signal
	AMK_TargetVelocity     *descriptor.Signal
	AMK_TorqueLimitPositiv *descriptor.Signal
	AMK_TorqueLimitNegativ *descriptor.Signal
}

type AMK0_ActualValues1Descriptor struct {
	*descriptor.Message
	AMK_bSystemReady       *descriptor.Signal
	AMK_bError             *descriptor.Signal
	AMK_bWarn              *descriptor.Signal
	AMK_bQuitDcOn          *descriptor.Signal
	AMK_bDcOn              *descriptor.Signal
	AMK_bQuitInverterOn    *descriptor.Signal
	AMK_bInverterOn        *descriptor.Signal
	AMK_bDerating          *descriptor.Signal
	AMK_ActualVelocity     *descriptor.Signal
	AMK_TorqueCurrent      *descriptor.Signal
	AMK_MagnetizingCurrent *descriptor.Signal
}

type AMK1_ActualValues1Descriptor struct {
	*descriptor.Message
	AMK_bSystemReady       *descriptor.Signal
	AMK_bError             *descriptor.Signal
	AMK_bWarn              *descriptor.Signal
	AMK_bQuitDcOn          *descriptor.Signal
	AMK_bDcOn              *descriptor.Signal
	AMK_bQuitInverterOn    *descriptor.Signal
	AMK_bInverterOn        *descriptor.Signal
	AMK_bDerating          *descriptor.Signal
	AMK_ActualVelocity     *descriptor.Signal
	AMK_TorqueCurrent      *descriptor.Signal
	AMK_MagnetizingCurrent *descriptor.Signal
}

type AMK0_ActualValues2Descriptor struct {
	*descriptor.Message
	AMK_TempMotor    *descriptor.Signal
	AMK_TempInverter *descriptor.Signal
	AMK_ErrorInfo    *descriptor.Signal
	AMK_TempIGBT     *descriptor.Signal
}

type AMK1_ActualValues2Descriptor struct {
	*descriptor.Message
	AMK_TempMotor    *descriptor.Signal
	AMK_TempInverter *descriptor.Signal
	AMK_ErrorInfo    *descriptor.Signal
	AMK_TempIGBT     *descriptor.Signal
}

// Database returns the pt database descriptor.
func (md *MessagesDescriptor) Database() *descriptor.Database {
	return d
}

var nd = &NodesDescriptor{
	AMK0: d.Nodes[0],
	AMK1: d.Nodes[1],
	FC:   d.Nodes[2],
}

var md = &MessagesDescriptor{
	AMK0_SetPoints1: &AMK0_SetPoints1Descriptor{
		Message:                d.Messages[0],
		AMK_bInverterOn:        d.Messages[0].Signals[0],
		AMK_bDcOn:              d.Messages[0].Signals[1],
		AMK_bEnable:            d.Messages[0].Signals[2],
		AMK_bErrorReset:        d.Messages[0].Signals[3],
		AMK_TargetVelocity:     d.Messages[0].Signals[4],
		AMK_TorqueLimitPositiv: d.Messages[0].Signals[5],
		AMK_TorqueLimitNegativ: d.Messages[0].Signals[6],
	},
	AMK1_SetPoints1: &AMK1_SetPoints1Descriptor{
		Message:                d.Messages[1],
		AMK_bInverterOn:        d.Messages[1].Signals[0],
		AMK_bDcOn:              d.Messages[1].Signals[1],
		AMK_bEnable:            d.Messages[1].Signals[2],
		AMK_bErrorReset:        d.Messages[1].Signals[3],
		AMK_TargetVelocity:     d.Messages[1].Signals[4],
		AMK_TorqueLimitPositiv: d.Messages[1].Signals[5],
		AMK_TorqueLimitNegativ: d.Messages[1].Signals[6],
	},
	AMK0_ActualValues1: &AMK0_ActualValues1Descriptor{
		Message:                d.Messages[2],
		AMK_bSystemReady:       d.Messages[2].Signals[0],
		AMK_bError:             d.Messages[2].Signals[1],
		AMK_bWarn:              d.Messages[2].Signals[2],
		AMK_bQuitDcOn:          d.Messages[2].Signals[3],
		AMK_bDcOn:              d.Messages[2].Signals[4],
		AMK_bQuitInverterOn:    d.Messages[2].Signals[5],
		AMK_bInverterOn:        d.Messages[2].Signals[6],
		AMK_bDerating:          d.Messages[2].Signals[7],
		AMK_ActualVelocity:     d.Messages[2].Signals[8],
		AMK_TorqueCurrent:      d.Messages[2].Signals[9],
		AMK_MagnetizingCurrent: d.Messages[2].Signals[10],
	},
	AMK1_ActualValues1: &AMK1_ActualValues1Descriptor{
		Message:                d.Messages[3],
		AMK_bSystemReady:       d.Messages[3].Signals[0],
		AMK_bError:             d.Messages[3].Signals[1],
		AMK_bWarn:              d.Messages[3].Signals[2],
		AMK_bQuitDcOn:          d.Messages[3].Signals[3],
		AMK_bDcOn:              d.Messages[3].Signals[4],
		AMK_bQuitInverterOn:    d.Messages[3].Signals[5],
		AMK_bInverterOn:        d.Messages[3].Signals[6],
		AMK_bDerating:          d.Messages[3].Signals[7],
		AMK_ActualVelocity:     d.Messages[3].Signals[8],
		AMK_TorqueCurrent:      d.Messages[3].Signals[9],
		AMK_MagnetizingCurrent: d.Messages[3].Signals[10],
	},
	AMK0_ActualValues2: &AMK0_ActualValues2Descriptor{
		Message:          d.Messages[4],
		AMK_TempMotor:    d.Messages[4].Signals[0],
		AMK_TempInverter: d.Messages[4].Signals[1],
		AMK_ErrorInfo:    d.Messages[4].Signals[2],
		AMK_TempIGBT:     d.Messages[4].Signals[3],
	},
	AMK1_ActualValues2: &AMK1_ActualValues2Descriptor{
		Message:          d.Messages[5],
		AMK_TempMotor:    d.Messages[5].Signals[0],
		AMK_TempInverter: d.Messages[5].Signals[1],
		AMK_ErrorInfo:    d.Messages[5].Signals[2],
		AMK_TempIGBT:     d.Messages[5].Signals[3],
	},
}

var d = (*descriptor.Database)(&descriptor.Database{
	SourceFile: (string)("temp/ptcan/pt.dbc"),
	Version:    (string)(""),
	Messages: ([]*descriptor.Message)([]*descriptor.Message{
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("AMK0_SetPoints1"),
			ID:          (uint32)(389),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bInverterOn"),
					Start:             (uint8)(8),
					Length:            (uint8)(1),
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
						(string)("AMK0"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bDcOn"),
					Start:             (uint8)(9),
					Length:            (uint8)(1),
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
						(string)("AMK0"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bEnable"),
					Start:             (uint8)(10),
					Length:            (uint8)(1),
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
						(string)("AMK0"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bErrorReset"),
					Start:             (uint8)(11),
					Length:            (uint8)(1),
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
						(string)("AMK0"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_TargetVelocity"),
					Start:             (uint8)(16),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("rpm"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("AMK0"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_TorqueLimitPositiv"),
					Start:             (uint8)(32),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("0.1%Mn"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("AMK0"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_TorqueLimitNegativ"),
					Start:             (uint8)(48),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("0.1%Mn"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("AMK0"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("FC"),
			CycleTime:  (time.Duration)(5000000),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("AMK1_SetPoints1"),
			ID:          (uint32)(390),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bInverterOn"),
					Start:             (uint8)(8),
					Length:            (uint8)(1),
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
						(string)("AMK1"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bDcOn"),
					Start:             (uint8)(9),
					Length:            (uint8)(1),
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
						(string)("AMK1"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bEnable"),
					Start:             (uint8)(10),
					Length:            (uint8)(1),
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
						(string)("AMK1"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bErrorReset"),
					Start:             (uint8)(11),
					Length:            (uint8)(1),
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
						(string)("AMK1"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_TargetVelocity"),
					Start:             (uint8)(16),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("rpm"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("AMK1"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_TorqueLimitPositiv"),
					Start:             (uint8)(32),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("0.1%Mn"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("AMK1"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_TorqueLimitNegativ"),
					Start:             (uint8)(48),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("0.1%Mn"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("AMK1"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("FC"),
			CycleTime:  (time.Duration)(5000000),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("AMK0_ActualValues1"),
			ID:          (uint32)(644),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bSystemReady"),
					Start:             (uint8)(8),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bError"),
					Start:             (uint8)(9),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bWarn"),
					Start:             (uint8)(10),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bQuitDcOn"),
					Start:             (uint8)(11),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bDcOn"),
					Start:             (uint8)(12),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bQuitInverterOn"),
					Start:             (uint8)(13),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bInverterOn"),
					Start:             (uint8)(14),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bDerating"),
					Start:             (uint8)(15),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_ActualVelocity"),
					Start:             (uint8)(16),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("rpm"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_TorqueCurrent"),
					Start:             (uint8)(32),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_MagnetizingCurrent"),
					Start:             (uint8)(48),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("AMK0"),
			CycleTime:  (time.Duration)(5000000),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("AMK1_ActualValues1"),
			ID:          (uint32)(645),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bSystemReady"),
					Start:             (uint8)(8),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bError"),
					Start:             (uint8)(9),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bWarn"),
					Start:             (uint8)(10),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bQuitDcOn"),
					Start:             (uint8)(11),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bDcOn"),
					Start:             (uint8)(12),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bQuitInverterOn"),
					Start:             (uint8)(13),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bInverterOn"),
					Start:             (uint8)(14),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_bDerating"),
					Start:             (uint8)(15),
					Length:            (uint8)(1),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_ActualVelocity"),
					Start:             (uint8)(16),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("rpm"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_TorqueCurrent"),
					Start:             (uint8)(32),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_MagnetizingCurrent"),
					Start:             (uint8)(48),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
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
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("AMK1"),
			CycleTime:  (time.Duration)(5000000),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("AMK0_ActualValues2"),
			ID:          (uint32)(646),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_TempMotor"),
					Start:             (uint8)(0),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(-3276.8),
					Max:               (float64)(3276.7),
					Unit:              (string)("�C"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_TempInverter"),
					Start:             (uint8)(16),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(-3276.8),
					Max:               (float64)(3276.7),
					Unit:              (string)("�C"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_ErrorInfo"),
					Start:             (uint8)(32),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(65535),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_TempIGBT"),
					Start:             (uint8)(48),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(-3276.8),
					Max:               (float64)(3276.7),
					Unit:              (string)("�C"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("AMK0"),
			CycleTime:  (time.Duration)(5000000),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("AMK1_ActualValues2"),
			ID:          (uint32)(647),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_TempMotor"),
					Start:             (uint8)(0),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(-3276.8),
					Max:               (float64)(3276.7),
					Unit:              (string)("�C"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_TempInverter"),
					Start:             (uint8)(16),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(-3276.8),
					Max:               (float64)(3276.7),
					Unit:              (string)("�C"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_ErrorInfo"),
					Start:             (uint8)(32),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(65535),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AMK_TempIGBT"),
					Start:             (uint8)(48),
					Length:            (uint8)(16),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(-3276.8),
					Max:               (float64)(3276.7),
					Unit:              (string)("�C"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("AMK1"),
			CycleTime:  (time.Duration)(5000000),
			DelayTime:  (time.Duration)(0),
		}),
	}),
	Nodes: ([]*descriptor.Node)([]*descriptor.Node{
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("AMK0"),
			Description: (string)(""),
		}),
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("AMK1"),
			Description: (string)(""),
		}),
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("FC"),
			Description: (string)(""),
		}),
	}),
})
