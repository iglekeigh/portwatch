package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/circuitbreaker"
)

func newFastBreaker(failures, successes int) *circuitbreaker.Breaker {
	return circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: failures,
		SuccessThreshold: successes,
		OpenDuration:     50 * time.Millisecond,
	})
}

func TestAllow_InitiallyPermits(t *testing.T) {
	b := newFastBreaker(3, 1)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	b := newFastBreaker(3, 1)
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if b.State() != circuitbreaker.StateOpen {
		t.Fatalf("expected StateOpen, got %v", b.State())
	}
	if err := b.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestOpenCircuit_TransitionsToHalfOpenAfterDuration(t *testing.T) {
	b := newFastBreaker(1, 1)
	b.RecordFailure()
	if b.State() != circuitbreaker.StateOpen {
		t.Fatal("expected open state")
	}
	time.Sleep(60 * time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil after open duration, got %v", err)
	}
	if b.State() != circuitbreaker.StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %v", b.State())
	}
}

func TestHalfOpen_SuccessCloses(t *testing.T) {
	b := newFastBreaker(1, 2)
	b.RecordFailure()
	time.Sleep(60 * time.Millisecond)
	_ = b.Allow() // transitions to half-open

	b.RecordSuccess()
	if b.State() != circuitbreaker.StateHalfOpen {
		t.Fatal("should still be half-open after one success")
	}
	b.RecordSuccess()
	if b.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected StateClosed after %d successes, got %v", 2, b.State())
	}
}

func TestHalfOpen_FailureReopens(t *testing.T) {
	b := newFastBreaker(1, 2)
	b.RecordFailure()
	time.Sleep(60 * time.Millisecond)
	_ = b.Allow()

	b.RecordFailure()
	if b.State() != circuitbreaker.StateOpen {
		t.Fatalf("expected StateOpen after failure in half-open, got %v", b.State())
	}
}

func TestDefaultConfig_Sensible(t *testing.T) {
	cfg := circuitbreaker.DefaultConfig()
	if cfg.FailureThreshold <= 0 {
		t.Error("FailureThreshold must be positive")
	}
	if cfg.SuccessThreshold <= 0 {
		t.Error("SuccessThreshold must be positive")
	}
	if cfg.OpenDuration <= 0 {
		t.Error("OpenDuration must be positive")
	}
}

func TestRecordSuccess_ResetFailCount(t *testing.T) {
	b := newFastBreaker(3, 1)
	b.RecordFailure()
	b.RecordFailure()
	b.RecordSuccess() // should reset consecutive failures
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != circuitbreaker.StateClosed {
		t.Fatal("circuit should remain closed; failure count was reset")
	}
}
