package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Store is the persistence interface for baseline snapshots.
type Store interface {
	Save(host string, snap *Snapshot) error
	Load(host string) (*Snapshot, error)
}

// FileStore persists baseline snapshots as JSON files in a directory.
type FileStore struct {
	dir string
}

// NewFileStore returns a FileStore rooted at dir, creating it if needed.
func NewFileStore(dir string) (*FileStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("baseline: mkdir %s: %w", dir, err)
	}
	return &FileStore{dir: dir}, nil
}

// Save writes the snapshot to disk as <dir>/<safe-host>.json.
func (fs *FileStore) Save(host string, snap *Snapshot) error {
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	return os.WriteFile(fs.path(host), data, 0o644)
}

// Load reads the snapshot for host from disk.
// Returns nil, nil if the file does not exist.
func (fs *FileStore) Load(host string) (*Snapshot, error) {
	data, err := os.ReadFile(fs.path(host))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("baseline: read: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &snap, nil
}

func (fs *FileStore) path(host string) string {
	safe := strings.NewReplacer(":", "_", "/", "_", "\\", "_", ".", "_").Replace(host)
	return filepath.Join(fs.dir, safe+".json")
}
