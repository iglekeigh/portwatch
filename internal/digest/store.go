package digest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Store persists per-host digests to a JSON file so that they survive
// process restarts.
type Store struct {
	mu   sync.RWMutex
	path string
	data map[string]string // host -> digest hex
}

// NewFileStore opens (or creates) a JSON digest store at the given path.
func NewFileStore(path string) (*Store, error) {
	s := &Store{path: path, data: make(map[string]string)}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

// Get returns the stored digest for host, or an empty string if none exists.
func (s *Store) Get(host string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[host]
}

// Set stores the digest for host and flushes to disk.
func (s *Store) Set(host, digest string) error {
	s.mu.Lock()
	s.data[host] = digest
	s.mu.Unlock()
	return s.save()
}

// safeFilename converts a host string into a filesystem-safe segment.
func safeFilename(host string) string {
	return strings.NewReplacer(":", "_", "/", "_", "\\", "_").Replace(host)
}

func (s *Store) load() error {
	b, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.data)
}

func (s *Store) save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}

// safeFilename is exported for tests via the unexported alias above.
var _ = safeFilename
