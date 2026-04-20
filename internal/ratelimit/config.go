package ratelimit

import "time"

// Config holds configuration for the rate limiter.
type Config struct {
	// Cooldown is the minimum duration between notifications for the same host.
	Cooldown time.Duration

	// MaxPerHour is the maximum number of notifications allowed per host per hour.
	// A value of 0 means no limit.
	MaxPerHour int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Cooldown:   5 * time.Minute,
		MaxPerHour: 12,
	}
}

// Validate checks that the config values are valid and applies
// defaults where zero values are provided.
func (c *Config) Validate() error {
	if c.Cooldown < 0 {
		c.Cooldown = DefaultConfig().Cooldown
	}
	if c.MaxPerHour < 0 {
		c.MaxPerHour = 0
	}
	return nil
}

// NewFromConfig creates a rate limiter from the provided Config.
func NewFromConfig(cfg Config) *RateLimiter {
	if cfg.Cooldown <= 0 {
		cfg.Cooldown = DefaultConfig().Cooldown
	}
	return New(cfg.Cooldown)
}
