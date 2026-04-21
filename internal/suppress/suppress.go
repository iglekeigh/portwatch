// Package suppress provides a mechanism to suppress duplicate alerts
// for a host when no changes have occurred between scan cycles.
package suppress

import (
	"sync"
	"time"
)

// Entry holds suppression state for a single host.
type Entry struct {
	LastAlerted time.Time
	LastHash    string
}

// Suppressor tracks per-host alert state and suppresses redundant notifications.
type Suppressor struct {
	mu      sync.Mutex
	entries map[string]*Entry
	ttl     time.Duration
}

// New creates a Suppressor with the given TTL. Hosts that have not been
// alerted within the TTL window are treated as unsuppressed.
func New(ttl time.Duration) *Suppressor {
	return &Suppressor{
		entries: make(map[string]*Entry),
		ttl:     ttl,
	}
}

// ShouldAlert returns true when the host should trigger a notification.
// It suppresses alerts when the port hash is identical to the last alerted
// state and the TTL has not yet expired.
func (s *Suppressor) ShouldAlert(host, portHash string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.entries[host]
	if !ok {
		s.entries[host] = &Entry{LastAlerted: time.Now(), LastHash: portHash}
		return true
	}

	expired := time.Since(e.LastAlerted) >= s.ttl
	changed := e.LastHash != portHash

	if changed || expired {
		e.LastAlerted = time.Now()
		e.LastHash = portHash
		return true
	}

	return false
}

// Reset clears suppression state for a host, forcing the next call to
// ShouldAlert to return true regardless of the hash.
func (s *Suppressor) Reset(host string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, host)
}

// Hosts returns all hosts currently tracked by the suppressor.
func (s *Suppressor) Hosts() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, 0, len(s.entries))
	for h := range s.entries {
		out = append(out, h)
	}
	return out
}
