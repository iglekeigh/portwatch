// Package window provides a time-windowed event counter for tracking
// scan result frequencies over a rolling duration.
package window

import (
	"sync"
	"time"
)

// Counter tracks how many events occurred for a host within a sliding window.
type Counter struct {
	mu      sync.Mutex
	window  time.Duration
	entries map[string][]time.Time
	now     func() time.Time
}

// New creates a Counter with the given rolling window duration.
func New(window time.Duration) *Counter {
	return &Counter{
		window:  window,
		entries: make(map[string][]time.Time),
		now:     time.Now,
	}
}

// Record adds an event for the given host at the current time.
func (c *Counter) Record(host string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[host] = c.prune(host)
	c.entries[host] = append(c.entries[host], c.now())
}

// Count returns the number of events recorded for a host within the window.
func (c *Counter) Count(host string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[host] = c.prune(host)
	return len(c.entries[host])
}

// Reset clears all recorded events for a host.
func (c *Counter) Reset(host string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, host)
}

// Hosts returns all hosts currently tracked.
func (c *Counter) Hosts() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]string, 0, len(c.entries))
	for h := range c.entries {
		out = append(out, h)
	}
	return out
}

// prune removes timestamps older than the window. Must be called with lock held.
func (c *Counter) prune(host string) []time.Time {
	cutoff := c.now().Add(-c.window)
	old := c.entries[host]
	var fresh []time.Time
	for _, t := range old {
		if t.After(cutoff) {
			fresh = append(fresh, t)
		}
	}
	return fresh
}
