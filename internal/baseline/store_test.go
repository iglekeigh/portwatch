package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/baseline"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "baseline-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestFileStore_SaveAndLoad(t *testing.T) {
	fs, err := baseline.NewFileStore(tempDir(t))
	if err != nil {
		t.Fatal(err)
	}
	snap := &baseline.Snapshot{
		Host:       "192.168.1.1",
		Ports:      []int{22, 80, 443},
		CapturedAt: time.Now().UTC(),
	}
	if err := fs.Save("192.168.1.1", snap); err != nil {
		t.Fatal(err)
	}
	loaded, err := fs.Load("192.168.1.1")
	if err != nil {
		t.Fatal(err)
	}
	if loaded == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if len(loaded.Ports) != 3 {
		t.Errorf("expected 3 ports, got %d", len(loaded.Ports))
	}
}

func TestFileStore_MissingHost_ReturnsNil(t *testing.T) {
	fs, _ := baseline.NewFileStore(tempDir(t))
	snap, err := fs.Load("nope")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap != nil {
		t.Errorf("expected nil, got %+v", snap)
	}
}

func TestFileStore_SafeFilename(t *testing.T) {
	dir := tempDir(t)
	fs, _ := baseline.NewFileStore(dir)
	_ = fs.Save("host:8080", &baseline.Snapshot{Host: "host:8080", Ports: []int{8080}})

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 file, got %d", len(entries))
	}
	name := entries[0].Name()
	if filepath.Ext(name) != ".json" {
		t.Errorf("expected .json extension, got %s", name)
	}
	for _, ch := range []string{":", "/"} {
		if contains(name, ch) {
			t.Errorf("filename contains unsafe char %q: %s", ch, name)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsRune(s, sub))
}

func containsRune(s, sub string) bool {
	for i := range s {
		if i+len(sub) <= len(s) && s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
