package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Exchanges []ExchangeConfig `yaml:"exchanges"`
}

type ExchangeConfig struct {
	Name    string   `yaml:"name"`
	URL     string   `yaml:"url"`
	Symbols []string `yaml:"symbols"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(f, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
