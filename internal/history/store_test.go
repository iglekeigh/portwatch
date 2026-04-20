package history_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/history"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "history-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestFileStore_SaveAndLoad(t *testing.T) {
	st, err := history.NewFileStore(tempDir(t))
	if err != nil {
		t.Fatal(err)
	}
	entries := []history.Entry{{Host: "localhost", Ports: []int{80, 443}}}
	if err := st.Save("localhost", entries); err != nil {
		t.Fatal(err)
	}
	loaded, err := st.Load("localhost")
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 1 || len(loaded[0].Ports) != 2 {
		t.Errorf("unexpected loaded entries: %+v", loaded)
	}
}

func TestFileStore_MissingHost_ReturnsNil(t *testing.T) {
	st, _ := history.NewFileStore(tempDir(t))
	entries, err := st.Load("ghost")
	if err != nil {
		t.Fatal(err)
	}
	if entries != nil {
		t.Errorf("expected nil, got %v", entries)
	}
}

func TestFileStore_SafeFilename(t *testing.T) {
	dir := tempDir(t)
	st, _ := history.NewFileStore(dir)
	_ = st.Save("192.168.1.1:8080", []history.Entry{{Host: "192.168.1.1:8080", Ports: []int{8080}}})
	expected := filepath.Join(dir, "192.168.1.1_8080.json")
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", expected)
	}
}
