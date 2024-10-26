package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config holds the application configuration
type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
		Mode string `yaml:"mode"` // Ajout du champ Mode
	} `yaml:"server"`
	Logging struct {
		Level     string `yaml:"level"`
		Directory string `yaml:"directory"`
	} `yaml:"logging"`
	Storage struct {
		TempDirectory string `yaml:"temp_directory"`
	} `yaml:"storage"`
}

// LoadConfig reads the config file and returns a Config struct
func LoadConfig(filename string) (*Config, error) {
	config := &Config{}

	// Read the config file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse the YAML
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config YAML: %w", err)
	}

	return config, nil
}
