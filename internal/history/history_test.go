package history_test

import (
	"testing"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/scanner"
)

type memStore struct {
	data map[string][]history.Entry
}

func newMem() *memStore { return &memStore{data: map[string][]history.Entry{}} }

func (m *memStore) Load(host string) ([]history.Entry, error) {
	return m.data[host], nil
}
func (m *memStore) Save(host string, entries []history.Entry) error {
	m.data[host] = entries
	return nil
}

func TestRecord_AddsEntry(t *testing.T) {
	st := newMem()
	err := history.Record(st, "localhost", scanner.Result{Open: []int{80, 443}})
	if err != nil {
		t.Fatal(err)
	}
	entries, _ := st.Load("localhost")
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if len(entries[0].Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(entries[0].Ports))
	}
}

func TestLatest_ReturnsLastEntry(t *testing.T) {
	st := newMem()
	_ = history.Record(st, "host1", scanner.Result{Open: []int{22}})
	_ = history.Record(st, "host1", scanner.Result{Open: []int{22, 80}})
	e, err := history.Latest(st, "host1")
	if err != nil {
		t.Fatal(err)
	}
	if e == nil {
		t.Fatal("expected entry, got nil")
	}
	if len(e.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(e.Ports))
	}
}

func TestLatest_NoEntries(t *testing.T) {
	st := newMem()
	e, err := history.Latest(st, "unknown")
	if err != nil {
		t.Fatal(err)
	}
	if e != nil {
		t.Errorf("expected nil, got %+v", e)
	}
}

func TestRecord_CapsAtMaxEntries(t *testing.T) {
	st := newMem()
	for i := 0; i <= history.MaxEntries; i++ {
		_ = history.Record(st, "h", scanner.Result{Open: []int{i + 1}})
	}
	entries, _ := st.Load("h")
	if len(entries) != history.MaxEntries {
		t.Errorf("expected %d entries, got %d", history.MaxEntries, len(entries))
	}
}
