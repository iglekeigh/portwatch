package cooldown_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/cooldown"
)

func TestAllow_FirstScanAlwaysPasses(t *testing.T) {
	c := cooldown.New(5 * time.Second)
	if !c.Allow("host1") {
		t.Fatal("expected first scan to be allowed")
	}
}

func TestAllow_SecondScanBlockedWithinInterval(t *testing.T) {
	c := cooldown.New(5 * time.Second)
	c.Allow("host1")
	if c.Allow("host1") {
		t.Fatal("expected second scan within interval to be blocked")
	}
}

func TestAllow_DifferentHostsAreIndependent(t *testing.T) {
	c := cooldown.New(5 * time.Second)
	c.Allow("host1")
	if !c.Allow("host2") {
		t.Fatal("expected different host to be allowed independently")
	}
}

func TestAllow_PassesAfterInterval(t *testing.T) {
	c := cooldown.New(10 * time.Millisecond)
	c.Allow("host1")
	time.Sleep(20 * time.Millisecond)
	if !c.Allow("host1") {
		t.Fatal("expected scan to be allowed after interval elapsed")
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	c := cooldown.New(5 * time.Second)
	c.Allow("host1")
	c.Reset("host1")
	if !c.Allow("host1") {
		t.Fatal("expected scan to be allowed after reset")
	}
}

func TestRemaining_ReturnsZeroWhenEligible(t *testing.T) {
	c := cooldown.New(5 * time.Second)
	if r := c.Remaining("host1"); r != 0 {
		t.Fatalf("expected 0 remaining for unseen host, got %v", r)
	}
}

func TestRemaining_ReturnsPositiveWithinInterval(t *testing.T) {
	c := cooldown.New(5 * time.Second)
	c.Allow("host1")
	if r := c.Remaining("host1"); r <= 0 {
		t.Fatalf("expected positive remaining, got %v", r)
	}
}

func TestHosts_ReturnsTrackedHosts(t *testing.T) {
	c := cooldown.New(5 * time.Second)
	c.Allow("alpha")
	c.Allow("beta")
	hosts := c.Hosts()
	if len(hosts) != 2 {
		t.Fatalf("expected 2 tracked hosts, got %d", len(hosts))
	}
}

func TestNew_DefaultsNegativeInterval(t *testing.T) {
	c := cooldown.New(-1 * time.Second)
	c.Allow("host1")
	if r := c.Remaining("host1"); r <= 0 {
		t.Fatal("expected default interval to be applied for negative input")
	}
}
