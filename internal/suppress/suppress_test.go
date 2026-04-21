package suppress_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/suppress"
)

func TestShouldAlert_FirstCallAlwaysPasses(t *testing.T) {
	s := suppress.New(5 * time.Minute)
	if !s.ShouldAlert("host1", "abc123") {
		t.Fatal("expected first call to return true")
	}
}

func TestShouldAlert_SameHashWithinTTL_Suppressed(t *testing.T) {
	s := suppress.New(5 * time.Minute)
	s.ShouldAlert("host1", "abc123")
	if s.ShouldAlert("host1", "abc123") {
		t.Fatal("expected second call with same hash to be suppressed")
	}
}

func TestShouldAlert_DifferentHash_Passes(t *testing.T) {
	s := suppress.New(5 * time.Minute)
	s.ShouldAlert("host1", "abc123")
	if !s.ShouldAlert("host1", "def456") {
		t.Fatal("expected different hash to pass through")
	}
}

func TestShouldAlert_ExpiredTTL_Passes(t *testing.T) {
	s := suppress.New(10 * time.Millisecond)
	s.ShouldAlert("host1", "abc123")
	time.Sleep(20 * time.Millisecond)
	if !s.ShouldAlert("host1", "abc123") {
		t.Fatal("expected expired TTL to allow alert")
	}
}

func TestShouldAlert_DifferentHosts_Independent(t *testing.T) {
	s := suppress.New(5 * time.Minute)
	s.ShouldAlert("host1", "abc123")
	s.ShouldAlert("host2", "abc123")

	if s.ShouldAlert("host1", "abc123") {
		t.Fatal("host1 should be suppressed")
	}
	if s.ShouldAlert("host2", "abc123") {
		t.Fatal("host2 should be suppressed")
	}
	if !s.ShouldAlert("host3", "abc123") {
		t.Fatal("host3 is new and should pass")
	}
}

func TestReset_AllowsNextAlert(t *testing.T) {
	s := suppress.New(5 * time.Minute)
	s.ShouldAlert("host1", "abc123")
	s.Reset("host1")
	if !s.ShouldAlert("host1", "abc123") {
		t.Fatal("expected reset to allow next alert")
	}
}

func TestHosts_ReturnsTrackedHosts(t *testing.T) {
	s := suppress.New(5 * time.Minute)
	s.ShouldAlert("alpha", "h1")
	s.ShouldAlert("beta", "h2")

	hosts := s.Hosts()
	if len(hosts) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(hosts))
	}
}
