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
		logger.Info(fmt.Sprintf("digital input pin: ecu (%s) sig name (%s)", digitalInputPin.EcuName, digitalInputPin.SigName))
		mapSet(digitalInputPins, digitalInputPin.EcuName, digitalInputPin.SigName, _defaultDigitalValue)
	}
	for _, digitalOutputPin := range digitalOutputs {
		logger.Info(fmt.Sprintf("dig output pin: ecu (%s) sig name (%s)", digitalOutputPin.EcuName, digitalOutputPin.SigName))
		mapSet(digitalOutputPins, digitalOutputPin.EcuName, digitalOutputPin.SigName, _defaultDigitalValue)
	}
	for _, analogInput := range analogInputs {
		mapSet(analogInputPins, analogInput.EcuName, analogInput.SigName, _defaultAnalogValue)
	}
	for _, analogOutput := range analogOutputs {
		mapSet(analogOutputPins, analogOutput.EcuName, analogOutput.SigName, _defaultAnalogValue)
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

func (p *PinModel) RegisterDigitalOutput(pin *DigitalPin) {
	mapSet(p.digitalOutputPins, pin.EcuName, pin.SigName, false)
}

func (p *PinModel) RegisterDigitalInput(pin *DigitalPin) {
	mapSet(p.digitalInputPins, pin.EcuName, pin.SigName, false)
}

func (p *PinModel) RegisterAnalogOutput(pin *AnalogPin) {
	mapSet(p.analogOutputPins, pin.EcuName, pin.SigName, 0.0)
}

func (p *PinModel) RegisterAnalogInput(pin *AnalogPin) {
	mapSet(p.analogInputPins, pin.EcuName, pin.SigName, 0.0)
}

// ReadDigital returns the level of a SIL digital pin.
func (p *PinModel) ReadDigitalInput(pin *DigitalPin) (bool, error) {
	p.digitalInputMtx.Lock()
	defer p.digitalInputMtx.Unlock()

	level, ok := mapLookup(p.digitalInputPins, pin.EcuName, pin.SigName)
	if !ok {
		return false, errors.Errorf("no entry for ecu name (%s) signal name (%s)",
			pin.EcuName, pin.SigName)
	}

	return level, nil
}

// ReadDigital returns the level of a SIL digital pin.
func (p *PinModel) ReadAnalogInput(pin *AnalogPin) (float64, error) {
	p.analogInputMtx.Lock()
	defer p.analogInputMtx.Unlock()

	voltage, ok := mapLookup(p.analogInputPins, pin.EcuName, pin.SigName)
	if !ok {
		return 0, errors.Errorf("no entry for ecu name (%s) signal name (%s)",
			pin.EcuName, pin.SigName)
	}

	return voltage, nil
}

// ReadDigital returns the level of a SIL digital pin.
func (p *PinModel) ReadDigitalOutput(pin *DigitalPin) (bool, error) {
	p.digitalOutputMtx.Lock()
	defer p.digitalOutputMtx.Unlock()

	level, ok := mapLookup(p.digitalOutputPins, pin.EcuName, pin.SigName)
	if !ok {
		return false, errors.Errorf("no entry for ecu name (%s) signal name (%s)",
			pin.EcuName, pin.SigName)
	}

	return level, nil
}

// ReadDigital returns the level of a SIL digital pin.
func (p *PinModel) ReadAnalogOutput(pin *AnalogPin) (float64, error) {
	p.analogOutputMtx.Lock()
	defer p.analogOutputMtx.Unlock()

	voltage, ok := mapLookup(p.analogOutputPins, pin.EcuName, pin.SigName)
	if !ok {
		return 0, errors.Errorf("no entry for ecu name (%s) signal name (%s)",
			pin.EcuName, pin.SigName)
	}

	return voltage, nil
}

// SetDigital sets an output digital pin for a SIL digital pin.
func (p *PinModel) SetDigitalInput(pin *DigitalPin, level bool) error {
	p.digitalInputMtx.Lock()
	defer p.digitalInputMtx.Unlock()

	_, ok := mapLookup(p.digitalOutputPins, pin.EcuName, pin.SigName)
	if !ok {
		p.RegisterDigitalInput(pin)
	}

	p.digitalInputPins[pin.EcuName][pin.SigName] = level

	return nil
}

func (p *PinModel) SetAnalogInput(pin *AnalogPin, voltage float64) error {
	p.analogInputMtx.Lock()
	defer p.analogInputMtx.Unlock()

	_, ok := mapLookup(p.analogInputPins, pin.EcuName, pin.SigName)
	if !ok {
		p.RegisterAnalogInput(pin)
	}

	p.analogInputPins[pin.EcuName][pin.SigName] = voltage

	return nil
}

// SetDigital sets an output digital pin for a SIL digital pin.
func (p *PinModel) SetDigitalOutput(pin *DigitalPin, level bool) error {
	p.digitalOutputMtx.Lock()
	defer p.digitalOutputMtx.Unlock()

	_, ok := mapLookup(p.digitalOutputPins, pin.EcuName, pin.SigName)
	if !ok {
		p.RegisterDigitalOutput(pin)
	}

	p.digitalOutputPins[pin.EcuName][pin.SigName] = level

	return nil
}

// SetDigital sets an output digital pin for a SIL digital pin.
func (p *PinModel) SetAnalogOutput(pin *AnalogPin, voltage float64) error {
	p.analogOutputMtx.Lock()
	defer p.analogOutputMtx.Unlock()

	_, ok := mapLookup(p.analogOutputPins, pin.EcuName, pin.SigName)
	if !ok {
		p.RegisterAnalogOutput(pin)
	}

	p.analogOutputPins[pin.EcuName][pin.SigName] = voltage

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
