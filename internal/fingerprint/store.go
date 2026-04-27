package fingerprint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Store persists fingerprint results keyed by host.
type Store interface {
	Get(host string) (Result, bool)
	Set(host string, r Result) error
}

// MemStore is an in-memory Store used for testing.
type MemStore struct {
	mu   sync.RWMutex
	data map[string]Result
}

// NewMemStore returns an empty in-memory store.
func NewMemStore() *MemStore {
	return &MemStore{data: make(map[string]Result)}
}

func (m *MemStore) Get(host string) (Result, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.data[host]
	return r, ok
}

func (m *MemStore) Set(host string, r Result) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[host] = r
	return nil
}

// FileStore persists fingerprints as JSON files under a directory.
type FileStore struct {
	dir string
}

// NewFileStore returns a FileStore rooted at dir, creating it if needed.
func NewFileStore(dir string) (*FileStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("fingerprint: mkdir %s: %w", dir, err)
	}
	return &FileStore{dir: dir}, nil
}

func safeFilename(host string) string {
	return strings.NewReplacer(":", "_", "/", "_", "\\", "_").Replace(host) + ".json"
}

func (f *FileStore) path(host string) string {
	return filepath.Join(f.dir, safeFilename(host))
}

func (f *FileStore) Get(host string) (Result, bool) {
	data, err := os.ReadFile(f.path(host))
	if err != nil {
		return Result{}, false
	}
	var r Result
	if err := json.Unmarshal(data, &r); err != nil {
		return Result{}, false
	}
	return r, true
}

func (f *FileStore) Set(host string, r Result) error {
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("fingerprint: marshal: %w", err)
	}
	return os.WriteFile(f.path(host), data, 0o644)
}
