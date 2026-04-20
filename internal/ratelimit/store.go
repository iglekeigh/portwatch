package ratelimit

import (
	"sync"
	"time"
)

// Entry holds rate limit state for a single host.
type Entry struct {
	LastSeen  time.Time
	CountHour int
	WindowStart time.Time
}

// Store is a thread-safe in-memory store for rate limit entries.
type Store struct {
	mu      sync.Mutex
	entries map[string]*Entry
}

// NewStore creates a new in-memory rate limit store.
func NewStore() *Store {
	return &Store{
		entries: make(map[string]*Entry),
	}
}

// Get returns the entry for a host, or nil if not found.
func (s *Store) Get(host string) *Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.entries[host]
}

// Set stores or updates the entry for a host.
func (s *Store) Set(host string, e *Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[host] = e
}

// Delete removes the entry for a host.
func (s *Store) Delete(host string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, host)
}

// Hosts returns a snapshot of all tracked host keys.
func (s *Store) Hosts() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.entries))
	for k := range s.entries {
		keys = append(keys, k)
	}
	return keys
}

// Purge removes entries whose last-seen time is older than the given duration.
func (s *Store) Purge(olderThan time.Duration) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	cutoff := time.Now().Add(-olderThan)
	removed := 0
	for host, e := range s.entries {
		if e.LastSeen.Before(cutoff) {
			delete(s.entries, host)
			removed++
		}
	}
	return removed
}
