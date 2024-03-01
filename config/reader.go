package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func Read(file string) (*Configuration, error) {
	yuml, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var cfg *Configuration
	err = yaml.Unmarshal(yuml, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
