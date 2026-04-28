package watchdog_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

func TestCheckin_PreventsAlert(t *testing.T) {
	var fired bool
	w := watchdog.New(5*time.Second, func(_ string, _ time.Duration) {
		fired = true
	})
	w.Checkin("host-a")
	w.Check(time.Now().Add(2 * time.Second))
	if fired {
		t.Fatal("expected no alert within deadline")
	}
}

func TestCheck_FiresAfterDeadline(t *testing.T) {
	var mu sync.Mutex
	var alerts []string
	w := watchdog.New(5*time.Second, func(host string, _ time.Duration) {
		mu.Lock()
		alerts = append(alerts, host)
		mu.Unlock()
	})
	w.Checkin("host-b")
	w.Check(time.Now().Add(10 * time.Second))
	mu.Lock()
	defer mu.Unlock()
	if len(alerts) != 1 || alerts[0] != "host-b" {
		t.Fatalf("expected alert for host-b, got %v", alerts)
	}
}

func TestRemove_StopsTracking(t *testing.T) {
	var fired bool
	w := watchdog.New(1*time.Second, func(_ string, _ time.Duration) {
		fired = true
	})
	w.Checkin("host-c")
	w.Remove("host-c")
	w.Check(time.Now().Add(10 * time.Second))
	if fired {
		t.Fatal("expected no alert after remove")
	}
}

func TestCheck_SilentDurationReported(t *testing.T) {
	var reported time.Duration
	w := watchdog.New(1*time.Second, func(_ string, d time.Duration) {
		reported = d
	})
	w.Checkin("host-d")
	advance := 7 * time.Second
	w.Check(time.Now().Add(advance))
	if reported < advance {
		t.Fatalf("expected silent duration >= %v, got %v", advance, reported)
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	w := watchdog.New(100*time.Millisecond, func(_ string, _ time.Duration) {})
	ctx, cancel := context.WithCancel(context.Background())
	w.Run(ctx, 20*time.Millisecond)
	cancel()
	// Allow goroutine to exit cleanly; no panic expected.
	time.Sleep(50 * time.Millisecond)
}
