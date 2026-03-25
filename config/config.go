package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type AlertConfig struct {
	Threshold int `yaml:"threshold"`
	Duration  int `yaml:"duration"`
}

type AlertsConfig struct {
	CPU        AlertConfig `yaml:"cpu"`
	RAM        AlertConfig `yaml:"ram"`
	ProcessCPU AlertConfig `yaml:"process_cpu"`
}

type Config struct {
	ServerKey        string       `yaml:"server_key"`
	ProjectID        string       `yaml:"project_id"`
	ServerIdentifier string       `yaml:"server_identifier"`
	SamplingInterval int          `yaml:"sampling_interval"`
	Alerts           AlertsConfig `yaml:"alerts"`
	APIEndpoint      string       `yaml:"api_endpoint"`
	Registered       bool         `yaml:"registered"`
}

var defaultConfig = Config{
	ServerKey:        "",
	ProjectID:        "",
	ServerIdentifier: "",
	SamplingInterval: 5,
	APIEndpoint:      "https://watchup.space",
	Registered:       false,
	Alerts: AlertsConfig{
		CPU: AlertConfig{
			Threshold: 80,
			Duration:  300,
		},
		RAM: AlertConfig{
			Threshold: 75,
			Duration:  600,
		},
		ProcessCPU: AlertConfig{
			Threshold: 60,
			Duration:  120,
		},
	},
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Config file not found at %s, using defaults\n", path)
			return &defaultConfig, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Apply defaults for missing values
	if cfg.SamplingInterval == 0 {
		cfg.SamplingInterval = defaultConfig.SamplingInterval
	}
	if cfg.APIEndpoint == "" {
		cfg.APIEndpoint = defaultConfig.APIEndpoint
	}

	return &cfg, nil
}

func (c *Config) GetSamplingDuration() time.Duration {
	return time.Duration(c.SamplingInterval) * time.Second
}


func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (c *Config) IsRegistered() bool {
	return c.Registered && c.ServerKey != "" && c.ProjectID != ""
}
