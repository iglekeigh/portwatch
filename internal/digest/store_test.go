package digest

import (
	"os"
	"path/filepath"
	"testing"
)

func tempDigestPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "digests.json")
}

func TestFileStore_SetAndGet(t *testing.T) {
	p := tempDigestPath(t)
	s, err := NewFileStore(p)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	if err := s.Set("localhost", "abc123"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if got := s.Get("localhost"); got != "abc123" {
		t.Errorf("Get = %q, want %q", got, "abc123")
	}
}

func TestFileStore_MissingHost_ReturnsEmpty(t *testing.T) {
	s, _ := NewFileStore(tempDigestPath(t))
	if got := s.Get("unknown"); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestFileStore_PersistsAcrossReload(t *testing.T) {
	p := tempDigestPath(t)
	s1, _ := NewFileStore(p)
	_ = s1.Set("10.0.0.1", "deadbeef")

	s2, err := NewFileStore(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got := s2.Get("10.0.0.1"); got != "deadbeef" {
		t.Errorf("after reload Get = %q, want %q", got, "deadbeef")
	}
}

func TestFileStore_MissingFile_NoError(t *testing.T) {
	p := filepath.Join(t.TempDir(), "nonexistent", "digests.json")
	// parent dir does not exist — NewFileStore should still succeed (load skips)
	_, err := NewFileStore(p)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
}

func TestFileStore_WritesJSON(t *testing.T) {
	p := tempDigestPath(t)
	s, _ := NewFileStore(p)
	_ = s.Set("host", "ff00")

	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty JSON file")
	}
}
