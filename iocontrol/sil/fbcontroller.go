package sil

import (
	"context"
	"fmt"
	"io"

	"net"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	flatbuffers "github.com/google/flatbuffers/go"
	pb "github.com/macformula/hil/iocontrol/sil/proto"
	signals "github.com/macformula/hil/iocontrol/sil/signals"
)

const (
	_unsetDigitalValue = false
	_unsetAnalogValue  = 0.0
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/signals.proto
type FbController struct {
	pb.UnimplementedSignalsServer

	l      *zap.Logger
	port   int
	server *grpc.Server
	pins   *PinModel
}

// NewFbController returns a new SIL FbController.
func NewFbController(port int, l *zap.Logger) *FbController {
	return &FbController{
		l:    l.Named(_loggerName),
		port: port,
		pins: NewPinModel(),
	}
}

// Open configures the gRPC server to communicate with firmware.
func (c *FbController) Open(ctx context.Context) error {
	c.l.Info("opening sil FbController")

	addr := fmt.Sprintf("localhost:%v", c.port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		c.l.Error(fmt.Sprintf("creating listener: %s", errors.Wrap(err, "creating listener")))
	}
	defer listener.Close()

	c.l.Info(fmt.Sprintf("sil listening on %s", addr))

	for {
		conn, err := listener.Accept()
		if err != nil {
			c.l.Fatal(fmt.Sprintf("accepting sil client: %s", err))
		}
		go c.handleConnection(conn)
	}
}

// Close stops the gRPC server.
func (c *FbController) Close() error {
	c.l.Info("closing sil FbController")
	c.server.GracefulStop()

	return nil
}

func (c *FbController) handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 2024)
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				c.l.Error(fmt.Sprintf("read error: %s", err))
			}
			break
		}
		request := signals.GetRootAsRequest(buffer, 0)

		unionTable := new(flatbuffers.Table)
		if request.Request(unionTable) {
			c.l.Info("recvieved request")
			requestType := request.RequestType()
			switch requestType {
			case signals.RequestTypeReadRequest:
				ecu, sigName, sigType, sigDirection := deserializeReadRequest(unionTable)
				c.l.Info(fmt.Sprintf("ecu: %s, signame: %s, sigType: %s, sigDir: %s", ecu, sigName, sigType, sigDirection))

				switch sigType {
				case signals.SIGNAL_TYPEDIGITAL:
					switch sigDirection {
					case signals.SIGNAL_DIRECTIONINPUT:
						c.pins.ReadDigitalInput(ecu, sigName)
					case signals.SIGNAL_DIRECTIONOUTPUT:
						c.pins.ReadDigitalOutput(ecu, sigName)
					}
				case signals.SIGNAL_TYPEANALOG:
					switch sigDirection {
					case signals.SIGNAL_DIRECTIONINPUT:
						c.pins.ReadAnalogInput(ecu, sigName)
					case signals.SIGNAL_DIRECTIONOUTPUT:
						c.pins.ReadAnalogOutput(ecu, sigName)
					}
				}
			case signals.RequestTypeSetRequest:
				ecu, sigName, sigType, sigDirection, value, voltage := deserializeSetRequest(unionTable)

				switch sigType {
				case signals.SIGNAL_TYPEDIGITAL:
					switch sigDirection {
					case signals.SIGNAL_DIRECTIONINPUT:
						c.pins.SetDigitalInput(ecu, sigName, value)
					case signals.SIGNAL_DIRECTIONOUTPUT:
						c.pins.SetDigitalOutput(ecu, sigName, value)
					}
				case signals.SIGNAL_TYPEANALOG:
					switch sigDirection {
					case signals.SIGNAL_DIRECTIONINPUT:
						c.pins.SetAnalogInput(ecu, sigName, voltage)
					case signals.SIGNAL_DIRECTIONOUTPUT:
						c.pins.SetAnalogOutput(ecu, sigName, voltage)
					}
				}

			case signals.RequestTypeRegisterRequest:
				ecu, sigName, sigType, sigDirection := deserializeRegisterRequest(unionTable)

				switch sigType {
				case signals.SIGNAL_TYPEDIGITAL:
					switch sigDirection {
					case signals.SIGNAL_DIRECTIONINPUT:
						c.pins.RegisterDigitalInput(ecu, sigName)
					case signals.SIGNAL_DIRECTIONOUTPUT:
						c.pins.RegisterDigitalOutput(ecu, sigName)
					}
				case signals.SIGNAL_TYPEANALOG:
					switch sigDirection {
					case signals.SIGNAL_DIRECTIONINPUT:
						c.pins.RegisterAnalogInput(ecu, sigName)
					case signals.SIGNAL_DIRECTIONOUTPUT:
						c.pins.RegisterAnalogOutput(ecu, sigName)
					}
				}
			}

		}
	}
}

// WriteCurrent sets the current of a SIL analog pin (unimplemented for SIL).
func (c *FbController) WriteCurrent(_ *AnalogPin, _ float64) error {
	return errors.New("unimplemented function on sil FbController")
}

// ReadCurrent returns the current of a SIL analog pin (unimplemented for SIL).
func (c *FbController) ReadCurrent(_ *AnalogPin) (float64, error) {
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
	errorString := builder.CreateString(err)
	signals.ReadResponseAddError(builder, errorString)
	read_response := signals.ReadResponseEnd(builder)

	signals.ResponseStart(builder)
	signals.ResponseAddResponseType(builder, signals.ResponseTypeReadResponse)
	signals.ResponseAddResponse(builder, read_response)
	response := signals.ResponseEnd(builder)
	builder.Finish(response)

	return builder.FinishedBytes()
}

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
