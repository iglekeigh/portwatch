package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all portwatch configuration.
type Config struct {
	Hosts    []string      `yaml:"hosts"`
	Ports    string        `yaml:"ports"`
	Interval time.Duration `yaml:"interval"`
	StorePath string       `yaml:"store_path"`
	Timeout  time.Duration `yaml:"timeout"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{
		Hosts:     []string{"127.0.0.1"},
		Ports:     "1-1024",
		Interval:  60 * time.Second,
		StorePath: "/tmp/portwatch_state.json",
		Timeout:   500 * time.Millisecond,
	}
}

// Load reads a YAML config file, falling back to defaults for missing fields.
func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	if cfg.Interval == 0 {
		cfg.Interval = Default().Interval
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = Default().Timeout
	}
	if cfg.StorePath == "" {
		cfg.StorePath = Default().StorePath
	}
	return cfg, nil
}
