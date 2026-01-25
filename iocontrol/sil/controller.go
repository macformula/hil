package sil

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	flatbuffers "github.com/google/flatbuffers/go"
	cobs "github.com/justincpresley/go-cobs"
	pb "github.com/macformula/hil/iocontrol/sil/gotobuf"
	signals "github.com/macformula/hil/iocontrol/sil/signals"
)

const (
	_unsetDigitalValue = false
	_unsetAnalogValue  = 0.0
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/signals.proto
type CircularBuffer struct {
	buffer     []byte
	size       int
	write_head int
	full       bool
}

func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		buffer: make([]byte, size),
		size:   size,
	}
}
func (b *CircularBuffer) Add(entry byte) {
	b.buffer[b.write_head] = entry
	b.write_head++
	if b.write_head == b.size {
		if !b.full {
			b.full = true
		}
		b.write_head = 0
	}
}
func (b *CircularBuffer) Get() (int, []byte) {
	return_buffer := make([]byte, b.size)
	len := 0
	if b.full {
		len = b.size
	} else {
		len = b.write_head
	}
	for i := range b.size {
		return_buffer[i] = b.buffer[(b.write_head+b.size-i-1)%b.size]
	}

	return len, return_buffer
}

type Controller struct {
	l        *zap.Logger
	port     int
	listener net.Listener
	Pins     *PinModel
	inputs   *pb.Input
	outputs  *pb.Output
	//Inputs   lvcontroller_inputs
	//Outputs  lvcontroller_outputs // find what these should be called
}

// NewController returns a new SIL Controller.
func NewController(port int, l *zap.Logger, digitalInputs []*DigitalPin, digitalOutputs []*DigitalPin, analogInputs []*AnalogPin, analogOutputs []*AnalogPin) *Controller {
	return &Controller{
		l:       l,
		port:    port,
		Pins:    NewPinModel(l, digitalInputs, digitalOutputs, analogInputs, analogOutputs),
		inputs:  &pb.Input{},
		outputs: &pb.Output{},
	}
}

func (c *Controller) Open(ctx context.Context) error {

	c.l.Info("opening sil FbController")

	addr := fmt.Sprintf("localhost:%v", c.port)

	listener, err := net.Listen("tcp", addr)
	c.listener = listener
	if err != nil {
		c.l.Error(fmt.Sprintf("creating listener: %s", errors.Wrap(err, "creating listener")))
		return errors.Wrap(err, "creating sil listener")
	}

	c.l.Info(fmt.Sprintf("sil listening on %s", addr))

	for {
		conn, err := c.listener.Accept()
		if err != nil {
			c.l.Fatal(fmt.Sprintf("accepting sil client: %s", err))
		}
		go c.handleConnection(conn)
	}
}

func (c *Controller) Close() error {
	c.l.Info("closing sil FbController")
	c.listener.Close()
	return nil
}

func (c *Controller) handleConnection(conn net.Conn) {
	WriteEncoder, err := cobs.NewEncoder(cobs.Config{
		SpecialByte: 0x00,
		Delimiter:   true,
		EndingSave:  true,
		Type:        cobs.Native,
	})
	if err != nil {
		c.l.Info(fmt.Sprintf("Error in creating the write encoder: %v", err))
	}
	ReadEncoder, err := cobs.NewEncoder(cobs.Config{
		SpecialByte: 0x00,
		Delimiter:   true,
		EndingSave:  true,
		Type:        cobs.ReversedNative,
	})
	if err != nil {
		c.l.Info(fmt.Sprintf("Error in creating the read encoder: %v", err))
	}
	go c.Writer(conn, WriteEncoder)
	go c.Reader(conn, ReadEncoder)
}

func (c *Controller) Reader(conn net.Conn, encoder cobs.Encoder) {
	defer conn.Close()
	c.l.Sync()
	buffer := NewCircularBuffer(40)
	output := &pb.Output{}

	for {
		readbuf := make([]byte, 2048)
		bytes_read, err := conn.Read(readbuf)
		if err != nil {
			if err != io.EOF {
				c.l.Error(fmt.Sprintf("read error: %s", err))
			}
		}
		c.l.Sync()
		for i := range bytes_read {
			buffer.Add(readbuf[i])
		}

		start_index := -1
		end_index := -1
		buf_len, msg_buf := buffer.Get()
		for i := range buf_len {
			if start_index != -1 && msg_buf[i] == 0 {
				end_index = i + 1
				break
			}
			if msg_buf[i] == 0 {
				start_index = i + 1
			}
		}
		if end_index != -1 {
			read_slice := msg_buf[start_index:end_index]
			cobs_decoded := encoder.Decode(read_slice)
			for i := 0; i < len(cobs_decoded)/2; i++ {
				cobs_decoded[i], cobs_decoded[len(cobs_decoded)-1-i] =
					cobs_decoded[len(cobs_decoded)-1-i], cobs_decoded[i]
			}
			err = proto.Unmarshal(cobs_decoded, output)
			if err != nil {
				c.l.Info(fmt.Sprintf("unmarshal error: %s", err))
			}
			c.l.Sync()
		}
		time.Sleep(time.Millisecond * 50)
	}
}
func (c *Controller) Writer(conn net.Conn, encoder cobs.Encoder) {
	defer conn.Close()
	//write
	for {
		c.inputs.ImdFault = !c.inputs.ImdFault
		c.inputs.BmsFault = !c.inputs.BmsFault
		proto_encoded, err := proto.Marshal(c.inputs)
		if err != nil {
			c.l.Info(fmt.Sprintf("marshalling error: %v", err))
		}
		cobs_encoded := encoder.Encode(proto_encoded)
		_, err = conn.Write(cobs_encoded)
		if err != nil {
			c.l.Info(fmt.Sprintf("write error: %s", err))
		}
		time.Sleep(time.Millisecond * 2000)
	}

}

// WriteCurrent sets the current of a SIL analog pin (unimplemented for SIL).
func (c *Controller) WriteCurrent(_ *AnalogPin, _ float64) error {
	return errors.New("unimplemented function on sil FbController")
}

// ReadCurrent returns the current of a SIL analog pin (unimplemented for SIL).
func (c *Controller) ReadCurrent(_ *AnalogPin) (float64, error) {
	return 0.00, errors.New("unimplemented function on sil FbController")
}

func deserializeReadRequest(unionTable *flatbuffers.Table) (string, string, signals.SIGNAL_TYPE, signals.SIGNAL_DIRECTION) {
	unionRequest := new(signals.ReadRequest)
	unionRequest.Init(unionTable.Bytes, unionTable.Pos)

	ecu := string(unionRequest.EcuName())
	sigName := string(unionRequest.SignalName())
	sigType := unionRequest.SignalType()
	sigDirection := unionRequest.SignalDirection()

	return ecu, sigName, sigType, sigDirection
}

func deserializeSetRequest(unionTable *flatbuffers.Table) (string, string, signals.SIGNAL_TYPE, signals.SIGNAL_DIRECTION, bool, float64) {
	unionRequest := new(signals.SetRequest)
	unionRequest.Init(unionTable.Bytes, unionTable.Pos)

	ecu := string(unionRequest.EcuName())
	sigName := string(unionRequest.SignalName())
	sigType := unionRequest.SignalType()
	sigDirection := unionRequest.SignalDirection()

	unionTable = new(flatbuffers.Table)
	if unionRequest.SignalValue(unionTable) {
		switch unionRequest.SignalValueType() {
		case signals.SignalValueDigital:
			unionSignalValue := new(signals.Digital)
			unionSignalValue.Init(unionTable.Bytes, unionTable.Pos)

			return ecu, sigName, sigType, sigDirection, unionSignalValue.Value(), _unsetAnalogValue

		case signals.SignalValueAnalog:
			unionSignalValue := new(signals.Analog)
			unionSignalValue.Init(unionTable.Bytes, unionTable.Pos)

			return ecu, sigName, sigType, sigDirection, _unsetDigitalValue, unionSignalValue.Voltage()
		}
	}
	return "", "", signals.SIGNAL_TYPEANALOG, signals.SIGNAL_DIRECTIONINPUT, _unsetDigitalValue, _unsetAnalogValue
}

func deserializeRegisterRequest(unionTable *flatbuffers.Table) (string, string, signals.SIGNAL_TYPE, signals.SIGNAL_DIRECTION) {
	unionRequest := new(signals.RegisterRequest)
	unionRequest.Init(unionTable.Bytes, unionTable.Pos)

	ecu := string(unionRequest.EcuName())
	sigName := string(unionRequest.SignalName())
	sigType := unionRequest.SignalType()
	sigDirection := unionRequest.SignalDirection()

	return ecu, sigName, sigType, sigDirection
}

func serializeReadResponse(sigType signals.SignalValue, level bool, voltage float64, ok bool, err string) []byte {
	builder := flatbuffers.NewBuilder(1024)
	errorString := builder.CreateString(err)

	var sigVal flatbuffers.UOffsetT
	if sigType == signals.SignalValueDigital {
		signals.DigitalStart(builder)
		signals.DigitalAddValue(builder, level)
		sigVal = signals.DigitalEnd(builder)

	} else if sigType == signals.SignalValueAnalog {
		signals.AnalogStart(builder)
		signals.AnalogAddVoltage(builder, voltage)
		sigVal = signals.AnalogEnd(builder)
	}

	signals.ReadResponseStart(builder)
	signals.ReadResponseAddSignalValueType(builder, sigType)
	signals.ReadResponseAddSignalValue(builder, sigVal)
	signals.ReadResponseAddOk(builder, ok)
	signals.ReadResponseAddError(builder, errorString)
	readResponse := signals.ReadResponseEnd(builder)

	signals.ResponseStart(builder)
	signals.ResponseAddResponseType(builder, signals.ResponseTypeReadResponse)
	signals.ResponseAddResponse(builder, readResponse)
	response := signals.ResponseEnd(builder)
	builder.Finish(response)

	return builder.FinishedBytes()
}

// Keep these last two functions. Depending on how firmware interacts with sil. It might not listen for a response back.
func serializeSetResponse(ok bool, err string) []byte {
	builder := flatbuffers.NewBuilder(1024)
	errorString := builder.CreateString(err)

	signals.SetResponseStart(builder)
	signals.SetResponseAddError(builder, errorString)
	signals.SetResponseAddOk(builder, ok)
	setResponse := signals.ReadRequestEnd(builder)

	signals.RequestStart(builder)
	signals.RequestAddRequestType(builder, signals.RequestTypeSetRequest)
	signals.RequestAddRequest(builder, setResponse)
	setRequest := signals.RequestEnd(builder)
	builder.Finish(setRequest)
	return builder.FinishedBytes()
}

func serializeRegisterResponse(ok bool, err string) []byte {
	builder := flatbuffers.NewBuilder(1024)
	errorString := builder.CreateString(err)

	signals.RegisterResponseStart(builder)
	signals.RegisterResponseAddError(builder, errorString)
	signals.RegisterResponseAddOk(builder, ok)
	registerResponse := signals.RegisterResponseEnd(builder)

	signals.RequestStart(builder)
	signals.RequestAddRequestType(builder, signals.RequestTypeRegisterRequest)
	signals.RequestAddRequest(builder, registerResponse)
	setRequest := signals.RequestEnd(builder)
	builder.Finish(setRequest)
	return builder.FinishedBytes()
}
