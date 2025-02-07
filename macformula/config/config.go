package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/pkg/errors"
)

// Config holds configuration variables for the HIL setup
type Config struct {
	Revision      string `yaml:"revision"`
	CanInterfaces struct {
		Veh string `yaml:"veh"`
		Pt  string `yaml:"pt"`
	} `yaml:"canInterfaces"`
	ResultsDir              string `yaml:"resultsDir"`
	LogsDir                 string `yaml:"logsDir"`
	TagsFilePath            string `yaml:"tagsFilePath"`
	CanTracerTimeoutMinutes int    `yaml:"canTracerTimeoutMinutes"`
	SilPort                 int    `yaml:"silPort"`
}

// NewConfig returns a new Config type
func NewConfig(path string) (*Config, error) {
	config := &Config{}

	// Open config file
	file, err := os.Open(path)
	if err != nil {
		return config, errors.Wrap(err, "failed to open config file")
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&config)
	if err != nil {
		return config, errors.Wrap(err, "failed to parse config file")
	}

	return config, nil
}
