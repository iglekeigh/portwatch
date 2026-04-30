package backoff_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/backoff"
)

func TestDefaultConfig_ReturnsValidConfig(t *testing.T) {
	cfg := backoff.DefaultConfig()
	if cfg.InitialInterval <= 0 {
		t.Fatal("expected positive InitialInterval")
	}
	if cfg.MaxInterval < cfg.InitialInterval {
		t.Fatal("MaxInterval should be >= InitialInterval")
	}
	if cfg.Multiplier <= 1.0 {
		t.Fatal("Multiplier should be > 1")
	}
}

func TestBackoff_ZeroAttempt_ReturnsNearInitial(t *testing.T) {
	cfg := backoff.Config{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
		JitterFraction:  0,
	}
	d := cfg.Backoff(0)
	if d != 100*time.Millisecond {
		t.Fatalf("expected 100ms, got %v", d)
	}
}

func TestBackoff_GrowsExponentially(t *testing.T) {
	cfg := backoff.Config{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
		JitterFraction:  0,
	}
	d0 := cfg.Backoff(0)
	d1 := cfg.Backoff(1)
	d2 := cfg.Backoff(2)
	if d1 != 2*d0 {
		t.Fatalf("expected d1=2*d0, got d0=%v d1=%v", d0, d1)
	}
	if d2 != 2*d1 {
		t.Fatalf("expected d2=2*d1, got d1=%v d2=%v", d1, d2)
	}
}

func TestBackoff_CapsAtMaxInterval(t *testing.T) {
	cfg := backoff.Config{
		InitialInterval: 1 * time.Second,
		MaxInterval:     4 * time.Second,
		Multiplier:      2.0,
		JitterFraction:  0,
	}
	for i := 5; i < 20; i++ {
		d := cfg.Backoff(i)
		if d > cfg.MaxInterval {
			t.Fatalf("attempt %d: %v exceeds MaxInterval %v", i, d, cfg.MaxInterval)
		}
	}
}

func TestBackoff_NegativeAttempt_TreatedAsZero(t *testing.T) {
	cfg := backoff.Config{
		InitialInterval: 200 * time.Millisecond,
		MaxInterval:     5 * time.Second,
		Multiplier:      2.0,
		JitterFraction:  0,
	}
	if cfg.Backoff(-3) != cfg.Backoff(0) {
		t.Fatal("negative attempt should behave like attempt 0")
	}
}

func TestSeries_LengthMatchesN(t *testing.T) {
	cfg := backoff.DefaultConfig()
	s := cfg.Series(5)
	if len(s) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(s))
	}
}

func TestBackoff_JitterAddsVariance(t *testing.T) {
	cfg := backoff.Config{
		InitialInterval: 1 * time.Second,
		MaxInterval:     10 * time.Second,
		Multiplier:      1.0,
		JitterFraction:  0.5,
	}
	seen := map[time.Duration]bool{}
	for i := 0; i < 20; i++ {
		seen[cfg.Backoff(0)] = true
	}
	if len(seen) < 2 {
		t.Fatal("expected jitter to produce varied durations")
	}
}
