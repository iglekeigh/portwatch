package throttle

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// mockNotifier records calls and optionally returns an error.
type mockNotifier struct {
	called int
	err    error
}

func (m *mockNotifier) Notify(_ alert.Event) error {
	m.called++
	return m.err
}

func makeEvent(host string) alert.Event {
	return alert.Event{
		Host:       host,
		NewPorts:   []int{80},
		HasDiff:    true,
	}
}

func TestMiddleware_FirstNotificationPasses(t *testing.T) {
	inner := &mockNotifier{}
	mw := NewMiddleware(inner, New(100*time.Millisecond))

	if err := mw.Notify(makeEvent("host1")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inner.called != 1 {
		t.Fatalf("expected inner called once, got %d", inner.called)
	}
}

func TestMiddleware_SecondNotificationSuppressed(t *testing.T) {
	inner := &mockNotifier{}
	mw := NewMiddleware(inner, New(100*time.Millisecond))

	_ = mw.Notify(makeEvent("host1"))
	err := mw.Notify(makeEvent("host1"))

	if err == nil {
		t.Fatal("expected suppression error, got nil")
	}
	if inner.called != 1 {
		t.Fatalf("expected inner called once, got %d", inner.called)
	}
}

func TestMiddleware_DifferentHostsAreIndependent(t *testing.T) {
	inner := &mockNotifier{}
	mw := NewMiddleware(inner, New(100*time.Millisecond))

	_ = mw.Notify(makeEvent("host1"))
	if err := mw.Notify(makeEvent("host2")); err != nil {
		t.Fatalf("expected host2 to pass, got %v", err)
	}
	if inner.called != 2 {
		t.Fatalf("expected inner called twice, got %d", inner.called)
	}
}

func TestMiddleware_PassesAfterDelay(t *testing.T) {
	inner := &mockNotifier{}
	mw := NewMiddleware(inner, New(20*time.Millisecond))

	_ = mw.Notify(makeEvent("host1"))
	time.Sleep(30 * time.Millisecond)

	if err := mw.Notify(makeEvent("host1")); err != nil {
		t.Fatalf("expected notification after delay, got %v", err)
	}
	if inner.called != 2 {
		t.Fatalf("expected inner called twice, got %d", inner.called)
	}
}

func TestMiddleware_PropagatesInnerError(t *testing.T) {
	expected := errors.New("send failed")
	inner := &mockNotifier{err: expected}
	mw := NewMiddleware(inner, New(100*time.Millisecond))

	err := mw.Notify(makeEvent("host1"))
	if !errors.Is(err, expected) {
		t.Fatalf("expected inner error %v, got %v", expected, err)
	}
}
