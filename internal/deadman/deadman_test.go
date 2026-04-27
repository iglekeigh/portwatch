package deadman_test

import (
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/deadman"
)

func TestCheckin_PreventsAlert(t *testing.T) {
	var fired []string
	w := deadman.New(50*time.Millisecond, func(host string, _ time.Duration) {
		fired = append(fired, host)
	})
	w.Checkin("host-a")
	w.Check()
	if len(fired) != 0 {
		t.Fatalf("expected no alerts, got %v", fired)
	}
}

func TestCheck_FiresAfterTimeout(t *testing.T) {
	var mu sync.Mutex
	var fired []string
	w := deadman.New(10*time.Millisecond, func(host string, _ time.Duration) {
		mu.Lock()
		fired = append(fired, host)
		mu.Unlock()
	})
	w.Checkin("host-b")
	time.Sleep(30 * time.Millisecond)
	w.Check()
	mu.Lock()
	defer mu.Unlock()
	if len(fired) != 1 || fired[0] != "host-b" {
		t.Fatalf("expected alert for host-b, got %v", fired)
	}
}

func TestRemove_StopsTracking(t *testing.T) {
	var fired []string
	w := deadman.New(10*time.Millisecond, func(host string, _ time.Duration) {
		fired = append(fired, host)
	})
	w.Checkin("host-c")
	w.Remove("host-c")
	time.Sleep(20 * time.Millisecond)
	w.Check()
	if len(fired) != 0 {
		t.Fatalf("expected no alerts after remove, got %v", fired)
	}
}

func TestHosts_ReturnTracked(t *testing.T) {
	w := deadman.New(time.Minute, func(_ string, _ time.Duration) {})
	w.Checkin("alpha")
	w.Checkin("beta")
	hosts := w.Hosts()
	sort.Strings(hosts)
	if len(hosts) != 2 || hosts[0] != "alpha" || hosts[1] != "beta" {
		t.Fatalf("unexpected hosts: %v", hosts)
	}
}

func TestCheck_SilentDurationReported(t *testing.T) {
	var reported time.Duration
	w := deadman.New(10*time.Millisecond, func(_ string, d time.Duration) {
		reported = d
	})
	w.Checkin("host-d")
	time.Sleep(30 * time.Millisecond)
	w.Check()
	if reported < 10*time.Millisecond {
		t.Fatalf("expected silent duration >= 10ms, got %v", reported)
	}
}
