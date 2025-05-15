package sil

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_defaultDigitalValue = false
	_defaultAnalogValue  = 0.0
)

type PinModel struct {
	l *zap.Logger

	digitalInputPins  map[string]map[string]bool
	digitalOutputPins map[string]map[string]bool
	analogInputPins   map[string]map[string]float64
	analogOutputPins  map[string]map[string]float64

	digitalInputMtx  *sync.Mutex
	digitalOutputMtx *sync.Mutex
	analogInputMtx   *sync.Mutex
	analogOutputMtx  *sync.Mutex
}

func NewPinModel(logger *zap.Logger, digitalInputs []*DigitalPin, digitalOutputs []*DigitalPin, analogInputs []*AnalogPin, analogOutputs []*AnalogPin) *PinModel {
	digitalInputPins := make(map[string]map[string]bool)
	digitalOutputPins := make(map[string]map[string]bool)
	analogInputPins := make(map[string]map[string]float64)
	analogOutputPins := make(map[string]map[string]float64)
	for _, digitalInputPin := range digitalInputs {
		logger.Info(fmt.Sprintf("digital input pin: ecu (%s) sig name (%s)", digitalInputPin.GetEcuName(), digitalInputPin.GetSigName()))
		mapSet(digitalInputPins, digitalInputPin.GetEcuName(), digitalInputPin.GetSigName(), _defaultDigitalValue)
	}
	for _, digitalOutputPin := range digitalOutputs {
		logger.Info(fmt.Sprintf("dig output pin: ecu (%s) sig name (%s)", digitalOutputPin.GetEcuName(), digitalOutputPin.GetSigName()))
		mapSet(digitalOutputPins, digitalOutputPin.GetEcuName(), digitalOutputPin.GetSigName(), _defaultDigitalValue)
	}
	for _, analogInput := range analogInputs {
		mapSet(analogInputPins, analogInput.GetEcuName(), analogInput.GetSigName(), _defaultAnalogValue)
	}
	for _, analogOutput := range analogOutputs {
		mapSet(analogOutputPins, analogOutput.GetEcuName(), analogOutput.GetSigName(), _defaultAnalogValue)
	}
	return &PinModel{
		digitalInputPins:  digitalInputPins,
		digitalOutputPins: digitalOutputPins,
		analogInputPins:   analogInputPins,
		analogOutputPins:  analogOutputPins,
		digitalInputMtx:   &sync.Mutex{},
		digitalOutputMtx:  &sync.Mutex{},
		analogInputMtx:    &sync.Mutex{},
		analogOutputMtx:   &sync.Mutex{},
	}
}

func (p *PinModel) RegisterDigitalOutput(ecuName, sigName string) {
	mapSet(p.digitalOutputPins, ecuName, sigName, false)
}

func (p *PinModel) RegisterDigitalInput(ecuName, sigName string) {
	mapSet(p.digitalInputPins, ecuName, sigName, false)
}

func (p *PinModel) RegisterAnalogOutput(ecuName, sigName string) {
	mapSet(p.analogOutputPins, ecuName, sigName, 0.0)
}

func (p *PinModel) RegisterAnalogInput(ecuName, sigName string) {
	mapSet(p.analogInputPins, ecuName, sigName, 0.0)
}

// ReadDigital returns the level of a SIL digital pin.
func (p *PinModel) ReadDigitalInput(ecu_name string, sig_name string) (bool, error) {
	p.digitalInputMtx.Lock()
	defer p.digitalInputMtx.Unlock()

	level, ok := mapLookup(p.digitalInputPins, ecu_name, sig_name)
	if !ok {
		return false, errors.Errorf("no entry for ecu name (%s) signal name (%s)",
			ecu_name, sig_name)
	}

	return level, nil
}

// ReadDigital returns the level of a SIL digital pin.
func (p *PinModel) ReadAnalogInput(ecu_name string, sig_name string) (float64, error) {
	p.analogInputMtx.Lock()
	defer p.analogInputMtx.Unlock()

	voltage, ok := mapLookup(p.analogInputPins, ecu_name, sig_name)
	if !ok {
		return 0, errors.Errorf("no entry for ecu name (%s) signal name (%s)",
			ecu_name, sig_name)
	}

	return voltage, nil
}

// ReadDigital returns the level of a SIL digital pin.
func (p *PinModel) ReadDigitalOutput(ecu_name string, sig_name string) (bool, error) {
	p.digitalOutputMtx.Lock()
	defer p.digitalOutputMtx.Unlock()

	level, ok := mapLookup(p.digitalOutputPins, ecu_name, sig_name)
	if !ok {
		return false, errors.Errorf("no entry for ecu name (%s) signal name (%s)",
			ecu_name, sig_name)
	}

	return level, nil
}

// ReadDigital returns the level of a SIL digital pin.
func (p *PinModel) ReadAnalogOutput(ecu_name string, sig_name string) (float64, error) {
	p.analogOutputMtx.Lock()
	defer p.analogOutputMtx.Unlock()

	voltage, ok := mapLookup(p.analogOutputPins, ecu_name, sig_name)
	if !ok {
		return 0, errors.Errorf("no entry for ecu name (%s) signal name (%s)",
			ecu_name, sig_name)
	}

	return voltage, nil
}

// SetDigital sets an output digital pin for a SIL digital pin.
func (p *PinModel) SetDigitalInput(ecu_name string, sig_name string, level bool) error {
	p.digitalInputMtx.Lock()
	defer p.digitalInputMtx.Unlock()

	_, ok := mapLookup(p.digitalOutputPins, ecu_name, sig_name)
	if !ok {
		p.RegisterDigitalInput(ecu_name, sig_name)
	}

	p.digitalInputPins[ecu_name][sig_name] = level

	return nil
}

func (p *PinModel) SetAnalogInput(ecu_name string, sig_name string, voltage float64) error {
	p.analogInputMtx.Lock()
	defer p.analogInputMtx.Unlock()

	_, ok := mapLookup(p.analogInputPins, ecu_name, sig_name)
	if !ok {
		p.RegisterAnalogInput(ecu_name, sig_name)
	}

	p.analogInputPins[ecu_name][sig_name] = voltage

	return nil
}

// SetDigital sets an output digital pin for a SIL digital pin.
func (p *PinModel) SetDigitalOutput(ecu_name string, sig_name string, level bool) error {
	p.digitalOutputMtx.Lock()
	defer p.digitalOutputMtx.Unlock()

	_, ok := mapLookup(p.digitalOutputPins, ecu_name, sig_name)
	if !ok {
		p.RegisterDigitalOutput(ecu_name, sig_name)
	}

	p.digitalOutputPins[ecu_name][sig_name] = level

	return nil
}

// SetDigital sets an output digital pin for a SIL digital pin.
func (p *PinModel) SetAnalogOutput(ecu_name string, sig_name string, voltage float64) error {
	p.analogOutputMtx.Lock()
	defer p.analogOutputMtx.Unlock()

	_, ok := mapLookup(p.analogOutputPins, ecu_name, sig_name)
	if !ok {
		p.RegisterDigitalOutput(ecu_name, sig_name)
	}

	p.analogOutputPins[ecu_name][sig_name] = voltage

	return nil
}

func mapSet[T any](m map[string]map[string]T, first, second string, value T) {
	if m[first] == nil {
		m[first] = make(map[string]T)
	}

	m[first][second] = value
}

func mapLookup[T any](m map[string]map[string]T, first, second string) (T, bool) {
	var ret T

	m1, ok := m[first]
	if !ok {
		return ret, false
	}

	ret, ok = m1[second]
	if !ok {
		return ret, false
	}

	return ret, true
}
