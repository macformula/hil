// Package vehcan provides primitives for encoding and decoding veh CAN messages.
//
// Source: temp/vehcan/veh.dbc
package vehcan

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
// VC_StatusReader provides read access to a VC_Status message.
type VC_StatusReader interface {
	// VC_govStatus returns the value of the VC_govStatus signal.
	VC_govStatus() VC_Status_VC_govStatus
}

// VC_StatusWriter provides write access to a VC_Status message.
type VC_StatusWriter interface {
	// CopyFrom copies all values from VC_StatusReader.
	CopyFrom(VC_StatusReader) *VC_Status
	// SetVC_govStatus sets the value of the VC_govStatus signal.
	SetVC_govStatus(VC_Status_VC_govStatus) *VC_Status
}

type VC_Status struct {
	xxx_VC_govStatus VC_Status_VC_govStatus
}

func NewVC_Status() *VC_Status {
	m := &VC_Status{}
	m.Reset()
	return m
}

func (m *VC_Status) Reset() {
	m.xxx_VC_govStatus = 0
}

func (m *VC_Status) CopyFrom(o VC_StatusReader) *VC_Status {
	m.xxx_VC_govStatus = o.VC_govStatus()
	return m
}

// Descriptor returns the VC_Status descriptor.
func (m *VC_Status) Descriptor() *descriptor.Message {
	return Messages().VC_Status.Message
}

// String returns a compact string representation of the message.
func (m *VC_Status) String() string {
	return cantext.MessageString(m)
}

func (m *VC_Status) VC_govStatus() VC_Status_VC_govStatus {
	return m.xxx_VC_govStatus
}

func (m *VC_Status) SetVC_govStatus(v VC_Status_VC_govStatus) *VC_Status {
	m.xxx_VC_govStatus = VC_Status_VC_govStatus(Messages().VC_Status.VC_govStatus.SaturatedCastUnsigned(uint64(v)))
	return m
}

// VC_Status_VC_govStatus models the VC_govStatus signal of the VC_Status message.
type VC_Status_VC_govStatus uint8

// Value descriptions for the VC_govStatus signal of the VC_Status message.
const (
	VC_Status_VC_govStatus_govinit              VC_Status_VC_govStatus = 0
	VC_Status_VC_govStatus_govstartup           VC_Status_VC_govStatus = 1
	VC_Status_VC_govStatus_govrunning           VC_Status_VC_govStatus = 2
	VC_Status_VC_govStatus_hvstartuperror       VC_Status_VC_govStatus = 3
	VC_Status_VC_govStatus_motorstartuperror    VC_Status_VC_govStatus = 4
	VC_Status_VC_govStatus_driverinterfaceerror VC_Status_VC_govStatus = 5
	VC_Status_VC_govStatus_hvrunerror           VC_Status_VC_govStatus = 6
	VC_Status_VC_govStatus_motorrunerror        VC_Status_VC_govStatus = 7
)

func (v VC_Status_VC_govStatus) String() string {
	switch v {
	case 0:
		return "gov_init"
	case 1:
		return "gov_startup"
	case 2:
		return "gov_running"
	case 3:
		return "hv_startup_error"
	case 4:
		return "motor_startup_error"
	case 5:
		return "driver_interface_error"
	case 6:
		return "hv_run_error"
	case 7:
		return "motor_run_error"
	default:
		return fmt.Sprintf("VC_Status_VC_govStatus(%d)", v)
	}
}

// Frame returns a CAN frame representing the message.
func (m *VC_Status) Frame() can.Frame {
	md := Messages().VC_Status
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.VC_govStatus.MarshalUnsigned(&f.Data, uint64(m.xxx_VC_govStatus))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *VC_Status) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *VC_Status) UnmarshalFrame(f can.Frame) error {
	md := Messages().VC_Status
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal VC_Status: expects ID 255 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal VC_Status: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal VC_Status: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal VC_Status: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_VC_govStatus = VC_Status_VC_govStatus(md.VC_govStatus.UnmarshalUnsigned(f.Data))
	return nil
}

// FC_cmdReader provides read access to a FC_cmd message.
type FC_cmdReader interface {
	// FC_brakeLight returns the value of the FC_brakeLight signal.
	FC_brakeLight() bool
	// FC_readyToDriveSpeaker returns the value of the FC_readyToDriveSpeaker signal.
	FC_readyToDriveSpeaker() bool
}

// FC_cmdWriter provides write access to a FC_cmd message.
type FC_cmdWriter interface {
	// CopyFrom copies all values from FC_cmdReader.
	CopyFrom(FC_cmdReader) *FC_cmd
	// SetFC_brakeLight sets the value of the FC_brakeLight signal.
	SetFC_brakeLight(bool) *FC_cmd
	// SetFC_readyToDriveSpeaker sets the value of the FC_readyToDriveSpeaker signal.
	SetFC_readyToDriveSpeaker(bool) *FC_cmd
}

type FC_cmd struct {
	xxx_FC_brakeLight          bool
	xxx_FC_readyToDriveSpeaker bool
}

func NewFC_cmd() *FC_cmd {
	m := &FC_cmd{}
	m.Reset()
	return m
}

func (m *FC_cmd) Reset() {
	m.xxx_FC_brakeLight = false
	m.xxx_FC_readyToDriveSpeaker = false
}

func (m *FC_cmd) CopyFrom(o FC_cmdReader) *FC_cmd {
	m.xxx_FC_brakeLight = o.FC_brakeLight()
	m.xxx_FC_readyToDriveSpeaker = o.FC_readyToDriveSpeaker()
	return m
}

// Descriptor returns the FC_cmd descriptor.
func (m *FC_cmd) Descriptor() *descriptor.Message {
	return Messages().FC_cmd.Message
}

// String returns a compact string representation of the message.
func (m *FC_cmd) String() string {
	return cantext.MessageString(m)
}

func (m *FC_cmd) FC_brakeLight() bool {
	return m.xxx_FC_brakeLight
}

func (m *FC_cmd) SetFC_brakeLight(v bool) *FC_cmd {
	m.xxx_FC_brakeLight = v
	return m
}

func (m *FC_cmd) FC_readyToDriveSpeaker() bool {
	return m.xxx_FC_readyToDriveSpeaker
}

func (m *FC_cmd) SetFC_readyToDriveSpeaker(v bool) *FC_cmd {
	m.xxx_FC_readyToDriveSpeaker = v
	return m
}

// Frame returns a CAN frame representing the message.
func (m *FC_cmd) Frame() can.Frame {
	md := Messages().FC_cmd
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.FC_brakeLight.MarshalBool(&f.Data, bool(m.xxx_FC_brakeLight))
	md.FC_readyToDriveSpeaker.MarshalBool(&f.Data, bool(m.xxx_FC_readyToDriveSpeaker))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *FC_cmd) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *FC_cmd) UnmarshalFrame(f can.Frame) error {
	md := Messages().FC_cmd
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal FC_cmd: expects ID 256 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal FC_cmd: expects length 1 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal FC_cmd: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal FC_cmd: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_FC_brakeLight = bool(md.FC_brakeLight.UnmarshalBool(f.Data))
	m.xxx_FC_readyToDriveSpeaker = bool(md.FC_readyToDriveSpeaker.UnmarshalBool(f.Data))
	return nil
}

// FC_msgReader provides read access to a FC_msg message.
type FC_msgReader interface {
	// FC_apps1 returns the value of the FC_apps1 signal.
	FC_apps1() uint16
	// FC_apps2 returns the value of the FC_apps2 signal.
	FC_apps2() uint16
	// FC_bpps returns the value of the FC_bpps signal.
	FC_bpps() uint16
	// FC_steeringAngle returns the value of the FC_steeringAngle signal.
	FC_steeringAngle() uint16
	// FC_hvilSts returns the value of the FC_hvilSts signal.
	FC_hvilSts() bool
	// FC_readyToDriveBtn_n returns the value of the FC_readyToDriveBtn_n signal.
	FC_readyToDriveBtn_n() bool
}

// FC_msgWriter provides write access to a FC_msg message.
type FC_msgWriter interface {
	// CopyFrom copies all values from FC_msgReader.
	CopyFrom(FC_msgReader) *FC_msg
	// SetFC_apps1 sets the value of the FC_apps1 signal.
	SetFC_apps1(uint16) *FC_msg
	// SetFC_apps2 sets the value of the FC_apps2 signal.
	SetFC_apps2(uint16) *FC_msg
	// SetFC_bpps sets the value of the FC_bpps signal.
	SetFC_bpps(uint16) *FC_msg
	// SetFC_steeringAngle sets the value of the FC_steeringAngle signal.
	SetFC_steeringAngle(uint16) *FC_msg
	// SetFC_hvilSts sets the value of the FC_hvilSts signal.
	SetFC_hvilSts(bool) *FC_msg
	// SetFC_readyToDriveBtn_n sets the value of the FC_readyToDriveBtn_n signal.
	SetFC_readyToDriveBtn_n(bool) *FC_msg
}

type FC_msg struct {
	xxx_FC_apps1             uint16
	xxx_FC_apps2             uint16
	xxx_FC_bpps              uint16
	xxx_FC_steeringAngle     uint16
	xxx_FC_hvilSts           bool
	xxx_FC_readyToDriveBtn_n bool
}

func NewFC_msg() *FC_msg {
	m := &FC_msg{}
	m.Reset()
	return m
}

func (m *FC_msg) Reset() {
	m.xxx_FC_apps1 = 0
	m.xxx_FC_apps2 = 0
	m.xxx_FC_bpps = 0
	m.xxx_FC_steeringAngle = 0
	m.xxx_FC_hvilSts = false
	m.xxx_FC_readyToDriveBtn_n = false
}

func (m *FC_msg) CopyFrom(o FC_msgReader) *FC_msg {
	m.xxx_FC_apps1 = o.FC_apps1()
	m.xxx_FC_apps2 = o.FC_apps2()
	m.xxx_FC_bpps = o.FC_bpps()
	m.xxx_FC_steeringAngle = o.FC_steeringAngle()
	m.xxx_FC_hvilSts = o.FC_hvilSts()
	m.xxx_FC_readyToDriveBtn_n = o.FC_readyToDriveBtn_n()
	return m
}

// Descriptor returns the FC_msg descriptor.
func (m *FC_msg) Descriptor() *descriptor.Message {
	return Messages().FC_msg.Message
}

// String returns a compact string representation of the message.
func (m *FC_msg) String() string {
	return cantext.MessageString(m)
}

func (m *FC_msg) FC_apps1() uint16 {
	return m.xxx_FC_apps1
}

func (m *FC_msg) SetFC_apps1(v uint16) *FC_msg {
	m.xxx_FC_apps1 = uint16(Messages().FC_msg.FC_apps1.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *FC_msg) FC_apps2() uint16 {
	return m.xxx_FC_apps2
}

func (m *FC_msg) SetFC_apps2(v uint16) *FC_msg {
	m.xxx_FC_apps2 = uint16(Messages().FC_msg.FC_apps2.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *FC_msg) FC_bpps() uint16 {
	return m.xxx_FC_bpps
}

func (m *FC_msg) SetFC_bpps(v uint16) *FC_msg {
	m.xxx_FC_bpps = uint16(Messages().FC_msg.FC_bpps.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *FC_msg) FC_steeringAngle() uint16 {
	return m.xxx_FC_steeringAngle
}

func (m *FC_msg) SetFC_steeringAngle(v uint16) *FC_msg {
	m.xxx_FC_steeringAngle = uint16(Messages().FC_msg.FC_steeringAngle.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *FC_msg) FC_hvilSts() bool {
	return m.xxx_FC_hvilSts
}

func (m *FC_msg) SetFC_hvilSts(v bool) *FC_msg {
	m.xxx_FC_hvilSts = v
	return m
}

func (m *FC_msg) FC_readyToDriveBtn_n() bool {
	return m.xxx_FC_readyToDriveBtn_n
}

func (m *FC_msg) SetFC_readyToDriveBtn_n(v bool) *FC_msg {
	m.xxx_FC_readyToDriveBtn_n = v
	return m
}

// Frame returns a CAN frame representing the message.
func (m *FC_msg) Frame() can.Frame {
	md := Messages().FC_msg
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.FC_apps1.MarshalUnsigned(&f.Data, uint64(m.xxx_FC_apps1))
	md.FC_apps2.MarshalUnsigned(&f.Data, uint64(m.xxx_FC_apps2))
	md.FC_bpps.MarshalUnsigned(&f.Data, uint64(m.xxx_FC_bpps))
	md.FC_steeringAngle.MarshalUnsigned(&f.Data, uint64(m.xxx_FC_steeringAngle))
	md.FC_hvilSts.MarshalBool(&f.Data, bool(m.xxx_FC_hvilSts))
	md.FC_readyToDriveBtn_n.MarshalBool(&f.Data, bool(m.xxx_FC_readyToDriveBtn_n))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *FC_msg) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *FC_msg) UnmarshalFrame(f can.Frame) error {
	md := Messages().FC_msg
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal FC_msg: expects ID 511 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal FC_msg: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal FC_msg: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal FC_msg: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_FC_apps1 = uint16(md.FC_apps1.UnmarshalUnsigned(f.Data))
	m.xxx_FC_apps2 = uint16(md.FC_apps2.UnmarshalUnsigned(f.Data))
	m.xxx_FC_bpps = uint16(md.FC_bpps.UnmarshalUnsigned(f.Data))
	m.xxx_FC_steeringAngle = uint16(md.FC_steeringAngle.UnmarshalUnsigned(f.Data))
	m.xxx_FC_hvilSts = bool(md.FC_hvilSts.UnmarshalBool(f.Data))
	m.xxx_FC_readyToDriveBtn_n = bool(md.FC_readyToDriveBtn_n.UnmarshalBool(f.Data))
	return nil
}

// GnssStatusReader provides read access to a GnssStatus message.
type GnssStatusReader interface {
	// FixType returns the physical value of the FixType signal.
	FixType() float64
	// Satellites returns the value of the Satellites signal.
	Satellites() uint8
}

// GnssStatusWriter provides write access to a GnssStatus message.
type GnssStatusWriter interface {
	// CopyFrom copies all values from GnssStatusReader.
	CopyFrom(GnssStatusReader) *GnssStatus
	// SetFixType sets the physical value of the FixType signal.
	SetFixType(float64) *GnssStatus
	// SetSatellites sets the value of the Satellites signal.
	SetSatellites(uint8) *GnssStatus
}

type GnssStatus struct {
	xxx_FixType    uint8
	xxx_Satellites uint8
}

func NewGnssStatus() *GnssStatus {
	m := &GnssStatus{}
	m.Reset()
	return m
}

func (m *GnssStatus) Reset() {
	m.xxx_FixType = 0
	m.xxx_Satellites = 0
}

func (m *GnssStatus) CopyFrom(o GnssStatusReader) *GnssStatus {
	m.SetFixType(o.FixType())
	m.xxx_Satellites = o.Satellites()
	return m
}

// Descriptor returns the GnssStatus descriptor.
func (m *GnssStatus) Descriptor() *descriptor.Message {
	return Messages().GnssStatus.Message
}

// String returns a compact string representation of the message.
func (m *GnssStatus) String() string {
	return cantext.MessageString(m)
}

func (m *GnssStatus) FixType() float64 {
	return Messages().GnssStatus.FixType.ToPhysical(float64(m.xxx_FixType))
}

func (m *GnssStatus) SetFixType(v float64) *GnssStatus {
	m.xxx_FixType = uint8(Messages().GnssStatus.FixType.FromPhysical(v))
	return m
}

func (m *GnssStatus) Satellites() uint8 {
	return m.xxx_Satellites
}

func (m *GnssStatus) SetSatellites(v uint8) *GnssStatus {
	m.xxx_Satellites = uint8(Messages().GnssStatus.Satellites.SaturatedCastUnsigned(uint64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *GnssStatus) Frame() can.Frame {
	md := Messages().GnssStatus
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.FixType.MarshalUnsigned(&f.Data, uint64(m.xxx_FixType))
	md.Satellites.MarshalUnsigned(&f.Data, uint64(m.xxx_Satellites))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *GnssStatus) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *GnssStatus) UnmarshalFrame(f can.Frame) error {
	md := Messages().GnssStatus
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal GnssStatus: expects ID 769 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal GnssStatus: expects length 1 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal GnssStatus: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal GnssStatus: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_FixType = uint8(md.FixType.UnmarshalUnsigned(f.Data))
	m.xxx_Satellites = uint8(md.Satellites.UnmarshalUnsigned(f.Data))
	return nil
}

// GnssTimeReader provides read access to a GnssTime message.
type GnssTimeReader interface {
	// TimeValid returns the value of the TimeValid signal.
	TimeValid() bool
	// TimeConfirmed returns the value of the TimeConfirmed signal.
	TimeConfirmed() bool
	// Epoch returns the physical value of the Epoch signal.
	Epoch() float64
}

// GnssTimeWriter provides write access to a GnssTime message.
type GnssTimeWriter interface {
	// CopyFrom copies all values from GnssTimeReader.
	CopyFrom(GnssTimeReader) *GnssTime
	// SetTimeValid sets the value of the TimeValid signal.
	SetTimeValid(bool) *GnssTime
	// SetTimeConfirmed sets the value of the TimeConfirmed signal.
	SetTimeConfirmed(bool) *GnssTime
	// SetEpoch sets the physical value of the Epoch signal.
	SetEpoch(float64) *GnssTime
}

type GnssTime struct {
	xxx_TimeValid     bool
	xxx_TimeConfirmed bool
	xxx_Epoch         uint64
}

func NewGnssTime() *GnssTime {
	m := &GnssTime{}
	m.Reset()
	return m
}

func (m *GnssTime) Reset() {
	m.xxx_TimeValid = false
	m.xxx_TimeConfirmed = false
	m.xxx_Epoch = 0
}

func (m *GnssTime) CopyFrom(o GnssTimeReader) *GnssTime {
	m.xxx_TimeValid = o.TimeValid()
	m.xxx_TimeConfirmed = o.TimeConfirmed()
	m.SetEpoch(o.Epoch())
	return m
}

// Descriptor returns the GnssTime descriptor.
func (m *GnssTime) Descriptor() *descriptor.Message {
	return Messages().GnssTime.Message
}

// String returns a compact string representation of the message.
func (m *GnssTime) String() string {
	return cantext.MessageString(m)
}

func (m *GnssTime) TimeValid() bool {
	return m.xxx_TimeValid
}

func (m *GnssTime) SetTimeValid(v bool) *GnssTime {
	m.xxx_TimeValid = v
	return m
}

func (m *GnssTime) TimeConfirmed() bool {
	return m.xxx_TimeConfirmed
}

func (m *GnssTime) SetTimeConfirmed(v bool) *GnssTime {
	m.xxx_TimeConfirmed = v
	return m
}

func (m *GnssTime) Epoch() float64 {
	return Messages().GnssTime.Epoch.ToPhysical(float64(m.xxx_Epoch))
}

func (m *GnssTime) SetEpoch(v float64) *GnssTime {
	m.xxx_Epoch = uint64(Messages().GnssTime.Epoch.FromPhysical(v))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *GnssTime) Frame() can.Frame {
	md := Messages().GnssTime
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.TimeValid.MarshalBool(&f.Data, bool(m.xxx_TimeValid))
	md.TimeConfirmed.MarshalBool(&f.Data, bool(m.xxx_TimeConfirmed))
	md.Epoch.MarshalUnsigned(&f.Data, uint64(m.xxx_Epoch))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *GnssTime) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *GnssTime) UnmarshalFrame(f can.Frame) error {
	md := Messages().GnssTime
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal GnssTime: expects ID 770 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal GnssTime: expects length 6 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal GnssTime: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal GnssTime: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_TimeValid = bool(md.TimeValid.UnmarshalBool(f.Data))
	m.xxx_TimeConfirmed = bool(md.TimeConfirmed.UnmarshalBool(f.Data))
	m.xxx_Epoch = uint64(md.Epoch.UnmarshalUnsigned(f.Data))
	return nil
}

// GnssPositionReader provides read access to a GnssPosition message.
type GnssPositionReader interface {
	// PositionValid returns the value of the PositionValid signal.
	PositionValid() bool
	// Latitude returns the physical value of the Latitude signal.
	Latitude() float64
	// Longitude returns the physical value of the Longitude signal.
	Longitude() float64
	// PositionAccuracy returns the value of the PositionAccuracy signal.
	PositionAccuracy() uint8
}

// GnssPositionWriter provides write access to a GnssPosition message.
type GnssPositionWriter interface {
	// CopyFrom copies all values from GnssPositionReader.
	CopyFrom(GnssPositionReader) *GnssPosition
	// SetPositionValid sets the value of the PositionValid signal.
	SetPositionValid(bool) *GnssPosition
	// SetLatitude sets the physical value of the Latitude signal.
	SetLatitude(float64) *GnssPosition
	// SetLongitude sets the physical value of the Longitude signal.
	SetLongitude(float64) *GnssPosition
	// SetPositionAccuracy sets the value of the PositionAccuracy signal.
	SetPositionAccuracy(uint8) *GnssPosition
}

type GnssPosition struct {
	xxx_PositionValid    bool
	xxx_Latitude         uint32
	xxx_Longitude        uint32
	xxx_PositionAccuracy uint8
}

func NewGnssPosition() *GnssPosition {
	m := &GnssPosition{}
	m.Reset()
	return m
}

func (m *GnssPosition) Reset() {
	m.xxx_PositionValid = false
	m.xxx_Latitude = 0
	m.xxx_Longitude = 0
	m.xxx_PositionAccuracy = 0
}

func (m *GnssPosition) CopyFrom(o GnssPositionReader) *GnssPosition {
	m.xxx_PositionValid = o.PositionValid()
	m.SetLatitude(o.Latitude())
	m.SetLongitude(o.Longitude())
	m.xxx_PositionAccuracy = o.PositionAccuracy()
	return m
}

// Descriptor returns the GnssPosition descriptor.
func (m *GnssPosition) Descriptor() *descriptor.Message {
	return Messages().GnssPosition.Message
}

// String returns a compact string representation of the message.
func (m *GnssPosition) String() string {
	return cantext.MessageString(m)
}

func (m *GnssPosition) PositionValid() bool {
	return m.xxx_PositionValid
}

func (m *GnssPosition) SetPositionValid(v bool) *GnssPosition {
	m.xxx_PositionValid = v
	return m
}

func (m *GnssPosition) Latitude() float64 {
	return Messages().GnssPosition.Latitude.ToPhysical(float64(m.xxx_Latitude))
}

func (m *GnssPosition) SetLatitude(v float64) *GnssPosition {
	m.xxx_Latitude = uint32(Messages().GnssPosition.Latitude.FromPhysical(v))
	return m
}

func (m *GnssPosition) Longitude() float64 {
	return Messages().GnssPosition.Longitude.ToPhysical(float64(m.xxx_Longitude))
}

func (m *GnssPosition) SetLongitude(v float64) *GnssPosition {
	m.xxx_Longitude = uint32(Messages().GnssPosition.Longitude.FromPhysical(v))
	return m
}

func (m *GnssPosition) PositionAccuracy() uint8 {
	return m.xxx_PositionAccuracy
}

func (m *GnssPosition) SetPositionAccuracy(v uint8) *GnssPosition {
	m.xxx_PositionAccuracy = uint8(Messages().GnssPosition.PositionAccuracy.SaturatedCastUnsigned(uint64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *GnssPosition) Frame() can.Frame {
	md := Messages().GnssPosition
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.PositionValid.MarshalBool(&f.Data, bool(m.xxx_PositionValid))
	md.Latitude.MarshalUnsigned(&f.Data, uint64(m.xxx_Latitude))
	md.Longitude.MarshalUnsigned(&f.Data, uint64(m.xxx_Longitude))
	md.PositionAccuracy.MarshalUnsigned(&f.Data, uint64(m.xxx_PositionAccuracy))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *GnssPosition) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *GnssPosition) UnmarshalFrame(f can.Frame) error {
	md := Messages().GnssPosition
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal GnssPosition: expects ID 771 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal GnssPosition: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal GnssPosition: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal GnssPosition: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_PositionValid = bool(md.PositionValid.UnmarshalBool(f.Data))
	m.xxx_Latitude = uint32(md.Latitude.UnmarshalUnsigned(f.Data))
	m.xxx_Longitude = uint32(md.Longitude.UnmarshalUnsigned(f.Data))
	m.xxx_PositionAccuracy = uint8(md.PositionAccuracy.UnmarshalUnsigned(f.Data))
	return nil
}

// GnssAltitudeReader provides read access to a GnssAltitude message.
type GnssAltitudeReader interface {
	// AltitudeValid returns the value of the AltitudeValid signal.
	AltitudeValid() bool
	// Altitude returns the physical value of the Altitude signal.
	Altitude() float64
	// AltitudeAccuracy returns the physical value of the AltitudeAccuracy signal.
	AltitudeAccuracy() float64
}

// GnssAltitudeWriter provides write access to a GnssAltitude message.
type GnssAltitudeWriter interface {
	// CopyFrom copies all values from GnssAltitudeReader.
	CopyFrom(GnssAltitudeReader) *GnssAltitude
	// SetAltitudeValid sets the value of the AltitudeValid signal.
	SetAltitudeValid(bool) *GnssAltitude
	// SetAltitude sets the physical value of the Altitude signal.
	SetAltitude(float64) *GnssAltitude
	// SetAltitudeAccuracy sets the physical value of the AltitudeAccuracy signal.
	SetAltitudeAccuracy(float64) *GnssAltitude
}

type GnssAltitude struct {
	xxx_AltitudeValid    bool
	xxx_Altitude         uint32
	xxx_AltitudeAccuracy uint16
}

func NewGnssAltitude() *GnssAltitude {
	m := &GnssAltitude{}
	m.Reset()
	return m
}

func (m *GnssAltitude) Reset() {
	m.xxx_AltitudeValid = false
	m.xxx_Altitude = 0
	m.xxx_AltitudeAccuracy = 0
}

func (m *GnssAltitude) CopyFrom(o GnssAltitudeReader) *GnssAltitude {
	m.xxx_AltitudeValid = o.AltitudeValid()
	m.SetAltitude(o.Altitude())
	m.SetAltitudeAccuracy(o.AltitudeAccuracy())
	return m
}

// Descriptor returns the GnssAltitude descriptor.
func (m *GnssAltitude) Descriptor() *descriptor.Message {
	return Messages().GnssAltitude.Message
}

// String returns a compact string representation of the message.
func (m *GnssAltitude) String() string {
	return cantext.MessageString(m)
}

func (m *GnssAltitude) AltitudeValid() bool {
	return m.xxx_AltitudeValid
}

func (m *GnssAltitude) SetAltitudeValid(v bool) *GnssAltitude {
	m.xxx_AltitudeValid = v
	return m
}

func (m *GnssAltitude) Altitude() float64 {
	return Messages().GnssAltitude.Altitude.ToPhysical(float64(m.xxx_Altitude))
}

func (m *GnssAltitude) SetAltitude(v float64) *GnssAltitude {
	m.xxx_Altitude = uint32(Messages().GnssAltitude.Altitude.FromPhysical(v))
	return m
}

func (m *GnssAltitude) AltitudeAccuracy() float64 {
	return Messages().GnssAltitude.AltitudeAccuracy.ToPhysical(float64(m.xxx_AltitudeAccuracy))
}

func (m *GnssAltitude) SetAltitudeAccuracy(v float64) *GnssAltitude {
	m.xxx_AltitudeAccuracy = uint16(Messages().GnssAltitude.AltitudeAccuracy.FromPhysical(v))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *GnssAltitude) Frame() can.Frame {
	md := Messages().GnssAltitude
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.AltitudeValid.MarshalBool(&f.Data, bool(m.xxx_AltitudeValid))
	md.Altitude.MarshalUnsigned(&f.Data, uint64(m.xxx_Altitude))
	md.AltitudeAccuracy.MarshalUnsigned(&f.Data, uint64(m.xxx_AltitudeAccuracy))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *GnssAltitude) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *GnssAltitude) UnmarshalFrame(f can.Frame) error {
	md := Messages().GnssAltitude
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal GnssAltitude: expects ID 772 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal GnssAltitude: expects length 4 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal GnssAltitude: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal GnssAltitude: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_AltitudeValid = bool(md.AltitudeValid.UnmarshalBool(f.Data))
	m.xxx_Altitude = uint32(md.Altitude.UnmarshalUnsigned(f.Data))
	m.xxx_AltitudeAccuracy = uint16(md.AltitudeAccuracy.UnmarshalUnsigned(f.Data))
	return nil
}

// GnssAttitudeReader provides read access to a GnssAttitude message.
type GnssAttitudeReader interface {
	// AttitudeValid returns the value of the AttitudeValid signal.
	AttitudeValid() bool
	// Roll returns the physical value of the Roll signal.
	Roll() float64
	// RollAccuracy returns the physical value of the RollAccuracy signal.
	RollAccuracy() float64
	// Pitch returns the physical value of the Pitch signal.
	Pitch() float64
	// PitchAccuracy returns the physical value of the PitchAccuracy signal.
	PitchAccuracy() float64
	// Heading returns the physical value of the Heading signal.
	Heading() float64
	// HeadingAccuracy returns the physical value of the HeadingAccuracy signal.
	HeadingAccuracy() float64
}

// GnssAttitudeWriter provides write access to a GnssAttitude message.
type GnssAttitudeWriter interface {
	// CopyFrom copies all values from GnssAttitudeReader.
	CopyFrom(GnssAttitudeReader) *GnssAttitude
	// SetAttitudeValid sets the value of the AttitudeValid signal.
	SetAttitudeValid(bool) *GnssAttitude
	// SetRoll sets the physical value of the Roll signal.
	SetRoll(float64) *GnssAttitude
	// SetRollAccuracy sets the physical value of the RollAccuracy signal.
	SetRollAccuracy(float64) *GnssAttitude
	// SetPitch sets the physical value of the Pitch signal.
	SetPitch(float64) *GnssAttitude
	// SetPitchAccuracy sets the physical value of the PitchAccuracy signal.
	SetPitchAccuracy(float64) *GnssAttitude
	// SetHeading sets the physical value of the Heading signal.
	SetHeading(float64) *GnssAttitude
	// SetHeadingAccuracy sets the physical value of the HeadingAccuracy signal.
	SetHeadingAccuracy(float64) *GnssAttitude
}

type GnssAttitude struct {
	xxx_AttitudeValid   bool
	xxx_Roll            uint16
	xxx_RollAccuracy    uint16
	xxx_Pitch           uint16
	xxx_PitchAccuracy   uint16
	xxx_Heading         uint16
	xxx_HeadingAccuracy uint16
}

func NewGnssAttitude() *GnssAttitude {
	m := &GnssAttitude{}
	m.Reset()
	return m
}

func (m *GnssAttitude) Reset() {
	m.xxx_AttitudeValid = false
	m.xxx_Roll = 0
	m.xxx_RollAccuracy = 0
	m.xxx_Pitch = 0
	m.xxx_PitchAccuracy = 0
	m.xxx_Heading = 0
	m.xxx_HeadingAccuracy = 0
}

func (m *GnssAttitude) CopyFrom(o GnssAttitudeReader) *GnssAttitude {
	m.xxx_AttitudeValid = o.AttitudeValid()
	m.SetRoll(o.Roll())
	m.SetRollAccuracy(o.RollAccuracy())
	m.SetPitch(o.Pitch())
	m.SetPitchAccuracy(o.PitchAccuracy())
	m.SetHeading(o.Heading())
	m.SetHeadingAccuracy(o.HeadingAccuracy())
	return m
}

// Descriptor returns the GnssAttitude descriptor.
func (m *GnssAttitude) Descriptor() *descriptor.Message {
	return Messages().GnssAttitude.Message
}

// String returns a compact string representation of the message.
func (m *GnssAttitude) String() string {
	return cantext.MessageString(m)
}

func (m *GnssAttitude) AttitudeValid() bool {
	return m.xxx_AttitudeValid
}

func (m *GnssAttitude) SetAttitudeValid(v bool) *GnssAttitude {
	m.xxx_AttitudeValid = v
	return m
}

func (m *GnssAttitude) Roll() float64 {
	return Messages().GnssAttitude.Roll.ToPhysical(float64(m.xxx_Roll))
}

func (m *GnssAttitude) SetRoll(v float64) *GnssAttitude {
	m.xxx_Roll = uint16(Messages().GnssAttitude.Roll.FromPhysical(v))
	return m
}

func (m *GnssAttitude) RollAccuracy() float64 {
	return Messages().GnssAttitude.RollAccuracy.ToPhysical(float64(m.xxx_RollAccuracy))
}

func (m *GnssAttitude) SetRollAccuracy(v float64) *GnssAttitude {
	m.xxx_RollAccuracy = uint16(Messages().GnssAttitude.RollAccuracy.FromPhysical(v))
	return m
}

func (m *GnssAttitude) Pitch() float64 {
	return Messages().GnssAttitude.Pitch.ToPhysical(float64(m.xxx_Pitch))
}

func (m *GnssAttitude) SetPitch(v float64) *GnssAttitude {
	m.xxx_Pitch = uint16(Messages().GnssAttitude.Pitch.FromPhysical(v))
	return m
}

func (m *GnssAttitude) PitchAccuracy() float64 {
	return Messages().GnssAttitude.PitchAccuracy.ToPhysical(float64(m.xxx_PitchAccuracy))
}

func (m *GnssAttitude) SetPitchAccuracy(v float64) *GnssAttitude {
	m.xxx_PitchAccuracy = uint16(Messages().GnssAttitude.PitchAccuracy.FromPhysical(v))
	return m
}

func (m *GnssAttitude) Heading() float64 {
	return Messages().GnssAttitude.Heading.ToPhysical(float64(m.xxx_Heading))
}

func (m *GnssAttitude) SetHeading(v float64) *GnssAttitude {
	m.xxx_Heading = uint16(Messages().GnssAttitude.Heading.FromPhysical(v))
	return m
}

func (m *GnssAttitude) HeadingAccuracy() float64 {
	return Messages().GnssAttitude.HeadingAccuracy.ToPhysical(float64(m.xxx_HeadingAccuracy))
}

func (m *GnssAttitude) SetHeadingAccuracy(v float64) *GnssAttitude {
	m.xxx_HeadingAccuracy = uint16(Messages().GnssAttitude.HeadingAccuracy.FromPhysical(v))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *GnssAttitude) Frame() can.Frame {
	md := Messages().GnssAttitude
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.AttitudeValid.MarshalBool(&f.Data, bool(m.xxx_AttitudeValid))
	md.Roll.MarshalUnsigned(&f.Data, uint64(m.xxx_Roll))
	md.RollAccuracy.MarshalUnsigned(&f.Data, uint64(m.xxx_RollAccuracy))
	md.Pitch.MarshalUnsigned(&f.Data, uint64(m.xxx_Pitch))
	md.PitchAccuracy.MarshalUnsigned(&f.Data, uint64(m.xxx_PitchAccuracy))
	md.Heading.MarshalUnsigned(&f.Data, uint64(m.xxx_Heading))
	md.HeadingAccuracy.MarshalUnsigned(&f.Data, uint64(m.xxx_HeadingAccuracy))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *GnssAttitude) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *GnssAttitude) UnmarshalFrame(f can.Frame) error {
	md := Messages().GnssAttitude
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal GnssAttitude: expects ID 773 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal GnssAttitude: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal GnssAttitude: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal GnssAttitude: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_AttitudeValid = bool(md.AttitudeValid.UnmarshalBool(f.Data))
	m.xxx_Roll = uint16(md.Roll.UnmarshalUnsigned(f.Data))
	m.xxx_RollAccuracy = uint16(md.RollAccuracy.UnmarshalUnsigned(f.Data))
	m.xxx_Pitch = uint16(md.Pitch.UnmarshalUnsigned(f.Data))
	m.xxx_PitchAccuracy = uint16(md.PitchAccuracy.UnmarshalUnsigned(f.Data))
	m.xxx_Heading = uint16(md.Heading.UnmarshalUnsigned(f.Data))
	m.xxx_HeadingAccuracy = uint16(md.HeadingAccuracy.UnmarshalUnsigned(f.Data))
	return nil
}

// GnssOdoReader provides read access to a GnssOdo message.
type GnssOdoReader interface {
	// DistanceValid returns the value of the DistanceValid signal.
	DistanceValid() bool
	// DistanceTrip returns the value of the DistanceTrip signal.
	DistanceTrip() uint32
	// DistanceAccuracy returns the value of the DistanceAccuracy signal.
	DistanceAccuracy() uint32
	// DistanceTotal returns the value of the DistanceTotal signal.
	DistanceTotal() uint32
}

// GnssOdoWriter provides write access to a GnssOdo message.
type GnssOdoWriter interface {
	// CopyFrom copies all values from GnssOdoReader.
	CopyFrom(GnssOdoReader) *GnssOdo
	// SetDistanceValid sets the value of the DistanceValid signal.
	SetDistanceValid(bool) *GnssOdo
	// SetDistanceTrip sets the value of the DistanceTrip signal.
	SetDistanceTrip(uint32) *GnssOdo
	// SetDistanceAccuracy sets the value of the DistanceAccuracy signal.
	SetDistanceAccuracy(uint32) *GnssOdo
	// SetDistanceTotal sets the value of the DistanceTotal signal.
	SetDistanceTotal(uint32) *GnssOdo
}

type GnssOdo struct {
	xxx_DistanceValid    bool
	xxx_DistanceTrip     uint32
	xxx_DistanceAccuracy uint32
	xxx_DistanceTotal    uint32
}

func NewGnssOdo() *GnssOdo {
	m := &GnssOdo{}
	m.Reset()
	return m
}

func (m *GnssOdo) Reset() {
	m.xxx_DistanceValid = false
	m.xxx_DistanceTrip = 0
	m.xxx_DistanceAccuracy = 0
	m.xxx_DistanceTotal = 0
}

func (m *GnssOdo) CopyFrom(o GnssOdoReader) *GnssOdo {
	m.xxx_DistanceValid = o.DistanceValid()
	m.xxx_DistanceTrip = o.DistanceTrip()
	m.xxx_DistanceAccuracy = o.DistanceAccuracy()
	m.xxx_DistanceTotal = o.DistanceTotal()
	return m
}

// Descriptor returns the GnssOdo descriptor.
func (m *GnssOdo) Descriptor() *descriptor.Message {
	return Messages().GnssOdo.Message
}

// String returns a compact string representation of the message.
func (m *GnssOdo) String() string {
	return cantext.MessageString(m)
}

func (m *GnssOdo) DistanceValid() bool {
	return m.xxx_DistanceValid
}

func (m *GnssOdo) SetDistanceValid(v bool) *GnssOdo {
	m.xxx_DistanceValid = v
	return m
}

func (m *GnssOdo) DistanceTrip() uint32 {
	return m.xxx_DistanceTrip
}

func (m *GnssOdo) SetDistanceTrip(v uint32) *GnssOdo {
	m.xxx_DistanceTrip = uint32(Messages().GnssOdo.DistanceTrip.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *GnssOdo) DistanceAccuracy() uint32 {
	return m.xxx_DistanceAccuracy
}

func (m *GnssOdo) SetDistanceAccuracy(v uint32) *GnssOdo {
	m.xxx_DistanceAccuracy = uint32(Messages().GnssOdo.DistanceAccuracy.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *GnssOdo) DistanceTotal() uint32 {
	return m.xxx_DistanceTotal
}

func (m *GnssOdo) SetDistanceTotal(v uint32) *GnssOdo {
	m.xxx_DistanceTotal = uint32(Messages().GnssOdo.DistanceTotal.SaturatedCastUnsigned(uint64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *GnssOdo) Frame() can.Frame {
	md := Messages().GnssOdo
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.DistanceValid.MarshalBool(&f.Data, bool(m.xxx_DistanceValid))
	md.DistanceTrip.MarshalUnsigned(&f.Data, uint64(m.xxx_DistanceTrip))
	md.DistanceAccuracy.MarshalUnsigned(&f.Data, uint64(m.xxx_DistanceAccuracy))
	md.DistanceTotal.MarshalUnsigned(&f.Data, uint64(m.xxx_DistanceTotal))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *GnssOdo) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *GnssOdo) UnmarshalFrame(f can.Frame) error {
	md := Messages().GnssOdo
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal GnssOdo: expects ID 774 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal GnssOdo: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal GnssOdo: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal GnssOdo: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_DistanceValid = bool(md.DistanceValid.UnmarshalBool(f.Data))
	m.xxx_DistanceTrip = uint32(md.DistanceTrip.UnmarshalUnsigned(f.Data))
	m.xxx_DistanceAccuracy = uint32(md.DistanceAccuracy.UnmarshalUnsigned(f.Data))
	m.xxx_DistanceTotal = uint32(md.DistanceTotal.UnmarshalUnsigned(f.Data))
	return nil
}

// GnssSpeedReader provides read access to a GnssSpeed message.
type GnssSpeedReader interface {
	// SpeedValid returns the value of the SpeedValid signal.
	SpeedValid() bool
	// Speed returns the physical value of the Speed signal.
	Speed() float64
	// SpeedAccuracy returns the physical value of the SpeedAccuracy signal.
	SpeedAccuracy() float64
}

// GnssSpeedWriter provides write access to a GnssSpeed message.
type GnssSpeedWriter interface {
	// CopyFrom copies all values from GnssSpeedReader.
	CopyFrom(GnssSpeedReader) *GnssSpeed
	// SetSpeedValid sets the value of the SpeedValid signal.
	SetSpeedValid(bool) *GnssSpeed
	// SetSpeed sets the physical value of the Speed signal.
	SetSpeed(float64) *GnssSpeed
	// SetSpeedAccuracy sets the physical value of the SpeedAccuracy signal.
	SetSpeedAccuracy(float64) *GnssSpeed
}

type GnssSpeed struct {
	xxx_SpeedValid    bool
	xxx_Speed         uint32
	xxx_SpeedAccuracy uint32
}

func NewGnssSpeed() *GnssSpeed {
	m := &GnssSpeed{}
	m.Reset()
	return m
}

func (m *GnssSpeed) Reset() {
	m.xxx_SpeedValid = false
	m.xxx_Speed = 0
	m.xxx_SpeedAccuracy = 0
}

func (m *GnssSpeed) CopyFrom(o GnssSpeedReader) *GnssSpeed {
	m.xxx_SpeedValid = o.SpeedValid()
	m.SetSpeed(o.Speed())
	m.SetSpeedAccuracy(o.SpeedAccuracy())
	return m
}

// Descriptor returns the GnssSpeed descriptor.
func (m *GnssSpeed) Descriptor() *descriptor.Message {
	return Messages().GnssSpeed.Message
}

// String returns a compact string representation of the message.
func (m *GnssSpeed) String() string {
	return cantext.MessageString(m)
}

func (m *GnssSpeed) SpeedValid() bool {
	return m.xxx_SpeedValid
}

func (m *GnssSpeed) SetSpeedValid(v bool) *GnssSpeed {
	m.xxx_SpeedValid = v
	return m
}

func (m *GnssSpeed) Speed() float64 {
	return Messages().GnssSpeed.Speed.ToPhysical(float64(m.xxx_Speed))
}

func (m *GnssSpeed) SetSpeed(v float64) *GnssSpeed {
	m.xxx_Speed = uint32(Messages().GnssSpeed.Speed.FromPhysical(v))
	return m
}

func (m *GnssSpeed) SpeedAccuracy() float64 {
	return Messages().GnssSpeed.SpeedAccuracy.ToPhysical(float64(m.xxx_SpeedAccuracy))
}

func (m *GnssSpeed) SetSpeedAccuracy(v float64) *GnssSpeed {
	m.xxx_SpeedAccuracy = uint32(Messages().GnssSpeed.SpeedAccuracy.FromPhysical(v))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *GnssSpeed) Frame() can.Frame {
	md := Messages().GnssSpeed
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.SpeedValid.MarshalBool(&f.Data, bool(m.xxx_SpeedValid))
	md.Speed.MarshalUnsigned(&f.Data, uint64(m.xxx_Speed))
	md.SpeedAccuracy.MarshalUnsigned(&f.Data, uint64(m.xxx_SpeedAccuracy))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *GnssSpeed) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *GnssSpeed) UnmarshalFrame(f can.Frame) error {
	md := Messages().GnssSpeed
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal GnssSpeed: expects ID 775 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal GnssSpeed: expects length 5 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal GnssSpeed: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal GnssSpeed: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_SpeedValid = bool(md.SpeedValid.UnmarshalBool(f.Data))
	m.xxx_Speed = uint32(md.Speed.UnmarshalUnsigned(f.Data))
	m.xxx_SpeedAccuracy = uint32(md.SpeedAccuracy.UnmarshalUnsigned(f.Data))
	return nil
}

// GnssGeofenceReader provides read access to a GnssGeofence message.
type GnssGeofenceReader interface {
	// FenceValid returns the value of the FenceValid signal.
	FenceValid() bool
	// FenceCombined returns the physical value of the FenceCombined signal.
	FenceCombined() float64
	// Fence1 returns the physical value of the Fence1 signal.
	Fence1() float64
	// Fence2 returns the physical value of the Fence2 signal.
	Fence2() float64
	// Fence3 returns the physical value of the Fence3 signal.
	Fence3() float64
	// Fence4 returns the physical value of the Fence4 signal.
	Fence4() float64
}

// GnssGeofenceWriter provides write access to a GnssGeofence message.
type GnssGeofenceWriter interface {
	// CopyFrom copies all values from GnssGeofenceReader.
	CopyFrom(GnssGeofenceReader) *GnssGeofence
	// SetFenceValid sets the value of the FenceValid signal.
	SetFenceValid(bool) *GnssGeofence
	// SetFenceCombined sets the physical value of the FenceCombined signal.
	SetFenceCombined(float64) *GnssGeofence
	// SetFence1 sets the physical value of the Fence1 signal.
	SetFence1(float64) *GnssGeofence
	// SetFence2 sets the physical value of the Fence2 signal.
	SetFence2(float64) *GnssGeofence
	// SetFence3 sets the physical value of the Fence3 signal.
	SetFence3(float64) *GnssGeofence
	// SetFence4 sets the physical value of the Fence4 signal.
	SetFence4(float64) *GnssGeofence
}

type GnssGeofence struct {
	xxx_FenceValid    bool
	xxx_FenceCombined uint8
	xxx_Fence1        uint8
	xxx_Fence2        uint8
	xxx_Fence3        uint8
	xxx_Fence4        uint8
}

func NewGnssGeofence() *GnssGeofence {
	m := &GnssGeofence{}
	m.Reset()
	return m
}

func (m *GnssGeofence) Reset() {
	m.xxx_FenceValid = false
	m.xxx_FenceCombined = 0
	m.xxx_Fence1 = 0
	m.xxx_Fence2 = 0
	m.xxx_Fence3 = 0
	m.xxx_Fence4 = 0
}

func (m *GnssGeofence) CopyFrom(o GnssGeofenceReader) *GnssGeofence {
	m.xxx_FenceValid = o.FenceValid()
	m.SetFenceCombined(o.FenceCombined())
	m.SetFence1(o.Fence1())
	m.SetFence2(o.Fence2())
	m.SetFence3(o.Fence3())
	m.SetFence4(o.Fence4())
	return m
}

// Descriptor returns the GnssGeofence descriptor.
func (m *GnssGeofence) Descriptor() *descriptor.Message {
	return Messages().GnssGeofence.Message
}

// String returns a compact string representation of the message.
func (m *GnssGeofence) String() string {
	return cantext.MessageString(m)
}

func (m *GnssGeofence) FenceValid() bool {
	return m.xxx_FenceValid
}

func (m *GnssGeofence) SetFenceValid(v bool) *GnssGeofence {
	m.xxx_FenceValid = v
	return m
}

func (m *GnssGeofence) FenceCombined() float64 {
	return Messages().GnssGeofence.FenceCombined.ToPhysical(float64(m.xxx_FenceCombined))
}

func (m *GnssGeofence) SetFenceCombined(v float64) *GnssGeofence {
	m.xxx_FenceCombined = uint8(Messages().GnssGeofence.FenceCombined.FromPhysical(v))
	return m
}

func (m *GnssGeofence) Fence1() float64 {
	return Messages().GnssGeofence.Fence1.ToPhysical(float64(m.xxx_Fence1))
}

func (m *GnssGeofence) SetFence1(v float64) *GnssGeofence {
	m.xxx_Fence1 = uint8(Messages().GnssGeofence.Fence1.FromPhysical(v))
	return m
}

func (m *GnssGeofence) Fence2() float64 {
	return Messages().GnssGeofence.Fence2.ToPhysical(float64(m.xxx_Fence2))
}

func (m *GnssGeofence) SetFence2(v float64) *GnssGeofence {
	m.xxx_Fence2 = uint8(Messages().GnssGeofence.Fence2.FromPhysical(v))
	return m
}

func (m *GnssGeofence) Fence3() float64 {
	return Messages().GnssGeofence.Fence3.ToPhysical(float64(m.xxx_Fence3))
}

func (m *GnssGeofence) SetFence3(v float64) *GnssGeofence {
	m.xxx_Fence3 = uint8(Messages().GnssGeofence.Fence3.FromPhysical(v))
	return m
}

func (m *GnssGeofence) Fence4() float64 {
	return Messages().GnssGeofence.Fence4.ToPhysical(float64(m.xxx_Fence4))
}

func (m *GnssGeofence) SetFence4(v float64) *GnssGeofence {
	m.xxx_Fence4 = uint8(Messages().GnssGeofence.Fence4.FromPhysical(v))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *GnssGeofence) Frame() can.Frame {
	md := Messages().GnssGeofence
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.FenceValid.MarshalBool(&f.Data, bool(m.xxx_FenceValid))
	md.FenceCombined.MarshalUnsigned(&f.Data, uint64(m.xxx_FenceCombined))
	md.Fence1.MarshalUnsigned(&f.Data, uint64(m.xxx_Fence1))
	md.Fence2.MarshalUnsigned(&f.Data, uint64(m.xxx_Fence2))
	md.Fence3.MarshalUnsigned(&f.Data, uint64(m.xxx_Fence3))
	md.Fence4.MarshalUnsigned(&f.Data, uint64(m.xxx_Fence4))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *GnssGeofence) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *GnssGeofence) UnmarshalFrame(f can.Frame) error {
	md := Messages().GnssGeofence
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal GnssGeofence: expects ID 776 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal GnssGeofence: expects length 2 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal GnssGeofence: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal GnssGeofence: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_FenceValid = bool(md.FenceValid.UnmarshalBool(f.Data))
	m.xxx_FenceCombined = uint8(md.FenceCombined.UnmarshalUnsigned(f.Data))
	m.xxx_Fence1 = uint8(md.Fence1.UnmarshalUnsigned(f.Data))
	m.xxx_Fence2 = uint8(md.Fence2.UnmarshalUnsigned(f.Data))
	m.xxx_Fence3 = uint8(md.Fence3.UnmarshalUnsigned(f.Data))
	m.xxx_Fence4 = uint8(md.Fence4.UnmarshalUnsigned(f.Data))
	return nil
}

// GnssImuReader provides read access to a GnssImu message.
type GnssImuReader interface {
	// ImuValid returns the value of the ImuValid signal.
	ImuValid() bool
	// AccelerationX returns the physical value of the AccelerationX signal.
	AccelerationX() float64
	// AccelerationY returns the physical value of the AccelerationY signal.
	AccelerationY() float64
	// AccelerationZ returns the physical value of the AccelerationZ signal.
	AccelerationZ() float64
	// AngularRateX returns the physical value of the AngularRateX signal.
	AngularRateX() float64
	// AngularRateY returns the physical value of the AngularRateY signal.
	AngularRateY() float64
	// AngularRateZ returns the physical value of the AngularRateZ signal.
	AngularRateZ() float64
}

// GnssImuWriter provides write access to a GnssImu message.
type GnssImuWriter interface {
	// CopyFrom copies all values from GnssImuReader.
	CopyFrom(GnssImuReader) *GnssImu
	// SetImuValid sets the value of the ImuValid signal.
	SetImuValid(bool) *GnssImu
	// SetAccelerationX sets the physical value of the AccelerationX signal.
	SetAccelerationX(float64) *GnssImu
	// SetAccelerationY sets the physical value of the AccelerationY signal.
	SetAccelerationY(float64) *GnssImu
	// SetAccelerationZ sets the physical value of the AccelerationZ signal.
	SetAccelerationZ(float64) *GnssImu
	// SetAngularRateX sets the physical value of the AngularRateX signal.
	SetAngularRateX(float64) *GnssImu
	// SetAngularRateY sets the physical value of the AngularRateY signal.
	SetAngularRateY(float64) *GnssImu
	// SetAngularRateZ sets the physical value of the AngularRateZ signal.
	SetAngularRateZ(float64) *GnssImu
}

type GnssImu struct {
	xxx_ImuValid      bool
	xxx_AccelerationX uint16
	xxx_AccelerationY uint16
	xxx_AccelerationZ uint16
	xxx_AngularRateX  uint16
	xxx_AngularRateY  uint16
	xxx_AngularRateZ  uint16
}

func NewGnssImu() *GnssImu {
	m := &GnssImu{}
	m.Reset()
	return m
}

func (m *GnssImu) Reset() {
	m.xxx_ImuValid = false
	m.xxx_AccelerationX = 0
	m.xxx_AccelerationY = 0
	m.xxx_AccelerationZ = 0
	m.xxx_AngularRateX = 0
	m.xxx_AngularRateY = 0
	m.xxx_AngularRateZ = 0
}

func (m *GnssImu) CopyFrom(o GnssImuReader) *GnssImu {
	m.xxx_ImuValid = o.ImuValid()
	m.SetAccelerationX(o.AccelerationX())
	m.SetAccelerationY(o.AccelerationY())
	m.SetAccelerationZ(o.AccelerationZ())
	m.SetAngularRateX(o.AngularRateX())
	m.SetAngularRateY(o.AngularRateY())
	m.SetAngularRateZ(o.AngularRateZ())
	return m
}

// Descriptor returns the GnssImu descriptor.
func (m *GnssImu) Descriptor() *descriptor.Message {
	return Messages().GnssImu.Message
}

// String returns a compact string representation of the message.
func (m *GnssImu) String() string {
	return cantext.MessageString(m)
}

func (m *GnssImu) ImuValid() bool {
	return m.xxx_ImuValid
}

func (m *GnssImu) SetImuValid(v bool) *GnssImu {
	m.xxx_ImuValid = v
	return m
}

func (m *GnssImu) AccelerationX() float64 {
	return Messages().GnssImu.AccelerationX.ToPhysical(float64(m.xxx_AccelerationX))
}

func (m *GnssImu) SetAccelerationX(v float64) *GnssImu {
	m.xxx_AccelerationX = uint16(Messages().GnssImu.AccelerationX.FromPhysical(v))
	return m
}

func (m *GnssImu) AccelerationY() float64 {
	return Messages().GnssImu.AccelerationY.ToPhysical(float64(m.xxx_AccelerationY))
}

func (m *GnssImu) SetAccelerationY(v float64) *GnssImu {
	m.xxx_AccelerationY = uint16(Messages().GnssImu.AccelerationY.FromPhysical(v))
	return m
}

func (m *GnssImu) AccelerationZ() float64 {
	return Messages().GnssImu.AccelerationZ.ToPhysical(float64(m.xxx_AccelerationZ))
}

func (m *GnssImu) SetAccelerationZ(v float64) *GnssImu {
	m.xxx_AccelerationZ = uint16(Messages().GnssImu.AccelerationZ.FromPhysical(v))
	return m
}

func (m *GnssImu) AngularRateX() float64 {
	return Messages().GnssImu.AngularRateX.ToPhysical(float64(m.xxx_AngularRateX))
}

func (m *GnssImu) SetAngularRateX(v float64) *GnssImu {
	m.xxx_AngularRateX = uint16(Messages().GnssImu.AngularRateX.FromPhysical(v))
	return m
}

func (m *GnssImu) AngularRateY() float64 {
	return Messages().GnssImu.AngularRateY.ToPhysical(float64(m.xxx_AngularRateY))
}

func (m *GnssImu) SetAngularRateY(v float64) *GnssImu {
	m.xxx_AngularRateY = uint16(Messages().GnssImu.AngularRateY.FromPhysical(v))
	return m
}

func (m *GnssImu) AngularRateZ() float64 {
	return Messages().GnssImu.AngularRateZ.ToPhysical(float64(m.xxx_AngularRateZ))
}

func (m *GnssImu) SetAngularRateZ(v float64) *GnssImu {
	m.xxx_AngularRateZ = uint16(Messages().GnssImu.AngularRateZ.FromPhysical(v))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *GnssImu) Frame() can.Frame {
	md := Messages().GnssImu
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.ImuValid.MarshalBool(&f.Data, bool(m.xxx_ImuValid))
	md.AccelerationX.MarshalUnsigned(&f.Data, uint64(m.xxx_AccelerationX))
	md.AccelerationY.MarshalUnsigned(&f.Data, uint64(m.xxx_AccelerationY))
	md.AccelerationZ.MarshalUnsigned(&f.Data, uint64(m.xxx_AccelerationZ))
	md.AngularRateX.MarshalUnsigned(&f.Data, uint64(m.xxx_AngularRateX))
	md.AngularRateY.MarshalUnsigned(&f.Data, uint64(m.xxx_AngularRateY))
	md.AngularRateZ.MarshalUnsigned(&f.Data, uint64(m.xxx_AngularRateZ))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *GnssImu) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *GnssImu) UnmarshalFrame(f can.Frame) error {
	md := Messages().GnssImu
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal GnssImu: expects ID 777 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal GnssImu: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal GnssImu: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal GnssImu: expects standard ID (got %s with extended ID)", f.String(),
		)
	}
	m.xxx_ImuValid = bool(md.ImuValid.UnmarshalBool(f.Data))
	m.xxx_AccelerationX = uint16(md.AccelerationX.UnmarshalUnsigned(f.Data))
	m.xxx_AccelerationY = uint16(md.AccelerationY.UnmarshalUnsigned(f.Data))
	m.xxx_AccelerationZ = uint16(md.AccelerationZ.UnmarshalUnsigned(f.Data))
	m.xxx_AngularRateX = uint16(md.AngularRateX.UnmarshalUnsigned(f.Data))
	m.xxx_AngularRateY = uint16(md.AngularRateY.UnmarshalUnsigned(f.Data))
	m.xxx_AngularRateZ = uint16(md.AngularRateZ.UnmarshalUnsigned(f.Data))
	return nil
}

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
	// Pack_Positive_Feedback returns the value of the Pack_Positive_Feedback signal.
	Pack_Positive_Feedback() bool
	// Pack_Negative_Feedback returns the value of the Pack_Negative_Feedback signal.
	Pack_Negative_Feedback() bool
	// Pack_Precharge_Feedback returns the value of the Pack_Precharge_Feedback signal.
	Pack_Precharge_Feedback() bool
}

// Contactor_FeedbackWriter provides write access to a Contactor_Feedback message.
type Contactor_FeedbackWriter interface {
	// CopyFrom copies all values from Contactor_FeedbackReader.
	CopyFrom(Contactor_FeedbackReader) *Contactor_Feedback
	// SetPack_Positive_Feedback sets the value of the Pack_Positive_Feedback signal.
	SetPack_Positive_Feedback(bool) *Contactor_Feedback
	// SetPack_Negative_Feedback sets the value of the Pack_Negative_Feedback signal.
	SetPack_Negative_Feedback(bool) *Contactor_Feedback
	// SetPack_Precharge_Feedback sets the value of the Pack_Precharge_Feedback signal.
	SetPack_Precharge_Feedback(bool) *Contactor_Feedback
}

type Contactor_Feedback struct {
	xxx_Pack_Positive_Feedback  bool
	xxx_Pack_Negative_Feedback  bool
	xxx_Pack_Precharge_Feedback bool
}

func NewContactor_Feedback() *Contactor_Feedback {
	m := &Contactor_Feedback{}
	m.Reset()
	return m
}

func (m *Contactor_Feedback) Reset() {
	m.xxx_Pack_Positive_Feedback = false
	m.xxx_Pack_Negative_Feedback = false
	m.xxx_Pack_Precharge_Feedback = false
}

func (m *Contactor_Feedback) CopyFrom(o Contactor_FeedbackReader) *Contactor_Feedback {
	m.xxx_Pack_Positive_Feedback = o.Pack_Positive_Feedback()
	m.xxx_Pack_Negative_Feedback = o.Pack_Negative_Feedback()
	m.xxx_Pack_Precharge_Feedback = o.Pack_Precharge_Feedback()
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

func (m *Contactor_Feedback) Pack_Positive_Feedback() bool {
	return m.xxx_Pack_Positive_Feedback
}

func (m *Contactor_Feedback) SetPack_Positive_Feedback(v bool) *Contactor_Feedback {
	m.xxx_Pack_Positive_Feedback = v
	return m
}

func (m *Contactor_Feedback) Pack_Negative_Feedback() bool {
	return m.xxx_Pack_Negative_Feedback
}

func (m *Contactor_Feedback) SetPack_Negative_Feedback(v bool) *Contactor_Feedback {
	m.xxx_Pack_Negative_Feedback = v
	return m
}

func (m *Contactor_Feedback) Pack_Precharge_Feedback() bool {
	return m.xxx_Pack_Precharge_Feedback
}

func (m *Contactor_Feedback) SetPack_Precharge_Feedback(v bool) *Contactor_Feedback {
	m.xxx_Pack_Precharge_Feedback = v
	return m
}

// Frame returns a CAN frame representing the message.
func (m *Contactor_Feedback) Frame() can.Frame {
	md := Messages().Contactor_Feedback
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.Pack_Positive_Feedback.MarshalBool(&f.Data, bool(m.xxx_Pack_Positive_Feedback))
	md.Pack_Negative_Feedback.MarshalBool(&f.Data, bool(m.xxx_Pack_Negative_Feedback))
	md.Pack_Precharge_Feedback.MarshalBool(&f.Data, bool(m.xxx_Pack_Precharge_Feedback))
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
			"unmarshal Contactor_Feedback: expects length 3 (got %s with length %d)", f.String(), f.Length,
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
	m.xxx_Pack_Positive_Feedback = bool(md.Pack_Positive_Feedback.UnmarshalBool(f.Data))
	m.xxx_Pack_Negative_Feedback = bool(md.Pack_Negative_Feedback.UnmarshalBool(f.Data))
	m.xxx_Pack_Precharge_Feedback = bool(md.Pack_Precharge_Feedback.UnmarshalBool(f.Data))
	return nil
}

// BMSBroadcastReader provides read access to a BMSBroadcast message.
type BMSBroadcastReader interface {
	// ThermModuleNum returns the value of the ThermModuleNum signal.
	ThermModuleNum() uint8
	// LowThermValue returns the value of the LowThermValue signal.
	LowThermValue() int8
	// HighThermValue returns the value of the HighThermValue signal.
	HighThermValue() int8
	// AvgThermValue returns the value of the AvgThermValue signal.
	AvgThermValue() int8
	// NumThermEn returns the value of the NumThermEn signal.
	NumThermEn() uint8
	// HighThermID returns the value of the HighThermID signal.
	HighThermID() uint8
	// LowThermID returns the value of the LowThermID signal.
	LowThermID() uint8
	// Checksum returns the value of the Checksum signal.
	Checksum() int8
}

// BMSBroadcastWriter provides write access to a BMSBroadcast message.
type BMSBroadcastWriter interface {
	// CopyFrom copies all values from BMSBroadcastReader.
	CopyFrom(BMSBroadcastReader) *BMSBroadcast
	// SetThermModuleNum sets the value of the ThermModuleNum signal.
	SetThermModuleNum(uint8) *BMSBroadcast
	// SetLowThermValue sets the value of the LowThermValue signal.
	SetLowThermValue(int8) *BMSBroadcast
	// SetHighThermValue sets the value of the HighThermValue signal.
	SetHighThermValue(int8) *BMSBroadcast
	// SetAvgThermValue sets the value of the AvgThermValue signal.
	SetAvgThermValue(int8) *BMSBroadcast
	// SetNumThermEn sets the value of the NumThermEn signal.
	SetNumThermEn(uint8) *BMSBroadcast
	// SetHighThermID sets the value of the HighThermID signal.
	SetHighThermID(uint8) *BMSBroadcast
	// SetLowThermID sets the value of the LowThermID signal.
	SetLowThermID(uint8) *BMSBroadcast
	// SetChecksum sets the value of the Checksum signal.
	SetChecksum(int8) *BMSBroadcast
}

type BMSBroadcast struct {
	xxx_ThermModuleNum uint8
	xxx_LowThermValue  int8
	xxx_HighThermValue int8
	xxx_AvgThermValue  int8
	xxx_NumThermEn     uint8
	xxx_HighThermID    uint8
	xxx_LowThermID     uint8
	xxx_Checksum       int8
}

func NewBMSBroadcast() *BMSBroadcast {
	m := &BMSBroadcast{}
	m.Reset()
	return m
}

func (m *BMSBroadcast) Reset() {
	m.xxx_ThermModuleNum = 0
	m.xxx_LowThermValue = 0
	m.xxx_HighThermValue = 0
	m.xxx_AvgThermValue = 0
	m.xxx_NumThermEn = 0
	m.xxx_HighThermID = 0
	m.xxx_LowThermID = 0
	m.xxx_Checksum = 0
}

func (m *BMSBroadcast) CopyFrom(o BMSBroadcastReader) *BMSBroadcast {
	m.xxx_ThermModuleNum = o.ThermModuleNum()
	m.xxx_LowThermValue = o.LowThermValue()
	m.xxx_HighThermValue = o.HighThermValue()
	m.xxx_AvgThermValue = o.AvgThermValue()
	m.xxx_NumThermEn = o.NumThermEn()
	m.xxx_HighThermID = o.HighThermID()
	m.xxx_LowThermID = o.LowThermID()
	m.xxx_Checksum = o.Checksum()
	return m
}

// Descriptor returns the BMSBroadcast descriptor.
func (m *BMSBroadcast) Descriptor() *descriptor.Message {
	return Messages().BMSBroadcast.Message
}

// String returns a compact string representation of the message.
func (m *BMSBroadcast) String() string {
	return cantext.MessageString(m)
}

func (m *BMSBroadcast) ThermModuleNum() uint8 {
	return m.xxx_ThermModuleNum
}

func (m *BMSBroadcast) SetThermModuleNum(v uint8) *BMSBroadcast {
	m.xxx_ThermModuleNum = uint8(Messages().BMSBroadcast.ThermModuleNum.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *BMSBroadcast) LowThermValue() int8 {
	return m.xxx_LowThermValue
}

func (m *BMSBroadcast) SetLowThermValue(v int8) *BMSBroadcast {
	m.xxx_LowThermValue = int8(Messages().BMSBroadcast.LowThermValue.SaturatedCastSigned(int64(v)))
	return m
}

func (m *BMSBroadcast) HighThermValue() int8 {
	return m.xxx_HighThermValue
}

func (m *BMSBroadcast) SetHighThermValue(v int8) *BMSBroadcast {
	m.xxx_HighThermValue = int8(Messages().BMSBroadcast.HighThermValue.SaturatedCastSigned(int64(v)))
	return m
}

func (m *BMSBroadcast) AvgThermValue() int8 {
	return m.xxx_AvgThermValue
}

func (m *BMSBroadcast) SetAvgThermValue(v int8) *BMSBroadcast {
	m.xxx_AvgThermValue = int8(Messages().BMSBroadcast.AvgThermValue.SaturatedCastSigned(int64(v)))
	return m
}

func (m *BMSBroadcast) NumThermEn() uint8 {
	return m.xxx_NumThermEn
}

func (m *BMSBroadcast) SetNumThermEn(v uint8) *BMSBroadcast {
	m.xxx_NumThermEn = uint8(Messages().BMSBroadcast.NumThermEn.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *BMSBroadcast) HighThermID() uint8 {
	return m.xxx_HighThermID
}

func (m *BMSBroadcast) SetHighThermID(v uint8) *BMSBroadcast {
	m.xxx_HighThermID = uint8(Messages().BMSBroadcast.HighThermID.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *BMSBroadcast) LowThermID() uint8 {
	return m.xxx_LowThermID
}

func (m *BMSBroadcast) SetLowThermID(v uint8) *BMSBroadcast {
	m.xxx_LowThermID = uint8(Messages().BMSBroadcast.LowThermID.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *BMSBroadcast) Checksum() int8 {
	return m.xxx_Checksum
}

func (m *BMSBroadcast) SetChecksum(v int8) *BMSBroadcast {
	m.xxx_Checksum = int8(Messages().BMSBroadcast.Checksum.SaturatedCastSigned(int64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *BMSBroadcast) Frame() can.Frame {
	md := Messages().BMSBroadcast
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.ThermModuleNum.MarshalUnsigned(&f.Data, uint64(m.xxx_ThermModuleNum))
	md.LowThermValue.MarshalSigned(&f.Data, int64(m.xxx_LowThermValue))
	md.HighThermValue.MarshalSigned(&f.Data, int64(m.xxx_HighThermValue))
	md.AvgThermValue.MarshalSigned(&f.Data, int64(m.xxx_AvgThermValue))
	md.NumThermEn.MarshalUnsigned(&f.Data, uint64(m.xxx_NumThermEn))
	md.HighThermID.MarshalUnsigned(&f.Data, uint64(m.xxx_HighThermID))
	md.LowThermID.MarshalUnsigned(&f.Data, uint64(m.xxx_LowThermID))
	md.Checksum.MarshalSigned(&f.Data, int64(m.xxx_Checksum))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *BMSBroadcast) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *BMSBroadcast) UnmarshalFrame(f can.Frame) error {
	md := Messages().BMSBroadcast
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal BMSBroadcast: expects ID 406451072 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal BMSBroadcast: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal BMSBroadcast: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal BMSBroadcast: expects extended ID (got %s with standard ID)", f.String(),
		)
	}
	m.xxx_ThermModuleNum = uint8(md.ThermModuleNum.UnmarshalUnsigned(f.Data))
	m.xxx_LowThermValue = int8(md.LowThermValue.UnmarshalSigned(f.Data))
	m.xxx_HighThermValue = int8(md.HighThermValue.UnmarshalSigned(f.Data))
	m.xxx_AvgThermValue = int8(md.AvgThermValue.UnmarshalSigned(f.Data))
	m.xxx_NumThermEn = uint8(md.NumThermEn.UnmarshalUnsigned(f.Data))
	m.xxx_HighThermID = uint8(md.HighThermID.UnmarshalUnsigned(f.Data))
	m.xxx_LowThermID = uint8(md.LowThermID.UnmarshalUnsigned(f.Data))
	m.xxx_Checksum = int8(md.Checksum.UnmarshalSigned(f.Data))
	return nil
}

// ThermistorBroadcastReader provides read access to a ThermistorBroadcast message.
type ThermistorBroadcastReader interface {
	// RelThermID returns the value of the RelThermID signal.
	RelThermID() uint16
	// ThermValue returns the value of the ThermValue signal.
	ThermValue() int8
	// NumEnTherm returns the value of the NumEnTherm signal.
	NumEnTherm() int8
	// LowThermValue returns the value of the LowThermValue signal.
	LowThermValue() int8
	// HighThermValue returns the value of the HighThermValue signal.
	HighThermValue() int8
	// HighThermID returns the value of the HighThermID signal.
	HighThermID() uint8
	// LowThermID returns the value of the LowThermID signal.
	LowThermID() uint8
}

// ThermistorBroadcastWriter provides write access to a ThermistorBroadcast message.
type ThermistorBroadcastWriter interface {
	// CopyFrom copies all values from ThermistorBroadcastReader.
	CopyFrom(ThermistorBroadcastReader) *ThermistorBroadcast
	// SetRelThermID sets the value of the RelThermID signal.
	SetRelThermID(uint16) *ThermistorBroadcast
	// SetThermValue sets the value of the ThermValue signal.
	SetThermValue(int8) *ThermistorBroadcast
	// SetNumEnTherm sets the value of the NumEnTherm signal.
	SetNumEnTherm(int8) *ThermistorBroadcast
	// SetLowThermValue sets the value of the LowThermValue signal.
	SetLowThermValue(int8) *ThermistorBroadcast
	// SetHighThermValue sets the value of the HighThermValue signal.
	SetHighThermValue(int8) *ThermistorBroadcast
	// SetHighThermID sets the value of the HighThermID signal.
	SetHighThermID(uint8) *ThermistorBroadcast
	// SetLowThermID sets the value of the LowThermID signal.
	SetLowThermID(uint8) *ThermistorBroadcast
}

type ThermistorBroadcast struct {
	xxx_RelThermID     uint16
	xxx_ThermValue     int8
	xxx_NumEnTherm     int8
	xxx_LowThermValue  int8
	xxx_HighThermValue int8
	xxx_HighThermID    uint8
	xxx_LowThermID     uint8
}

func NewThermistorBroadcast() *ThermistorBroadcast {
	m := &ThermistorBroadcast{}
	m.Reset()
	return m
}

func (m *ThermistorBroadcast) Reset() {
	m.xxx_RelThermID = 0
	m.xxx_ThermValue = 0
	m.xxx_NumEnTherm = 0
	m.xxx_LowThermValue = 0
	m.xxx_HighThermValue = 0
	m.xxx_HighThermID = 0
	m.xxx_LowThermID = 0
}

func (m *ThermistorBroadcast) CopyFrom(o ThermistorBroadcastReader) *ThermistorBroadcast {
	m.xxx_RelThermID = o.RelThermID()
	m.xxx_ThermValue = o.ThermValue()
	m.xxx_NumEnTherm = o.NumEnTherm()
	m.xxx_LowThermValue = o.LowThermValue()
	m.xxx_HighThermValue = o.HighThermValue()
	m.xxx_HighThermID = o.HighThermID()
	m.xxx_LowThermID = o.LowThermID()
	return m
}

// Descriptor returns the ThermistorBroadcast descriptor.
func (m *ThermistorBroadcast) Descriptor() *descriptor.Message {
	return Messages().ThermistorBroadcast.Message
}

// String returns a compact string representation of the message.
func (m *ThermistorBroadcast) String() string {
	return cantext.MessageString(m)
}

func (m *ThermistorBroadcast) RelThermID() uint16 {
	return m.xxx_RelThermID
}

func (m *ThermistorBroadcast) SetRelThermID(v uint16) *ThermistorBroadcast {
	m.xxx_RelThermID = uint16(Messages().ThermistorBroadcast.RelThermID.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *ThermistorBroadcast) ThermValue() int8 {
	return m.xxx_ThermValue
}

func (m *ThermistorBroadcast) SetThermValue(v int8) *ThermistorBroadcast {
	m.xxx_ThermValue = int8(Messages().ThermistorBroadcast.ThermValue.SaturatedCastSigned(int64(v)))
	return m
}

func (m *ThermistorBroadcast) NumEnTherm() int8 {
	return m.xxx_NumEnTherm
}

func (m *ThermistorBroadcast) SetNumEnTherm(v int8) *ThermistorBroadcast {
	m.xxx_NumEnTherm = int8(Messages().ThermistorBroadcast.NumEnTherm.SaturatedCastSigned(int64(v)))
	return m
}

func (m *ThermistorBroadcast) LowThermValue() int8 {
	return m.xxx_LowThermValue
}

func (m *ThermistorBroadcast) SetLowThermValue(v int8) *ThermistorBroadcast {
	m.xxx_LowThermValue = int8(Messages().ThermistorBroadcast.LowThermValue.SaturatedCastSigned(int64(v)))
	return m
}

func (m *ThermistorBroadcast) HighThermValue() int8 {
	return m.xxx_HighThermValue
}

func (m *ThermistorBroadcast) SetHighThermValue(v int8) *ThermistorBroadcast {
	m.xxx_HighThermValue = int8(Messages().ThermistorBroadcast.HighThermValue.SaturatedCastSigned(int64(v)))
	return m
}

func (m *ThermistorBroadcast) HighThermID() uint8 {
	return m.xxx_HighThermID
}

func (m *ThermistorBroadcast) SetHighThermID(v uint8) *ThermistorBroadcast {
	m.xxx_HighThermID = uint8(Messages().ThermistorBroadcast.HighThermID.SaturatedCastUnsigned(uint64(v)))
	return m
}

func (m *ThermistorBroadcast) LowThermID() uint8 {
	return m.xxx_LowThermID
}

func (m *ThermistorBroadcast) SetLowThermID(v uint8) *ThermistorBroadcast {
	m.xxx_LowThermID = uint8(Messages().ThermistorBroadcast.LowThermID.SaturatedCastUnsigned(uint64(v)))
	return m
}

// Frame returns a CAN frame representing the message.
func (m *ThermistorBroadcast) Frame() can.Frame {
	md := Messages().ThermistorBroadcast
	f := can.Frame{ID: md.ID, IsExtended: md.IsExtended, Length: md.Length}
	md.RelThermID.MarshalUnsigned(&f.Data, uint64(m.xxx_RelThermID))
	md.ThermValue.MarshalSigned(&f.Data, int64(m.xxx_ThermValue))
	md.NumEnTherm.MarshalSigned(&f.Data, int64(m.xxx_NumEnTherm))
	md.LowThermValue.MarshalSigned(&f.Data, int64(m.xxx_LowThermValue))
	md.HighThermValue.MarshalSigned(&f.Data, int64(m.xxx_HighThermValue))
	md.HighThermID.MarshalUnsigned(&f.Data, uint64(m.xxx_HighThermID))
	md.LowThermID.MarshalUnsigned(&f.Data, uint64(m.xxx_LowThermID))
	return f
}

// MarshalFrame encodes the message as a CAN frame.
func (m *ThermistorBroadcast) MarshalFrame() (can.Frame, error) {
	return m.Frame(), nil
}

// UnmarshalFrame decodes the message from a CAN frame.
func (m *ThermistorBroadcast) UnmarshalFrame(f can.Frame) error {
	md := Messages().ThermistorBroadcast
	switch {
	case f.ID != md.ID:
		return fmt.Errorf(
			"unmarshal ThermistorBroadcast: expects ID 419361278 (got %s with ID %d)", f.String(), f.ID,
		)
	case f.Length != md.Length:
		return fmt.Errorf(
			"unmarshal ThermistorBroadcast: expects length 8 (got %s with length %d)", f.String(), f.Length,
		)
	case f.IsRemote:
		return fmt.Errorf(
			"unmarshal ThermistorBroadcast: expects non-remote frame (got remote frame %s)", f.String(),
		)
	case f.IsExtended != md.IsExtended:
		return fmt.Errorf(
			"unmarshal ThermistorBroadcast: expects extended ID (got %s with standard ID)", f.String(),
		)
	}
	m.xxx_RelThermID = uint16(md.RelThermID.UnmarshalUnsigned(f.Data))
	m.xxx_ThermValue = int8(md.ThermValue.UnmarshalSigned(f.Data))
	m.xxx_NumEnTherm = int8(md.NumEnTherm.UnmarshalSigned(f.Data))
	m.xxx_LowThermValue = int8(md.LowThermValue.UnmarshalSigned(f.Data))
	m.xxx_HighThermValue = int8(md.HighThermValue.UnmarshalSigned(f.Data))
	m.xxx_HighThermID = uint8(md.HighThermID.UnmarshalUnsigned(f.Data))
	m.xxx_LowThermID = uint8(md.LowThermID.UnmarshalUnsigned(f.Data))
	return nil
}

// Nodes returns the veh node descriptors.
func Nodes() *NodesDescriptor {
	return nd
}

// NodesDescriptor contains all veh node descriptors.
type NodesDescriptor struct {
	BMS    *descriptor.Node
	CANmod *descriptor.Node
	FC     *descriptor.Node
	PC_SG  *descriptor.Node
	TMS    *descriptor.Node
}

// Messages returns the veh message descriptors.
func Messages() *MessagesDescriptor {
	return md
}

// MessagesDescriptor contains all veh message descriptors.
type MessagesDescriptor struct {
	VC_Status           *VC_StatusDescriptor
	FC_cmd              *FC_cmdDescriptor
	FC_msg              *FC_msgDescriptor
	GnssStatus          *GnssStatusDescriptor
	GnssTime            *GnssTimeDescriptor
	GnssPosition        *GnssPositionDescriptor
	GnssAltitude        *GnssAltitudeDescriptor
	GnssAttitude        *GnssAttitudeDescriptor
	GnssOdo             *GnssOdoDescriptor
	GnssSpeed           *GnssSpeedDescriptor
	GnssGeofence        *GnssGeofenceDescriptor
	GnssImu             *GnssImuDescriptor
	Contactor_States    *Contactor_StatesDescriptor
	Pack_Current_Limits *Pack_Current_LimitsDescriptor
	Pack_State          *Pack_StateDescriptor
	Pack_SOC            *Pack_SOCDescriptor
	Contactor_Feedback  *Contactor_FeedbackDescriptor
	BMSBroadcast        *BMSBroadcastDescriptor
	ThermistorBroadcast *ThermistorBroadcastDescriptor
}

// UnmarshalFrame unmarshals the provided veh CAN frame.
func (md *MessagesDescriptor) UnmarshalFrame(f can.Frame) (generated.Message, error) {
	switch f.ID {
	case md.VC_Status.ID:
		var msg VC_Status
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.FC_cmd.ID:
		var msg FC_cmd
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.FC_msg.ID:
		var msg FC_msg
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.GnssStatus.ID:
		var msg GnssStatus
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.GnssTime.ID:
		var msg GnssTime
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.GnssPosition.ID:
		var msg GnssPosition
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.GnssAltitude.ID:
		var msg GnssAltitude
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.GnssAttitude.ID:
		var msg GnssAttitude
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.GnssOdo.ID:
		var msg GnssOdo
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.GnssSpeed.ID:
		var msg GnssSpeed
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.GnssGeofence.ID:
		var msg GnssGeofence
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.GnssImu.ID:
		var msg GnssImu
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.Contactor_States.ID:
		var msg Contactor_States
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.Pack_Current_Limits.ID:
		var msg Pack_Current_Limits
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.Pack_State.ID:
		var msg Pack_State
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.Pack_SOC.ID:
		var msg Pack_SOC
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.Contactor_Feedback.ID:
		var msg Contactor_Feedback
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.BMSBroadcast.ID:
		var msg BMSBroadcast
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	case md.ThermistorBroadcast.ID:
		var msg ThermistorBroadcast
		if err := msg.UnmarshalFrame(f); err != nil {
			return nil, fmt.Errorf("unmarshal veh frame: %w", err)
		}
		return &msg, nil
	default:
		return nil, fmt.Errorf("unmarshal veh frame: ID not in database: %d", f.ID)
	}
}

type VC_StatusDescriptor struct {
	*descriptor.Message
	VC_govStatus *descriptor.Signal
}

type FC_cmdDescriptor struct {
	*descriptor.Message
	FC_brakeLight          *descriptor.Signal
	FC_readyToDriveSpeaker *descriptor.Signal
}

type FC_msgDescriptor struct {
	*descriptor.Message
	FC_apps1             *descriptor.Signal
	FC_apps2             *descriptor.Signal
	FC_bpps              *descriptor.Signal
	FC_steeringAngle     *descriptor.Signal
	FC_hvilSts           *descriptor.Signal
	FC_readyToDriveBtn_n *descriptor.Signal
}

type GnssStatusDescriptor struct {
	*descriptor.Message
	FixType    *descriptor.Signal
	Satellites *descriptor.Signal
}

type GnssTimeDescriptor struct {
	*descriptor.Message
	TimeValid     *descriptor.Signal
	TimeConfirmed *descriptor.Signal
	Epoch         *descriptor.Signal
}

type GnssPositionDescriptor struct {
	*descriptor.Message
	PositionValid    *descriptor.Signal
	Latitude         *descriptor.Signal
	Longitude        *descriptor.Signal
	PositionAccuracy *descriptor.Signal
}

type GnssAltitudeDescriptor struct {
	*descriptor.Message
	AltitudeValid    *descriptor.Signal
	Altitude         *descriptor.Signal
	AltitudeAccuracy *descriptor.Signal
}

type GnssAttitudeDescriptor struct {
	*descriptor.Message
	AttitudeValid   *descriptor.Signal
	Roll            *descriptor.Signal
	RollAccuracy    *descriptor.Signal
	Pitch           *descriptor.Signal
	PitchAccuracy   *descriptor.Signal
	Heading         *descriptor.Signal
	HeadingAccuracy *descriptor.Signal
}

type GnssOdoDescriptor struct {
	*descriptor.Message
	DistanceValid    *descriptor.Signal
	DistanceTrip     *descriptor.Signal
	DistanceAccuracy *descriptor.Signal
	DistanceTotal    *descriptor.Signal
}

type GnssSpeedDescriptor struct {
	*descriptor.Message
	SpeedValid    *descriptor.Signal
	Speed         *descriptor.Signal
	SpeedAccuracy *descriptor.Signal
}

type GnssGeofenceDescriptor struct {
	*descriptor.Message
	FenceValid    *descriptor.Signal
	FenceCombined *descriptor.Signal
	Fence1        *descriptor.Signal
	Fence2        *descriptor.Signal
	Fence3        *descriptor.Signal
	Fence4        *descriptor.Signal
}

type GnssImuDescriptor struct {
	*descriptor.Message
	ImuValid      *descriptor.Signal
	AccelerationX *descriptor.Signal
	AccelerationY *descriptor.Signal
	AccelerationZ *descriptor.Signal
	AngularRateX  *descriptor.Signal
	AngularRateY  *descriptor.Signal
	AngularRateZ  *descriptor.Signal
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
	Pack_Positive_Feedback  *descriptor.Signal
	Pack_Negative_Feedback  *descriptor.Signal
	Pack_Precharge_Feedback *descriptor.Signal
}

type BMSBroadcastDescriptor struct {
	*descriptor.Message
	ThermModuleNum *descriptor.Signal
	LowThermValue  *descriptor.Signal
	HighThermValue *descriptor.Signal
	AvgThermValue  *descriptor.Signal
	NumThermEn     *descriptor.Signal
	HighThermID    *descriptor.Signal
	LowThermID     *descriptor.Signal
	Checksum       *descriptor.Signal
}

type ThermistorBroadcastDescriptor struct {
	*descriptor.Message
	RelThermID     *descriptor.Signal
	ThermValue     *descriptor.Signal
	NumEnTherm     *descriptor.Signal
	LowThermValue  *descriptor.Signal
	HighThermValue *descriptor.Signal
	HighThermID    *descriptor.Signal
	LowThermID     *descriptor.Signal
}

// Database returns the veh database descriptor.
func (md *MessagesDescriptor) Database() *descriptor.Database {
	return d
}

var nd = &NodesDescriptor{
	BMS:    d.Nodes[0],
	CANmod: d.Nodes[1],
	FC:     d.Nodes[2],
	PC_SG:  d.Nodes[3],
	TMS:    d.Nodes[4],
}

var md = &MessagesDescriptor{
	VC_Status: &VC_StatusDescriptor{
		Message:      d.Messages[0],
		VC_govStatus: d.Messages[0].Signals[0],
	},
	FC_cmd: &FC_cmdDescriptor{
		Message:                d.Messages[1],
		FC_brakeLight:          d.Messages[1].Signals[0],
		FC_readyToDriveSpeaker: d.Messages[1].Signals[1],
	},
	FC_msg: &FC_msgDescriptor{
		Message:              d.Messages[2],
		FC_apps1:             d.Messages[2].Signals[0],
		FC_apps2:             d.Messages[2].Signals[1],
		FC_bpps:              d.Messages[2].Signals[2],
		FC_steeringAngle:     d.Messages[2].Signals[3],
		FC_hvilSts:           d.Messages[2].Signals[4],
		FC_readyToDriveBtn_n: d.Messages[2].Signals[5],
	},
	GnssStatus: &GnssStatusDescriptor{
		Message:    d.Messages[3],
		FixType:    d.Messages[3].Signals[0],
		Satellites: d.Messages[3].Signals[1],
	},
	GnssTime: &GnssTimeDescriptor{
		Message:       d.Messages[4],
		TimeValid:     d.Messages[4].Signals[0],
		TimeConfirmed: d.Messages[4].Signals[1],
		Epoch:         d.Messages[4].Signals[2],
	},
	GnssPosition: &GnssPositionDescriptor{
		Message:          d.Messages[5],
		PositionValid:    d.Messages[5].Signals[0],
		Latitude:         d.Messages[5].Signals[1],
		Longitude:        d.Messages[5].Signals[2],
		PositionAccuracy: d.Messages[5].Signals[3],
	},
	GnssAltitude: &GnssAltitudeDescriptor{
		Message:          d.Messages[6],
		AltitudeValid:    d.Messages[6].Signals[0],
		Altitude:         d.Messages[6].Signals[1],
		AltitudeAccuracy: d.Messages[6].Signals[2],
	},
	GnssAttitude: &GnssAttitudeDescriptor{
		Message:         d.Messages[7],
		AttitudeValid:   d.Messages[7].Signals[0],
		Roll:            d.Messages[7].Signals[1],
		RollAccuracy:    d.Messages[7].Signals[2],
		Pitch:           d.Messages[7].Signals[3],
		PitchAccuracy:   d.Messages[7].Signals[4],
		Heading:         d.Messages[7].Signals[5],
		HeadingAccuracy: d.Messages[7].Signals[6],
	},
	GnssOdo: &GnssOdoDescriptor{
		Message:          d.Messages[8],
		DistanceValid:    d.Messages[8].Signals[0],
		DistanceTrip:     d.Messages[8].Signals[1],
		DistanceAccuracy: d.Messages[8].Signals[2],
		DistanceTotal:    d.Messages[8].Signals[3],
	},
	GnssSpeed: &GnssSpeedDescriptor{
		Message:       d.Messages[9],
		SpeedValid:    d.Messages[9].Signals[0],
		Speed:         d.Messages[9].Signals[1],
		SpeedAccuracy: d.Messages[9].Signals[2],
	},
	GnssGeofence: &GnssGeofenceDescriptor{
		Message:       d.Messages[10],
		FenceValid:    d.Messages[10].Signals[0],
		FenceCombined: d.Messages[10].Signals[1],
		Fence1:        d.Messages[10].Signals[2],
		Fence2:        d.Messages[10].Signals[3],
		Fence3:        d.Messages[10].Signals[4],
		Fence4:        d.Messages[10].Signals[5],
	},
	GnssImu: &GnssImuDescriptor{
		Message:       d.Messages[11],
		ImuValid:      d.Messages[11].Signals[0],
		AccelerationX: d.Messages[11].Signals[1],
		AccelerationY: d.Messages[11].Signals[2],
		AccelerationZ: d.Messages[11].Signals[3],
		AngularRateX:  d.Messages[11].Signals[4],
		AngularRateY:  d.Messages[11].Signals[5],
		AngularRateZ:  d.Messages[11].Signals[6],
	},
	Contactor_States: &Contactor_StatesDescriptor{
		Message:        d.Messages[12],
		Pack_Positive:  d.Messages[12].Signals[0],
		Pack_Precharge: d.Messages[12].Signals[1],
		Pack_Negative:  d.Messages[12].Signals[2],
	},
	Pack_Current_Limits: &Pack_Current_LimitsDescriptor{
		Message:  d.Messages[13],
		Pack_CCL: d.Messages[13].Signals[0],
		Pack_DCL: d.Messages[13].Signals[1],
	},
	Pack_State: &Pack_StateDescriptor{
		Message:           d.Messages[14],
		Pack_Current:      d.Messages[14].Signals[0],
		Pack_Inst_Voltage: d.Messages[14].Signals[1],
		Avg_Cell_Voltage:  d.Messages[14].Signals[2],
		Populated_Cells:   d.Messages[14].Signals[3],
	},
	Pack_SOC: &Pack_SOCDescriptor{
		Message:              d.Messages[15],
		Pack_SOC:             d.Messages[15].Signals[0],
		Maximum_Pack_Voltage: d.Messages[15].Signals[1],
	},
	Contactor_Feedback: &Contactor_FeedbackDescriptor{
		Message:                 d.Messages[16],
		Pack_Positive_Feedback:  d.Messages[16].Signals[0],
		Pack_Negative_Feedback:  d.Messages[16].Signals[1],
		Pack_Precharge_Feedback: d.Messages[16].Signals[2],
	},
	BMSBroadcast: &BMSBroadcastDescriptor{
		Message:        d.Messages[17],
		ThermModuleNum: d.Messages[17].Signals[0],
		LowThermValue:  d.Messages[17].Signals[1],
		HighThermValue: d.Messages[17].Signals[2],
		AvgThermValue:  d.Messages[17].Signals[3],
		NumThermEn:     d.Messages[17].Signals[4],
		HighThermID:    d.Messages[17].Signals[5],
		LowThermID:     d.Messages[17].Signals[6],
		Checksum:       d.Messages[17].Signals[7],
	},
	ThermistorBroadcast: &ThermistorBroadcastDescriptor{
		Message:        d.Messages[18],
		RelThermID:     d.Messages[18].Signals[0],
		ThermValue:     d.Messages[18].Signals[1],
		NumEnTherm:     d.Messages[18].Signals[2],
		LowThermValue:  d.Messages[18].Signals[3],
		HighThermValue: d.Messages[18].Signals[4],
		HighThermID:    d.Messages[18].Signals[5],
		LowThermID:     d.Messages[18].Signals[6],
	},
}

var d = (*descriptor.Database)(&descriptor.Database{
	SourceFile: (string)("temp/vehcan/veh.dbc"),
	Version:    (string)(""),
	Messages: ([]*descriptor.Message)([]*descriptor.Message{
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("VC_Status"),
			ID:          (uint32)(255),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:             (string)("VC_govStatus"),
					Start:            (uint8)(0),
					Length:           (uint8)(3),
					IsBigEndian:      (bool)(false),
					IsSigned:         (bool)(false),
					IsMultiplexer:    (bool)(false),
					IsMultiplexed:    (bool)(false),
					MultiplexerValue: (uint)(0),
					Offset:           (float64)(0),
					Scale:            (float64)(1),
					Min:              (float64)(0),
					Max:              (float64)(7),
					Unit:             (string)(""),
					Description:      (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)([]*descriptor.ValueDescription{
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(0),
							Description: (string)("gov_init"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(1),
							Description: (string)("gov_startup"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(2),
							Description: (string)("gov_running"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(3),
							Description: (string)("hv_startup_error"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(4),
							Description: (string)("motor_startup_error"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(5),
							Description: (string)("driver_interface_error"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(6),
							Description: (string)("hv_run_error"),
						}),
						(*descriptor.ValueDescription)(&descriptor.ValueDescription{
							Value:       (int64)(7),
							Description: (string)("motor_run_error"),
						}),
					}),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("PC_SG"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("FC_cmd"),
			ID:          (uint32)(256),
			IsExtended:  (bool)(false),
			Length:      (uint8)(1),
			SendType:    (descriptor.SendType)(0),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("FC_brakeLight"),
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
					Name:              (string)("FC_readyToDriveSpeaker"),
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
			}),
			SenderNode: (string)("PC_SG"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("FC_msg"),
			ID:          (uint32)(511),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)(""),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("FC_apps1"),
					Start:             (uint8)(0),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(4095),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("PC_SG"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("FC_apps2"),
					Start:             (uint8)(16),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(4095),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("PC_SG"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("FC_bpps"),
					Start:             (uint8)(32),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(4095),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("PC_SG"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("FC_steeringAngle"),
					Start:             (uint8)(48),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(4095),
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("PC_SG"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("FC_hvilSts"),
					Start:             (uint8)(60),
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
						(string)("PC_SG"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("FC_readyToDriveBtn_n"),
					Start:             (uint8)(61),
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
						(string)("PC_SG"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("FC"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("GnssStatus"),
			ID:          (uint32)(769),
			IsExtended:  (bool)(false),
			Length:      (uint8)(1),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("GNSS information"),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("FixType"),
					Start:             (uint8)(0),
					Length:            (uint8)(3),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(5),
					Unit:              (string)(""),
					Description:       (string)("Fix type"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Satellites"),
					Start:             (uint8)(3),
					Length:            (uint8)(5),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(31),
					Unit:              (string)(""),
					Description:       (string)("Number of satellites used"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("CANmod"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("GnssTime"),
			ID:          (uint32)(770),
			IsExtended:  (bool)(false),
			Length:      (uint8)(6),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("GNSS time"),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("TimeValid"),
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
					Description:       (string)("Time validity"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("TimeConfirmed"),
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
					Description:       (string)("Time confirmed"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Epoch"),
					Start:             (uint8)(8),
					Length:            (uint8)(40),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(1.5778404e+09),
					Scale:             (float64)(0.001),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)("sec"),
					Description:       (string)("Epoch time"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("CANmod"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("GnssPosition"),
			ID:          (uint32)(771),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("GNSS position"),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("PositionValid"),
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
					Description:       (string)("Position validity"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Latitude"),
					Start:             (uint8)(1),
					Length:            (uint8)(28),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(-90),
					Scale:             (float64)(1e-06),
					Min:               (float64)(-90),
					Max:               (float64)(178.435455),
					Unit:              (string)("deg"),
					Description:       (string)("Latitude"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Longitude"),
					Start:             (uint8)(29),
					Length:            (uint8)(29),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(-180),
					Scale:             (float64)(1e-06),
					Min:               (float64)(-180),
					Max:               (float64)(356.870911),
					Unit:              (string)("deg"),
					Description:       (string)("Longitude"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("PositionAccuracy"),
					Start:             (uint8)(58),
					Length:            (uint8)(6),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(63),
					Unit:              (string)("m"),
					Description:       (string)("Accuracy of position"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("CANmod"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("GnssAltitude"),
			ID:          (uint32)(772),
			IsExtended:  (bool)(false),
			Length:      (uint8)(4),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("GNSS altitude"),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AltitudeValid"),
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
					Description:       (string)("Altitude validity"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Altitude"),
					Start:             (uint8)(1),
					Length:            (uint8)(18),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(-6000),
					Scale:             (float64)(0.1),
					Min:               (float64)(-6000),
					Max:               (float64)(20000),
					Unit:              (string)("m"),
					Description:       (string)("Altitude"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AltitudeAccuracy"),
					Start:             (uint8)(19),
					Length:            (uint8)(13),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(8000),
					Unit:              (string)("m"),
					Description:       (string)("Accuracy of altitude"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("CANmod"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("GnssAttitude"),
			ID:          (uint32)(773),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("GNSS attitude"),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AttitudeValid"),
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
					Description:       (string)("Attitude validity"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Roll"),
					Start:             (uint8)(1),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(-180),
					Scale:             (float64)(0.1),
					Min:               (float64)(-180),
					Max:               (float64)(180),
					Unit:              (string)("deg"),
					Description:       (string)("Vehicle roll"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("RollAccuracy"),
					Start:             (uint8)(13),
					Length:            (uint8)(9),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(50),
					Unit:              (string)("deg"),
					Description:       (string)("Vehicle roll accuracy"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Pitch"),
					Start:             (uint8)(22),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(-90),
					Scale:             (float64)(0.1),
					Min:               (float64)(-90),
					Max:               (float64)(90),
					Unit:              (string)("deg"),
					Description:       (string)("Vehicle pitch"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("PitchAccuracy"),
					Start:             (uint8)(34),
					Length:            (uint8)(9),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(50),
					Unit:              (string)("deg"),
					Description:       (string)("Vehicle pitch accuracy"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Heading"),
					Start:             (uint8)(43),
					Length:            (uint8)(12),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(360),
					Unit:              (string)("deg"),
					Description:       (string)("Vehicle heading"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("HeadingAccuracy"),
					Start:             (uint8)(55),
					Length:            (uint8)(9),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.1),
					Min:               (float64)(0),
					Max:               (float64)(50),
					Unit:              (string)("deg"),
					Description:       (string)("Vehicle heading accuracy"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("CANmod"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("GnssOdo"),
			ID:          (uint32)(774),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("GNSS odometer"),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("DistanceValid"),
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
					Name:              (string)("DistanceTrip"),
					Start:             (uint8)(1),
					Length:            (uint8)(22),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(4.194303e+06),
					Unit:              (string)("m"),
					Description:       (string)("Distance traveled since last reset"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("DistanceAccuracy"),
					Start:             (uint8)(23),
					Length:            (uint8)(19),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(524287),
					Unit:              (string)("m"),
					Description:       (string)("Distance accuracy (1-sigma)"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("DistanceTotal"),
					Start:             (uint8)(42),
					Length:            (uint8)(22),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(4.194303e+06),
					Unit:              (string)("km"),
					Description:       (string)("Distance traveled in total"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("CANmod"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("GnssSpeed"),
			ID:          (uint32)(775),
			IsExtended:  (bool)(false),
			Length:      (uint8)(5),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("GNSS speed"),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("SpeedValid"),
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
					Description:       (string)("Speed valid"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Speed"),
					Start:             (uint8)(1),
					Length:            (uint8)(20),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.001),
					Min:               (float64)(0),
					Max:               (float64)(1048.575),
					Unit:              (string)("m/s"),
					Description:       (string)("Speed m/s"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("SpeedAccuracy"),
					Start:             (uint8)(21),
					Length:            (uint8)(19),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(0.001),
					Min:               (float64)(0),
					Max:               (float64)(524.287),
					Unit:              (string)("m/s"),
					Description:       (string)("Speed accuracy"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("CANmod"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("GnssGeofence"),
			ID:          (uint32)(776),
			IsExtended:  (bool)(false),
			Length:      (uint8)(2),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("GNSS geofence(s)"),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("FenceValid"),
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
					Description:       (string)("Geofencing status"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("FenceCombined"),
					Start:             (uint8)(1),
					Length:            (uint8)(2),
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
					Description:       (string)("Combined (logical OR) state of all geofences"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Fence1"),
					Start:             (uint8)(8),
					Length:            (uint8)(2),
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
					Description:       (string)("Geofence 1 state"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Fence2"),
					Start:             (uint8)(10),
					Length:            (uint8)(2),
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
					Description:       (string)("Geofence 2 state"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Fence3"),
					Start:             (uint8)(12),
					Length:            (uint8)(2),
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
					Description:       (string)("Geofence 3 state"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Fence4"),
					Start:             (uint8)(14),
					Length:            (uint8)(2),
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
					Description:       (string)("Geofence 4 state"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("CANmod"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("GnssImu"),
			ID:          (uint32)(777),
			IsExtended:  (bool)(false),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("GNSS IMU"),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("ImuValid"),
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
					Name:              (string)("AccelerationX"),
					Start:             (uint8)(1),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(-64),
					Scale:             (float64)(0.125),
					Min:               (float64)(-64),
					Max:               (float64)(63.875),
					Unit:              (string)("m/s^2"),
					Description:       (string)("X acceleration with a resolution of 0.125 m/s^2"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AccelerationY"),
					Start:             (uint8)(11),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(-64),
					Scale:             (float64)(0.125),
					Min:               (float64)(-64),
					Max:               (float64)(63.875),
					Unit:              (string)("m/s^2"),
					Description:       (string)("Y acceleration with a resolution of 0.125 m/s^2"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AccelerationZ"),
					Start:             (uint8)(21),
					Length:            (uint8)(10),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(-64),
					Scale:             (float64)(0.125),
					Min:               (float64)(-64),
					Max:               (float64)(63.875),
					Unit:              (string)("m/s^2"),
					Description:       (string)("Z acceleration with a resolution of 0.125 m/s^2"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AngularRateX"),
					Start:             (uint8)(31),
					Length:            (uint8)(11),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(-256),
					Scale:             (float64)(0.25),
					Min:               (float64)(-256),
					Max:               (float64)(255.75),
					Unit:              (string)("deg/s"),
					Description:       (string)("X angular rate with a resolution of 0.25 deg/s"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AngularRateY"),
					Start:             (uint8)(42),
					Length:            (uint8)(11),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(-256),
					Scale:             (float64)(0.25),
					Min:               (float64)(-256),
					Max:               (float64)(255.75),
					Unit:              (string)("deg/s"),
					Description:       (string)("Y angular rate with a resolution of 0.25 deg/s"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AngularRateZ"),
					Start:             (uint8)(53),
					Length:            (uint8)(11),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(false),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(-256),
					Scale:             (float64)(0.25),
					Min:               (float64)(-256),
					Max:               (float64)(255.75),
					Unit:              (string)("deg/s"),
					Description:       (string)("Z angular rate with a resolution of 0.25 deg/s"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("FC"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("CANmod"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
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
			Length:      (uint8)(3),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("This ID Transmits at 8 ms."),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Pack_Positive_Feedback"),
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
					Name:              (string)("Pack_Precharge_Feedback"),
					Start:             (uint8)(16),
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
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("BMSBroadcast"),
			ID:          (uint32)(406451072),
			IsExtended:  (bool)(true),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("Thermistor Module - BMS Broadcast"),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("ThermModuleNum"),
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
					Description:       (string)("Thermistor Module Number"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("LowThermValue"),
					Start:             (uint8)(8),
					Length:            (uint8)(8),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(" C"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("HighThermValue"),
					Start:             (uint8)(16),
					Length:            (uint8)(8),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(" C"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("AvgThermValue"),
					Start:             (uint8)(24),
					Length:            (uint8)(8),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(" C"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("NumThermEn"),
					Start:             (uint8)(32),
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
					Name:              (string)("HighThermID"),
					Start:             (uint8)(40),
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
					Name:              (string)("LowThermID"),
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
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("Checksum"),
					Start:             (uint8)(56),
					Length:            (uint8)(8),
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
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
			}),
			SenderNode: (string)("TMS"),
			CycleTime:  (time.Duration)(0),
			DelayTime:  (time.Duration)(0),
		}),
		(*descriptor.Message)(&descriptor.Message{
			Name:        (string)("ThermistorBroadcast"),
			ID:          (uint32)(419361278),
			IsExtended:  (bool)(true),
			Length:      (uint8)(8),
			SendType:    (descriptor.SendType)(0),
			Description: (string)("Thermistor General Broadcast"),
			Signals: ([]*descriptor.Signal)([]*descriptor.Signal{
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("RelThermID"),
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
					Unit:              (string)(""),
					Description:       (string)("Thermistor ID relative to all configured Thermistor Modules"),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("ThermValue"),
					Start:             (uint8)(16),
					Length:            (uint8)(8),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(" C"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("NumEnTherm"),
					Start:             (uint8)(24),
					Length:            (uint8)(8),
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
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("LowThermValue"),
					Start:             (uint8)(32),
					Length:            (uint8)(8),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(" C"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("HighThermValue"),
					Start:             (uint8)(40),
					Length:            (uint8)(8),
					IsBigEndian:       (bool)(false),
					IsSigned:          (bool)(true),
					IsMultiplexer:     (bool)(false),
					IsMultiplexed:     (bool)(false),
					MultiplexerValue:  (uint)(0),
					Offset:            (float64)(0),
					Scale:             (float64)(1),
					Min:               (float64)(0),
					Max:               (float64)(0),
					Unit:              (string)(" C"),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("HighThermID"),
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
					Unit:              (string)(""),
					Description:       (string)(""),
					ValueDescriptions: ([]*descriptor.ValueDescription)(nil),
					ReceiverNodes: ([]string)([]string{
						(string)("BMS"),
					}),
					DefaultValue: (int)(0),
				}),
				(*descriptor.Signal)(&descriptor.Signal{
					Name:              (string)("LowThermID"),
					Start:             (uint8)(56),
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
			SenderNode: (string)("TMS"),
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
			Name:        (string)("CANmod"),
			Description: (string)(""),
		}),
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("FC"),
			Description: (string)(""),
		}),
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("PC_SG"),
			Description: (string)(""),
		}),
		(*descriptor.Node)(&descriptor.Node{
			Name:        (string)("TMS"),
			Description: (string)(""),
		}),
	}),
})
