package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Owner       string `yaml:"owner"`
	Repo        string `yaml:"repo"`
	MonitorPath string `yaml:"monitorPath"`
	Script      string `yaml:"script"`
	AgeKey      string `yaml:"ageSecret"`
	GithubToken string `yaml:"githubToken"`
	Interval    int    `yaml:"interval"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}
