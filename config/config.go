package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

// Config holds configuration variables for the HIL setup
type Config struct {
	TracerDirectory string `yaml:"tracerDirectory"`
	CanInterface    string `yaml:"canInterface"`
	BusName         string `yaml:"busName"`
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
