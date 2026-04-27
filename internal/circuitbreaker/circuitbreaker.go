// Package circuitbreaker provides a simple circuit breaker for wrapping
// notifier calls so that a repeatedly failing backend is temporarily disabled.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit breaker is in the open state and
// the call is rejected without attempting the underlying operation.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // failing; requests rejected
	StateHalfOpen              // probe request allowed
)

// Config holds tunable parameters for the circuit breaker.
type Config struct {
	// FailureThreshold is the number of consecutive failures before opening.
	FailureThreshold int
	// SuccessThreshold is the number of consecutive successes in half-open
	// state required to close the circuit again.
	SuccessThreshold int
	// OpenDuration is how long the circuit stays open before moving to half-open.
	OpenDuration time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		OpenDuration:     30 * time.Second,
	}
}

// Breaker is a thread-safe circuit breaker.
type Breaker struct {
	cfg            Config
	mu             sync.Mutex
	state          State
	consecFails    int
	consecSuccesses int
	openedAt       time.Time
}

// New creates a new Breaker with the given configuration.
func New(cfg Config) *Breaker {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = DefaultConfig().FailureThreshold
	}
	if cfg.SuccessThreshold <= 0 {
		cfg.SuccessThreshold = DefaultConfig().SuccessThreshold
	}
	if cfg.OpenDuration <= 0 {
		cfg.OpenDuration = DefaultConfig().OpenDuration
	}
	return &Breaker{cfg: cfg}
}

// Allow reports whether a call should be attempted. It returns ErrOpen when
// the circuit is open and the open duration has not yet elapsed.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return nil
	case StateOpen:
		if time.Since(b.openedAt) >= b.cfg.OpenDuration {
			b.state = StateHalfOpen
			b.consecSuccesses = 0
			return nil
		}
		return ErrOpen
	case StateHalfOpen:
		return nil
	}
	return nil
}

// RecordSuccess records a successful call and potentially closes the circuit.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.consecFails = 0
	if b.state == StateHalfOpen {
		b.consecSuccesses++
		if b.consecSuccesses >= b.cfg.SuccessThreshold {
			b.state = StateClosed
		}
	}
}

// RecordFailure records a failed call and potentially opens the circuit.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.consecSuccesses = 0
	b.consecFails++
	if b.state == StateHalfOpen || b.consecFails >= b.cfg.FailureThreshold {
		b.state = StateOpen
		b.openedAt = time.Now()
		b.consecFails = 0
	}
}

// State returns the current state of the circuit breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
