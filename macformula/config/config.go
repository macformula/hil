package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/pkg/errors"
)

// Config holds configuration variables for the HIL setup
type Config struct {
	Revision                    string `yaml:"revision"`
	VehCanInterface             string `yaml:"vehCanInterface"`
	PtCanInterface              string `yaml:"ptCanInterface"`
	TraceDir                    string `yaml:"traceDir"`
	LogsDir                     string `yaml:"logsDir"`
	ResultProcessorAutoStart    bool   `yaml:"resultProcessorAutoStart"`
	ResultProcessorAddr         string `yaml:"resultProcessorAddr"`
	ResultProcessorPath         string `yaml:"resultProcessorPath"`
	ResultProcessorPushToGithub bool   `yaml:"resultProcessorPushToGithub"`
	CanTracerTimeoutMinutes     int    `yaml:"canTracerTimeoutMinutes"`
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
