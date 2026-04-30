package window

import (
	"testing"
	"time"
)

func TestRecord_IncreasesCount(t *testing.T) {
	c := New(10 * time.Second)
	c.Record("host1")
	c.Record("host1")
	if got := c.Count("host1"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestCount_ZeroForUnknownHost(t *testing.T) {
	c := New(10 * time.Second)
	if got := c.Count("unknown"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestCount_PrunesExpiredEntries(t *testing.T) {
	now := time.Now()
	c := New(5 * time.Second)
	// inject a fake clock
	c.now = func() time.Time { return now }
	c.Record("host1")
	c.Record("host1")
	// advance past window
	c.now = func() time.Time { return now.Add(10 * time.Second) }
	if got := c.Count("host1"); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestReset_ClearsHost(t *testing.T) {
	c := New(10 * time.Second)
	c.Record("host1")
	c.Reset("host1")
	if got := c.Count("host1"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestHosts_ReturnsTrackedHosts(t *testing.T) {
	c := New(10 * time.Second)
	c.Record("alpha")
	c.Record("beta")
	hosts := c.Hosts()
	if len(hosts) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(hosts))
	}
}

func TestDifferentHosts_CountedIndependently(t *testing.T) {
	c := New(10 * time.Second)
	c.Record("a")
	c.Record("a")
	c.Record("b")
	if c.Count("a") != 2 {
		t.Fatalf("expected 2 for a")
	}
	if c.Count("b") != 1 {
		t.Fatalf("expected 1 for b")
	}
}
