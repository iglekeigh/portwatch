// Package deadman provides a dead-man's switch that fires an alert
// when a host has not been scanned successfully within a configured interval.
package deadman

import (
	"sync"
	"time"
)

// Watcher tracks the last successful scan time per host and fires
// a callback when a host goes silent beyond the configured timeout.
type Watcher struct {
	mu      sync.Mutex
	last    map[string]time.Time
	timeout time.Duration
	onDead  func(host string, silent time.Duration)
}

// New creates a Watcher with the given silence timeout and dead callback.
func New(timeout time.Duration, onDead func(host string, silent time.Duration)) *Watcher {
	return &Watcher{
		last:    make(map[string]time.Time),
		timeout: timeout,
		onDead:  onDead,
	}
}

// Checkin records a successful scan for the given host.
func (w *Watcher) Checkin(host string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.last[host] = time.Now()
}

// Check evaluates all tracked hosts and fires onDead for any that
// have exceeded the silence timeout. It is safe to call concurrently.
func (w *Watcher) Check() {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now()
	for host, t := range w.last {
		silent := now.Sub(t)
		if silent > w.timeout {
			w.onDead(host, silent)
		}
	}
}

// Remove stops tracking the given host.
func (w *Watcher) Remove(host string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.last, host)
}

// Hosts returns all currently tracked host names.
func (w *Watcher) Hosts() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]string, 0, len(w.last))
	for h := range w.last {
		out = append(out, h)
	}
	return out
}
