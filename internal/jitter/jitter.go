// Package jitter adds randomised delay to scheduled scans to avoid
// thundering-herd problems when many hosts are monitored simultaneously.
package jitter

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

// Jitter holds configuration for randomised scan delays.
type Jitter struct {
	mu      sync.Mutex
	max     time.Duration
	source  *rand.Rand
}

// New returns a Jitter that will delay scans by a random duration in
// [0, maxDelay). Passing zero disables jitter entirely.
func New(maxDelay time.Duration) *Jitter {
	//nolint:gosec // non-cryptographic randomness is intentional here
	return &Jitter{
		max:    maxDelay,
		source: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Wait blocks for a random duration in [0, max) and then returns nil.
// If the context is cancelled before the delay elapses, Wait returns
// ctx.Err() immediately.
func (j *Jitter) Wait(ctx context.Context) error {
	if j.max <= 0 {
		return nil
	}

	j.mu.Lock()
	delay := time.Duration(j.source.Int63n(int64(j.max)))
	j.mu.Unlock()

	if delay == 0 {
		return nil
	}

	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Sample returns a single random duration in [0, max) without blocking.
// Useful for testing or pre-computing delays.
func (j *Jitter) Sample() time.Duration {
	if j.max <= 0 {
		return 0
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	return time.Duration(j.source.Int63n(int64(j.max)))
}
