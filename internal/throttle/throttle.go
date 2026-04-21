// Package throttle provides per-host scan throttling to prevent
// overwhelming targets with too-frequent port scans.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks the last scan time per host and enforces a minimum
// interval between successive scans of the same host.
type Throttle struct {
	mu       sync.Mutex
	lastScan map[string]time.Time
	minDelay time.Duration
}

// New creates a Throttle that enforces at least minDelay between scans
// of the same host. If minDelay is zero or negative, a default of 30s is used.
func New(minDelay time.Duration) *Throttle {
	if minDelay <= 0 {
		minDelay = 30 * time.Second
	}
	return &Throttle{
		lastScan: make(map[string]time.Time),
		minDelay: minDelay,
	}
}

// Allow reports whether a scan of host is permitted at now.
// If permitted, it records the scan time and returns true.
// If the host was scanned too recently, it returns false and the
// remaining wait duration.
func (t *Throttle) Allow(host string) (bool, time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if last, ok := t.lastScan[host]; ok {
		elapsed := now.Sub(last)
		if elapsed < t.minDelay {
			return false, t.minDelay - elapsed
		}
	}
	t.lastScan[host] = now
	return true, 0
}

// Reset clears the recorded scan time for host, allowing an immediate scan.
func (t *Throttle) Reset(host string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastScan, host)
}

// Hosts returns all hosts currently tracked by the throttle.
func (t *Throttle) Hosts() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	hosts := make([]string, 0, len(t.lastScan))
	for h := range t.lastScan {
		hosts = append(hosts, h)
	}
	return hosts
}
