package ratelimit_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/ratelimit"
)

func TestAllow_FirstEventAlwaysPasses(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	if !l.Allow("host1") {
		t.Fatal("expected first event to be allowed")
	}
}

func TestAllow_SecondEventBlockedWithinCooldown(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("host1")
	if l.Allow("host1") {
		t.Fatal("expected second event to be blocked within cooldown")
	}
}

func TestAllow_DifferentHostsAreIndependent(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("host1")
	if !l.Allow("host2") {
		t.Fatal("expected different host to be allowed independently")
	}
}

func TestAllow_PassesAfterCooldown(t *testing.T) {
	l := ratelimit.New(10 * time.Millisecond)
	l.Allow("host1")
	time.Sleep(20 * time.Millisecond)
	if !l.Allow("host1") {
		t.Fatal("expected event to pass after cooldown expired")
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("host1")
	l.Reset("host1")
	if !l.Allow("host1") {
		t.Fatal("expected allow after reset")
	}
}

func TestResetAll_ClearsAllHosts(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("host1")
	l.Allow("host2")
	l.ResetAll()
	if !l.Allow("host1") || !l.Allow("host2") {
		t.Fatal("expected all hosts to be allowed after ResetAll")
	}
}
