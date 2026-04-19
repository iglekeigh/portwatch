package store

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func TestNew_CreatesEmptyStore(t *testing.T) {
	s, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("localhost")
	if ok {
		t.Fatal("expected no snapshot for unknown host")
	}
}

func TestSave_AndGet(t *testing.T) {
	path := tempPath(t)
	s, _ := New(path)

	snap := Snapshot{Host: "192.168.1.1", Ports: []int{22, 80, 443}}
	if err := s.Save(snap); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	got, ok := s.Get("192.168.1.1")
	if !ok {
		t.Fatal("expected snapshot to exist")
	}
	if len(got.Ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(got.Ports))
	}
}

func TestStore_PersistsAcrossReload(t *testing.T) {
	path := tempPath(t)

	s1, _ := New(path)
	_ = s1.Save(Snapshot{Host: "10.0.0.1", Ports: []int{8080}})

	s2, err := New(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	snap, ok := s2.Get("10.0.0.1")
	if !ok {
		t.Fatal("expected snapshot after reload")
	}
	if len(snap.Ports) != 1 || snap.Ports[0] != 8080 {
		t.Fatalf("unexpected ports: %v", snap.Ports)
	}
}

func TestStore_MissingFile_NoError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent", "state.json")
	_, err := New(path)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("unexpected error: %v", err)
	}
}
