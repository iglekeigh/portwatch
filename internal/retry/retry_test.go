package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/retry"
)

func fastCfg() retry.Config {
	return retry.Config{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Millisecond,
		Multiplier:   1.0,
	}
}

func TestDo_SucceedsFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastCfg(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastCfg(), func() error {
		calls++
		if calls < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAllAttempts(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastCfg(), func() error {
		calls++
		return errors.New("permanent failure")
	})
	if err == nil {
		t.Fatal("expected an error after exhausting retries")
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	calls := 0
	err := retry.Do(ctx, fastCfg(), func() error {
		calls++
		return nil
	})
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
	if calls != 0 {
		t.Fatalf("expected 0 calls with cancelled context, got %d", calls)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := retry.DefaultConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.InitialDelay != 500*time.Millisecond {
		t.Errorf("unexpected InitialDelay: %v", cfg.InitialDelay)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", cfg.Multiplier)
	}
}
