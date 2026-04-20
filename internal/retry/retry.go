// Package retry provides configurable retry logic for notifier operations.
package retry

import (
	"context"
	"errors"
	"time"
)

// Config holds retry behaviour settings.
type Config struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// InitialDelay is the wait time before the second attempt.
	InitialDelay time.Duration
	// Multiplier scales the delay after each failure (exponential back-off).
	Multiplier float64
}

// DefaultConfig returns a sensible default retry configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts:  3,
		InitialDelay: 500 * time.Millisecond,
		Multiplier:   2.0,
	}
}

// Do executes fn up to cfg.MaxAttempts times, backing off between failures.
// It returns the last error if all attempts are exhausted.
// The context is checked before every attempt; a cancelled context aborts immediately.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}
	if cfg.Multiplier <= 0 {
		cfg.Multiplier = 1
	}

	delay := cfg.InitialDelay
	var lastErr error

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		// No sleep after the final attempt.
		if attempt < cfg.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * cfg.Multiplier)
		}
	}

	return errors.New("all retry attempts failed: " + lastErr.Error())
}
