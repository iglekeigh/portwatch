// Package ratelimit provides a simple token-bucket rate limiter
// to prevent alert flooding when ports change frequently.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls how frequently alerts can be sent per host.
type Limiter struct {
	mu       sync.Mutex
	last     map[string]time.Time
	cooldown time.Duration
}

// New creates a Limiter with the given cooldown duration between
// allowed events for the same host.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		last:     make(map[string]time.Time),
		cooldown: cooldown,
	}
}

// Allow reports whether an event for the given host should be
// allowed through. If allowed, the timestamp is recorded.
func (l *Limiter) Allow(host string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if t, ok := l.last[host]; ok {
		if now.Sub(t) < l.cooldown {
			return false
		}
	}
	l.last[host] = now
	return true
}

// Reset clears the recorded timestamp for a host, allowing the
// next event to pass immediately.
func (l *Limiter) Reset(host string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, host)
}

// ResetAll clears all recorded timestamps.
func (l *Limiter) ResetAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}
