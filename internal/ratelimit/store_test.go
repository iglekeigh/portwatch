package ratelimit

import (
	"testing"
	"time"
)

func TestStore_SetAndGet(t *testing.T) {
	s := NewStore()
	e := &Entry{
		LastSeen:    time.Now(),
		CountHour:   3,
		WindowStart: time.Now(),
	}
	s.Set("host-a", e)

	got := s.Get("host-a")
	if got == nil {
		t.Fatal("expected entry, got nil")
	}
	if got.CountHour != 3 {
		t.Errorf("expected CountHour=3, got %d", got.CountHour)
	}
}

func TestStore_Get_Missing(t *testing.T) {
	s := NewStore()
	if got := s.Get("nonexistent"); got != nil {
		t.Errorf("expected nil for missing host, got %+v", got)
	}
}

func TestStore_Delete(t *testing.T) {
	s := NewStore()
	s.Set("host-b", &Entry{LastSeen: time.Now()})
	s.Delete("host-b")
	if got := s.Get("host-b"); got != nil {
		t.Error("expected nil after delete")
	}
}

func TestStore_Hosts(t *testing.T) {
	s := NewStore()
	s.Set("alpha", &Entry{LastSeen: time.Now()})
	s.Set("beta", &Entry{LastSeen: time.Now()})

	hosts := s.Hosts()
	if len(hosts) != 2 {
		t.Errorf("expected 2 hosts, got %d", len(hosts))
	}
}

func TestStore_Purge_RemovesOldEntries(t *testing.T) {
	s := NewStore()
	old := time.Now().Add(-2 * time.Hour)
	s.Set("stale", &Entry{LastSeen: old})
	s.Set("fresh", &Entry{LastSeen: time.Now()})

	removed := s.Purge(1 * time.Hour)
	if removed != 1 {
		t.Errorf("expected 1 removed, got %d", removed)
	}
	if s.Get("stale") != nil {
		t.Error("stale entry should have been purged")
	}
	if s.Get("fresh") == nil {
		t.Error("fresh entry should remain")
	}
}

func TestStore_Purge_NothingToRemove(t *testing.T) {
	s := NewStore()
	s.Set("recent", &Entry{LastSeen: time.Now()})

	removed := s.Purge(1 * time.Hour)
	if removed != 0 {
		t.Errorf("expected 0 removed, got %d", removed)
	}
}
