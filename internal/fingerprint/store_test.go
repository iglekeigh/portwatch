package fingerprint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
)

func tempFingerprintDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "fingerprint-test-*")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestMemStore_SetAndGet(t *testing.T) {
	s := fingerprint.NewMemStore()
	r := fingerprint.Compute("host", []int{80, 443})

	if err := s.Set("host", r); err != nil {
		t.Fatalf("Set: %v", err)
	}

	got, ok := s.Get("host")
	if !ok {
		t.Fatal("expected Get to find result")
	}
	if got.Hash != r.Hash {
		t.Errorf("hash mismatch: %s != %s", got.Hash, r.Hash)
	}
}

func TestMemStore_MissingHost(t *testing.T) {
	s := fingerprint.NewMemStore()
	_, ok := s.Get("unknown")
	if ok {
		t.Error("expected Get to return false for unknown host")
	}
}

func TestFileStore_SetAndGet(t *testing.T) {
	dir := tempFingerprintDir(t)
	fs, err := fingerprint.NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	r := fingerprint.Compute("192.168.1.1", []int{22, 80})
	if err := fs.Set("192.168.1.1", r); err != nil {
		t.Fatalf("Set: %v", err)
	}

	got, ok := fs.Get("192.168.1.1")
	if !ok {
		t.Fatal("expected Get to return result")
	}
	if got.Hash != r.Hash {
		t.Errorf("hash mismatch: %s != %s", got.Hash, r.Hash)
	}
}

func TestFileStore_MissingHost_ReturnsFalse(t *testing.T) {
	dir := tempFingerprintDir(t)
	fs, _ := fingerprint.NewFileStore(dir)

	_, ok := fs.Get("ghost")
	if ok {
		t.Error("expected false for unknown host")
	}
}

func TestFileStore_SafeFilename(t *testing.T) {
	dir := tempFingerprintDir(t)
	fs, _ := fingerprint.NewFileStore(dir)

	r := fingerprint.Compute("host:8080", []int{8080})
	_ = fs.Set("host:8080", r)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 file, got %d", len(entries))
	}
	name := entries[0].Name()
	if filepath.Ext(name) != ".json" {
		t.Errorf("expected .json extension, got %s", name)
	}
}
