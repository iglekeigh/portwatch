package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Snapshot holds the last known open ports for a host.
type Snapshot struct {
	Host  string `json:"host"`
	Ports []int  `json:"ports"`
}

// Store persists port snapshots to disk.
type Store struct {
	mu   sync.RWMutex
	path string
	data map[string]Snapshot
}

// New opens (or creates) a JSON store at the given file path.
func New(path string) (*Store, error) {
	s := &Store{
		path: path,
		data: make(map[string]Snapshot),
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// Get returns the stored snapshot for a host, and whether it existed.
func (s *Store) Get(host string) (Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap, ok := s.data[host]
	return snap, ok
}

// Save stores a snapshot for a host and flushes to disk.
func (s *Store) Save(snap Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[snap.Host] = snap
	return s.flush()
}

// Delete removes the snapshot for a host and flushes to disk.
// It returns false (without error) if the host was not present.
func (s *Store) Delete(host string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[host]; !ok {
		return false, nil
	}
	delete(s.data, host)
	return true, s.flush()
}

// Hosts returns a slice of all host names currently in the store.
func (s *Store) Hosts() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	hosts := make([]string, 0, len(s.data))
	for h := range s.data {
		hosts = append(hosts, h)
	}
	return hosts
}

func (s *Store) load() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.data)
}

func (s *Store) flush() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s.data)
}
