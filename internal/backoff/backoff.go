// Package backoff provides exponential backoff with jitter for retry scheduling.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Config holds parameters for exponential backoff.
type Config struct {
	// InitialInterval is the delay before the first retry.
	InitialInterval time.Duration
	// MaxInterval caps the computed delay.
	MaxInterval time.Duration
	// Multiplier is applied to the interval after each attempt.
	Multiplier float64
	// JitterFraction adds randomness in [0, JitterFraction) of the interval.
	JitterFraction float64
}

// DefaultConfig returns a sensible default backoff configuration.
func DefaultConfig() Config {
	return Config{
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		JitterFraction:  0.2,
	}
}

// Backoff computes the delay for a given attempt number (0-indexed).
// It applies exponential growth capped at MaxInterval, then adds jitter.
func (c Config) Backoff(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	base := float64(c.InitialInterval) * math.Pow(c.Multiplier, float64(attempt))
	if base > float64(c.MaxInterval) {
		base = float64(c.MaxInterval)
	}
	jitter := 0.0
	if c.JitterFraction > 0 {
		jitter = rand.Float64() * c.JitterFraction * base //nolint:gosec
	}
	return time.Duration(base + jitter)
}

// Series returns the first n backoff durations for inspection or testing.
func (c Config) Series(n int) []time.Duration {
	out := make([]time.Duration, n)
	for i := range out {
		out[i] = c.Backoff(i)
	}
	return out
}
