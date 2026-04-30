package jitter_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/jitter"
)

func TestNew_ZeroMax_NoDelay(t *testing.T) {
	j := jitter.New(0)
	start := time.Now()
	if err := j.Wait(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 10*time.Millisecond {
		t.Errorf("expected no delay, got %v", elapsed)
	}
}

func TestSample_WithinBounds(t *testing.T) {
	max := 100 * time.Millisecond
	j := jitter.New(max)
	for i := 0; i < 50; i++ {
		d := j.Sample()
		if d < 0 || d >= max {
			t.Errorf("sample %v out of [0, %v)", d, max)
		}
	}
}

func TestSample_ZeroMax_AlwaysZero(t *testing.T) {
	j := jitter.New(0)
	for i := 0; i < 10; i++ {
		if d := j.Sample(); d != 0 {
			t.Errorf("expected 0, got %v", d)
		}
	}
}

func TestWait_CancelledContext_ReturnsError(t *testing.T) {
	j := jitter.New(10 * time.Second) // large delay to ensure cancel fires first
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	start := time.Now()
	err := j.Wait(ctx)
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Errorf("Wait did not return promptly on cancel: %v", elapsed)
	}
}

func TestWait_SmallMax_CompletesQuickly(t *testing.T) {
	j := jitter.New(20 * time.Millisecond)
	start := time.Now()
	if err := j.Wait(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Errorf("Wait took too long: %v", elapsed)
	}
}
