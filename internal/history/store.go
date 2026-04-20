package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Store defines persistence for scan history entries.
type Store interface {
	Load(host string) ([]Entry, error)
	Save(host string, entries []Entry) error
}

// FileStore persists history as JSON files under a directory.
type FileStore struct {
	Dir string
}

// NewFileStore creates a FileStore rooted at dir.
func NewFileStore(dir string) (*FileStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &FileStore{Dir: dir}, nil
}

func (f *FileStore) path(host string) string {
	safe := strings.NewReplacer(":", "_", "/", "_", "\\", "_").Replace(host)
	return filepath.Join(f.Dir, safe+".json")
}

// Load reads entries for host from disk.
func (f *FileStore) Load(host string) ([]Entry, error) {
	data, err := os.ReadFile(f.path(host))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// Save writes entries for host to disk.
func (f *FileStore) Save(host string, entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.path(host), data, 0o644)
}
