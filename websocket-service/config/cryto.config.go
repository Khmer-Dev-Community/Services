package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ConfigScrete struct {
	SecretKey string `yaml:"secret_key"`
}

func LoadConfigCryto(configFilePath string) (*ConfigScrete, error) {
	// Read the YAML file
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %v", err)
	}

	// Initialize a Config struct
	var config ConfigScrete

	// Unmarshal the YAML file into the Config struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	// Check if the SecretKey is loaded
	if config.SecretKey == "" {
		return nil, fmt.Errorf("secret_key not found in YAML file")
	}

	return &config, nil
}
