package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the top-level portwatch configuration.
type Config struct {
	Hosts   []HostConfig  `json:"hosts"`
	Interval time.Duration `json:"-"`
	// IntervalSeconds is used for JSON marshaling.
	IntervalSeconds int    `json:"interval_seconds"`
	StorePath       string `json:"store_path"`
}

// HostConfig describes a single host and the ports to scan.
type HostConfig struct {
	Address   string `json:"address"`
	PortRange string `json:"port_range"`
	Label     string `json:"label,omitempty"`
}

// Load reads a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	if cfg.IntervalSeconds <= 0 {
		cfg.IntervalSeconds = 60
	}
	cfg.Interval = time.Duration(cfg.IntervalSeconds) * time.Second

	if cfg.StorePath == "" {
		cfg.StorePath = "portwatch.db"
	}

	return &cfg, nil
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Hosts:           []HostConfig{},
		IntervalSeconds: 60,
		Interval:        60 * time.Second,
		StorePath:       "portwatch.db",
	}
}
