// Package cooldown provides per-host scan cooldown enforcement,
// preventing scans from running more frequently than a configured interval.
package cooldown

import (
	"sync"
	"time"
)

// Cooldown tracks the last scan time per host and enforces a minimum
// interval between successive scans of the same host.
type Cooldown struct {
	mu       sync.Mutex
	last     map[string]time.Time
	interval time.Duration
}

// New creates a Cooldown with the given minimum interval between scans.
func New(interval time.Duration) *Cooldown {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &Cooldown{
		last:     make(map[string]time.Time),
		interval: interval,
	}
}

// Allow returns true if enough time has passed since the last scan of host.
// It records the current time as the new last-scan time when returning true.
func (c *Cooldown) Allow(host string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	if t, ok := c.last[host]; ok && now.Sub(t) < c.interval {
		return false
	}
	c.last[host] = now
	return true
}

// Reset clears the last-scan record for host, allowing an immediate scan.
func (c *Cooldown) Reset(host string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.last, host)
}

// Remaining returns how long until the host is eligible for another scan.
// Returns 0 if the host is already eligible.
func (c *Cooldown) Remaining(host string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	t, ok := c.last[host]
	if !ok {
		return 0
	}
	elapsed := time.Since(t)
	if elapsed >= c.interval {
		return 0
	}
	return c.interval - elapsed
}

// Hosts returns all hosts currently tracked by the cooldown.
func (c *Cooldown) Hosts() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	hosts := make([]string, 0, len(c.last))
	for h := range c.last {
		hosts = append(hosts, h)
	}
	return hosts
}
