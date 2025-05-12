package main

import (
	flatbuffers "github.com/google/flatbuffers/go"
	signals "github.com/macformula/hil/cmd/basicio/signals"
)

func serializeRegisterRequest(ecu_name string, signal_name string, signalType signals.SIGNAL_TYPE, signal_direction signals.SIGNAL_DIRECTION) []byte {
	builder := flatbuffers.NewBuilder(1024)
	ecu2 := builder.CreateString(ecu_name)
	signal_name2 := builder.CreateString(signal_name)

	signals.RegisterRequestStart(builder)
	signals.RegisterRequestAddEcuName(builder, ecu2)
	signals.RegisterRequestAddSignalName(builder, signal_name2)
	signals.RegisterRequestAddSignalType(builder, signalType)
	signals.RegisterRequestAddSignalDirection(builder, signal_direction)

	registerRequest := signals.RegisterRequestEnd(builder)

	signals.RequestStart(builder)
	signals.RequestAddRequestType(builder, signals.RequestTypeRegisterRequest)
	signals.RequestAddRequest(builder, registerRequest)
	request := signals.RequestEnd(builder)
	builder.Finish(request)
	return builder.FinishedBytes()
}

func serializeReadRequest(ecu_name string, signal_name string, signalType signals.SIGNAL_TYPE, sigDirection signals.SIGNAL_DIRECTION) []byte {
	builder := flatbuffers.NewBuilder(1024)
	ecu := builder.CreateString(ecu_name)
	sig_name := builder.CreateString(signal_name)

	signals.ReadRequestStart(builder)
	signals.ReadRequestAddEcuName(builder, ecu)
	signals.ReadRequestAddSignalName(builder, sig_name)
	signals.ReadRequestAddSignalType(builder, signalType)
	signals.ReadRequestAddSignalDirection(builder, sigDirection)
	readRequest := signals.ReadRequestEnd(builder)

	signals.RequestStart(builder)
	signals.RequestAddRequestType(builder, signals.RequestTypeReadRequest)
	signals.RequestAddRequest(builder, readRequest)
	request := signals.RequestEnd(builder)
	builder.Finish(request)
	return builder.FinishedBytes()
}

func serializeSetRequest(ecu_name string, signal_name string, signalType signals.SIGNAL_TYPE, sigDirection signals.SIGNAL_DIRECTION, voltage float64, level bool) []byte {
	builder := flatbuffers.NewBuilder(1024)
	ecu2 := builder.CreateString(ecu_name)
	signal_name2 := builder.CreateString(signal_name)

	switch signalType {
	case signals.SIGNAL_TYPEDIGITAL:
		signals.DigitalStart(builder)
		signals.DigitalAddValue(builder, level)
		dig_sig := signals.DigitalEnd(builder)

		signals.SetRequestStart(builder)
		signals.SetRequestAddEcuName(builder, ecu2)
		signals.SetRequestAddSignalName(builder, signal_name2)
		signals.SetRequestAddSignalType(builder, signals.SIGNAL_TYPEDIGITAL)
		signals.SetRequestAddSignalDirection(builder, sigDirection)
		signals.SetRequestAddSignalValueType(builder, signals.SignalValueDigital)
		signals.SetRequestAddSignalValue(builder, dig_sig)
		setRequest2 := signals.ReadRequestEnd(builder)

		signals.RequestStart(builder)
		signals.RequestAddRequestType(builder, signals.RequestTypeSetRequest)
		signals.RequestAddRequest(builder, setRequest2)
		setRequest := signals.RequestEnd(builder)
		builder.Finish(setRequest)
		return builder.FinishedBytes()
	case signals.SIGNAL_TYPEANALOG:
		signals.AnalogStart(builder)
		signals.AnalogAddVoltage(builder, voltage)
		dig_sig := signals.AnalogEnd(builder)

		signals.SetRequestStart(builder)
		signals.SetRequestAddEcuName(builder, ecu2)
		signals.SetRequestAddSignalName(builder, signal_name2)
		signals.SetRequestAddSignalType(builder, signals.SIGNAL_TYPEANALOG)
		signals.SetRequestAddSignalValueType(builder, signals.SignalValueAnalog)
		signals.SetRequestAddSignalValue(builder, dig_sig)
		setRequest2 := signals.ReadRequestEnd(builder)

		signals.RequestStart(builder)
		signals.RequestAddRequestType(builder, signals.RequestTypeSetRequest)
		signals.RequestAddRequest(builder, setRequest2)
		setRequest := signals.RequestEnd(builder)
		builder.Finish(setRequest)
		return builder.FinishedBytes()
	}
	return nil
}

func deserializeReadResponse(unionTable *flatbuffers.Table) (bool, string, bool, float64) {
	unionResponse := new(signals.ReadResponse)
	unionResponse.Init(unionTable.Bytes, unionTable.Pos)

	ok := unionResponse.Ok()
	errorString := string(unionResponse.Error())

	unionTable = new(flatbuffers.Table)
	if unionResponse.SignalValue(unionTable) {
		switch unionResponse.SignalValueType() {
		case signals.SignalValueDigital:
			unionSignalValue := new(signals.Digital)
			unionSignalValue.Init(unionTable.Bytes, unionTable.Pos)

			return ok, errorString, unionSignalValue.Value(), 0.0

		case signals.SignalValueAnalog:
			unionSignalValue := new(signals.Analog)
			unionSignalValue.Init(unionTable.Bytes, unionTable.Pos)

			return ok, errorString, false, unionSignalValue.Voltage()
		}
	}

	return false, "", false, 0.0
}
