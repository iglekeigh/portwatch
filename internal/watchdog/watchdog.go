// Package watchdog provides a self-monitoring mechanism that alerts when
// the portwatch scanner itself stops producing results within an expected interval.
package watchdog

import (
	"context"
	"sync"
	"time"
)

// AlertFunc is called when a host scan has not produced results within the deadline.
type AlertFunc func(host string, silent time.Duration)

// Watchdog tracks the last scan time per host and fires an alert when
// no scan result has been received within the configured deadline.
type Watchdog struct {
	mu       sync.Mutex
	deadline time.Duration
	lastSeen map[string]time.Time
	alert    AlertFunc
}

// New creates a Watchdog with the given deadline and alert callback.
func New(deadline time.Duration, alert AlertFunc) *Watchdog {
	return &Watchdog{
		deadline: deadline,
		lastSeen: make(map[string]time.Time),
		alert:    alert,
	}
}

// Checkin records that a scan result was received for the given host.
func (w *Watchdog) Checkin(host string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastSeen[host] = time.Now()
}

// Remove stops tracking the given host.
func (w *Watchdog) Remove(host string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.lastSeen, host)
}

// Check evaluates all tracked hosts and fires the alert for any that
// have exceeded the deadline since their last check-in.
func (w *Watchdog) Check(now time.Time) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for host, last := range w.lastSeen {
		silent := now.Sub(last)
		if silent >= w.deadline {
			w.alert(host, silent)
		}
	}
}

// Run starts a background goroutine that calls Check on the given interval
// until the context is cancelled.
func (w *Watchdog) Run(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case t := <-ticker.C:
				w.Check(t)
			case <-ctx.Done():
				return
			}
		}
	}()
}
