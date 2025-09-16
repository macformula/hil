package sil

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	signals "github.com/macformula/hil/cmd/basicio/signals"
	"google.golang.org/protobuf/proto"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/signals.proto
type Controller struct {
	l        *zap.Logger
	port     int
	listener net.Listener
	Pins     *PinModel
}

// NewController returns a new SIL Controller.
func NewController(port int, l *zap.Logger, digitalInputs []*DigitalPin, digitalOutputs []*DigitalPin, analogInputs []*AnalogPin, analogOutputs []*AnalogPin) *Controller {
	return &Controller{
		l:    l,
		port: port,
		Pins: NewPinModel(l, digitalInputs, digitalOutputs, analogInputs, analogOutputs),
	}
}

func (c *Controller) Open(ctx context.Context) error {
	c.l.Info("opening sil Controller")

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

		var request signals.Request
		err = proto.Unmarshal(buffer, &request)

		if err != nil {
			c.l.Error(fmt.Sprintf("deserialize error: %v", err))
			break
		}

		signal := request.Signal
		var response *signals.Response = nil

		switch action := request.GetRequest().(type) {
		case *signals.Request_Read:
			switch signal.Type {
			case signals.SignalType_SIGNAL_TYPE_DIGITAL:
				pin := &DigitalPin{EcuName: signal.EcuName, SigName: signal.SignalName}

				switch signal.Direction {
				case signals.SignalDirection_SIGNAL_DIRECTION_INPUT:
					level, err := c.Pins.ReadDigitalInput(pin)
					if err == nil {
						response = &signals.Response{
							Response: &signals.Response_Read{
								Read: &signals.ReadResponse{Value: &signals.SignalValue{
									Value: &signals.SignalValue_Digital{Digital: level}},
								},
							},
							Ok:    true,
							Error: "",
						}
					} else {
						err_s := fmt.Sprintf("read digital input (%s) error: %s", pin_str(pin), err)
						c.l.Error(err_s)

						response = &signals.Response{
							Response: nil,
							Ok:       false,
							Error:    err_s,
						}
					}

				case signals.SignalDirection_SIGNAL_DIRECTION_OUTPUT:
					level, err := c.Pins.ReadDigitalOutput(pin)
					if err == nil {
						response = &signals.Response{
							Response: &signals.Response_Read{
								Read: &signals.ReadResponse{Value: &signals.SignalValue{
									Value: &signals.SignalValue_Digital{Digital: level}},
								},
							},
							Ok:    true,
							Error: "",
						}
					} else {
						err_s := fmt.Sprintf("read digital output %s", pin_str(pin))
						c.l.Error(err_s)

						response = &signals.Response{
							Response: nil,
							Ok:       false,
							Error:    err_s,
						}
					}
				}

			case signals.SIGNAL_TYPEANALOG:
				switch sigDirection {
				case signals.SIGNAL_DIRECTIONINPUT:
					pin := NewAnalogInputPin(ecu, sigName)
					voltage, err := c.Pins.ReadAnalogInput(pin)
					if err != nil {
						c.l.Error(fmt.Sprintf("read analog input ecu (%s) signal name (%s)", ecu, sigName))

						response := serializeReadResponse(signals.SignalValueAnalog, _unsetDigitalValue, _unsetAnalogValue, false, fmt.Sprintf("read digital output ecu (%s) signal name (%s)", ecu, sigName))
						_, err = conn.Write(response)
						if err != nil {
							c.l.Error(fmt.Sprintf("write sil response (%s)", err.Error()))
						}
					}

					response := serializeReadResponse(signals.SignalValueDigital, _unsetDigitalValue, voltage, true, "")
					_, err = conn.Write(response)
					if err != nil {
						c.l.Error(fmt.Sprintf("write sil response (%s)", err.Error()))
					}
				case signals.SIGNAL_DIRECTIONOUTPUT:
					pin := NewAnalogOutputPin(ecu, sigName)
					voltage, err := c.Pins.ReadAnalogOutput(pin)
					if err != nil {
						c.l.Error(fmt.Sprintf("read analog output ecu (%s) signal name (%s)", ecu, sigName))

						response := serializeReadResponse(signals.SignalValueDigital, _unsetDigitalValue, _unsetAnalogValue, false, fmt.Sprintf("read analog output ecu (%s) signal name (%s)", ecu, sigName))
						_, err = conn.Write(response)
						if err != nil {
							c.l.Error(fmt.Sprintf("write sil response (%s)", err.Error()))
						}
					}

					response := serializeReadResponse(signals.SignalValueDigital, _unsetDigitalValue, voltage, true, "")
					_, err = conn.Write(response)
					if err != nil {
						c.l.Error(fmt.Sprintf("write sil response (%s)", err.Error()))
					}
				}
			}
		case *signals.Request_Set:

			ecu, sigName, sigType, sigDirection, value, voltage := deserializeSetRequest(unionTable)

			switch sigType {
			case signals.SIGNAL_TYPEDIGITAL:
				switch sigDirection {
				case signals.SIGNAL_DIRECTIONINPUT:
					pin := NewDigitalInputPin(ecu, sigName)
					c.Pins.SetDigitalInput(pin, value)
				case signals.SIGNAL_DIRECTIONOUTPUT:
					pin := NewDigitalOutputPin(ecu, sigName)
					c.Pins.SetDigitalOutput(pin, value)
				}
			case signals.SIGNAL_TYPEANALOG:
				switch sigDirection {
				case signals.SIGNAL_DIRECTIONINPUT:
					pin := NewAnalogInputPin(ecu, sigName)
					c.Pins.SetAnalogInput(pin, voltage)
				case signals.SIGNAL_DIRECTIONOUTPUT:
					pin := NewAnalogOutputPin(ecu, sigName)
					c.Pins.SetAnalogOutput(pin, voltage)
				}
			}

		case *signals.Request_Register:

			ecu, sigName, sigType, sigDirection := deserializeRegisterRequest(unionTable)

			switch sigType {
			case signals.SIGNAL_TYPEDIGITAL:
				switch sigDirection {
				case signals.SIGNAL_DIRECTIONINPUT:
					pin := NewDigitalInputPin(ecu, sigName)
					c.Pins.RegisterDigitalInput(pin)
				case signals.SIGNAL_DIRECTIONOUTPUT:
					pin := NewDigitalOutputPin(ecu, sigName)
					c.Pins.RegisterDigitalOutput(pin)
				}
			case signals.SIGNAL_TYPEANALOG:
				switch sigDirection {
				case signals.SIGNAL_DIRECTIONINPUT:
					pin := NewAnalogInputPin(ecu, sigName)
					c.Pins.RegisterAnalogInput(pin)
				case signals.SIGNAL_DIRECTIONOUTPUT:
					pin := NewAnalogOutputPin(ecu, sigName)
					c.Pins.RegisterAnalogOutput(pin)
				}
			}
		}

		bytes, err := proto.Marshal(response)
		if err != nil {
			c.l.Error(fmt.Sprintf("Failed to serialize response (%s)", err.Error()))
		}

		_, err = conn.Write(bytes)
		if err != nil {
			c.l.Error(fmt.Sprintf("write sil response (%s)", err.Error()))
		}
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

func pin_str(pin *DigitalPin) string {
	return fmt.Sprintf("ecu (%s) signal name (%s)", pin.EcuName, pin.SigName)
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
