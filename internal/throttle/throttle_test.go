package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestAllow_FirstScanAlwaysPasses(t *testing.T) {
	th := throttle.New(5 * time.Second)
	ok, wait := th.Allow("localhost")
	if !ok {
		t.Fatalf("expected first scan to be allowed, got wait=%v", wait)
	}
	if wait != 0 {
		t.Errorf("expected zero wait, got %v", wait)
	}
}

func TestAllow_SecondScanBlockedWithinDelay(t *testing.T) {
	th := throttle.New(5 * time.Second)
	th.Allow("localhost")

	ok, wait := th.Allow("localhost")
	if ok {
		t.Fatal("expected second scan to be blocked")
	}
	if wait <= 0 || wait > 5*time.Second {
		t.Errorf("unexpected wait duration: %v", wait)
	}
}

func TestAllow_DifferentHostsAreIndependent(t *testing.T) {
	th := throttle.New(5 * time.Second)
	th.Allow("host-a")

	ok, _ := th.Allow("host-b")
	if !ok {
		t.Error("expected host-b to be allowed independently of host-a")
	}
}

func TestAllow_PassesAfterDelay(t *testing.T) {
	th := throttle.New(10 * time.Millisecond)
	th.Allow("localhost")

	time.Sleep(20 * time.Millisecond)

	ok, wait := th.Allow("localhost")
	if !ok {
		t.Errorf("expected scan to pass after delay, got wait=%v", wait)
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	th := throttle.New(1 * time.Hour)
	th.Allow("localhost")

	th.Reset("localhost")

	ok, _ := th.Allow("localhost")
	if !ok {
		t.Error("expected scan to be allowed after reset")
	}
}

func TestHosts_ReturnsTrackedHosts(t *testing.T) {
	th := throttle.New(5 * time.Second)
	th.Allow("alpha")
	th.Allow("beta")

	hosts := th.Hosts()
	if len(hosts) != 2 {
		t.Errorf("expected 2 hosts, got %d", len(hosts))
	}
}

func TestNew_DefaultDelay(t *testing.T) {
	// zero or negative delay should default to 30s
	th := throttle.New(0)
	th.Allow("localhost")
	ok, wait := th.Allow("localhost")
	if ok {
		t.Fatal("expected scan to be blocked with default delay")
	}
	if wait <= 0 {
		t.Errorf("expected positive wait with default delay, got %v", wait)
	}
}
