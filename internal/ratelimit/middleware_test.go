package ratelimit

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// stubNotifier records every event it receives.
type stubNotifier struct {
	events []alert.Event
	err    error
}

func (s *stubNotifier) Notify(event alert.Event) error {
	s.events = append(s.events, event)
	return s.err
}

func makeEvent(host string) alert.Event {
	return alert.Event{Host: host}
}

func TestMiddleware_FirstNotificationPasses(t *testing.T) {
	stub := &stubNotifier{}
	mw := NewMiddleware(stub, 5*time.Minute)

	if err := mw.Notify(makeEvent("host-a")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(stub.events) != 1 {
		t.Fatalf("expected 1 event forwarded, got %d", len(stub.events))
	}
}

func TestMiddleware_SecondNotificationSuppressed(t *testing.T) {
	stub := &stubNotifier{}
	mw := NewMiddleware(stub, 5*time.Minute)

	_ = mw.Notify(makeEvent("host-a"))
	err := mw.Notify(makeEvent("host-a"))

	if err == nil {
		t.Fatal("expected suppression error, got nil")
	}
	if len(stub.events) != 1 {
		t.Fatalf("expected 1 forwarded event, got %d", len(stub.events))
	}
}

func TestMiddleware_DifferentHostsAreIndependent(t *testing.T) {
	stub := &stubNotifier{}
	mw := NewMiddleware(stub, 5*time.Minute)

	_ = mw.Notify(makeEvent("host-a"))
	if err := mw.Notify(makeEvent("host-b")); err != nil {
		t.Fatalf("expected host-b to pass, got %v", err)
	}
	if len(stub.events) != 2 {
		t.Fatalf("expected 2 forwarded events, got %d", len(stub.events))
	}
}

func TestMiddleware_ResetAllowsImmediateRetry(t *testing.T) {
	stub := &stubNotifier{}
	mw := NewMiddleware(stub, 5*time.Minute)

	_ = mw.Notify(makeEvent("host-a"))
	mw.Reset("host-a")

	if err := mw.Notify(makeEvent("host-a")); err != nil {
		t.Fatalf("expected notification after reset to pass, got %v", err)
	}
	if len(stub.events) != 2 {
		t.Fatalf("expected 2 forwarded events after reset, got %d", len(stub.events))
	}
}

func TestMiddleware_InnerErrorPropagated(t *testing.T) {
	want := errors.New("smtp failure")
	stub := &stubNotifier{err: want}
	mw := NewMiddleware(stub, 5*time.Minute)

	if err := mw.Notify(makeEvent("host-a")); !errors.Is(err, want) {
		t.Fatalf("expected inner error %v, got %v", want, err)
	}
}
