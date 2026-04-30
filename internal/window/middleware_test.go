package window

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

type captureNotifier struct {
	calls []alert.Event
}

func (c *captureNotifier) Notify(_ context.Context, e alert.Event) error {
	c.calls = append(c.calls, e)
	return nil
}

func makeWindowEvent(host string) alert.Event {
	return alert.Event{
		Host: host,
		Diff: scanner.Diff{
			New:    []int{80},
			Closed: []int{},
		},
	}
}

func TestMiddleware_SuppressesBelowThreshold(t *testing.T) {
	cap := &captureNotifier{}
	m := NewMiddleware(cap, 10*time.Second, 3)
	m.Notify(context.Background(), makeWindowEvent("host1"))
	m.Notify(context.Background(), makeWindowEvent("host1"))
	if len(cap.calls) != 0 {
		t.Fatalf("expected 0 calls before threshold, got %d", len(cap.calls))
	}
}

func TestMiddleware_ForwardsAtThreshold(t *testing.T) {
	cap := &captureNotifier{}
	m := NewMiddleware(cap, 10*time.Second, 3)
	for i := 0; i < 3; i++ {
		m.Notify(context.Background(), makeWindowEvent("host1"))
	}
	if len(cap.calls) != 1 {
		t.Fatalf("expected 1 call at threshold, got %d", len(cap.calls))
	}
}

func TestMiddleware_SkipsNoDiff(t *testing.T) {
	cap := &captureNotifier{}
	m := NewMiddleware(cap, 10*time.Second, 1)
	e := alert.Event{Host: "host1", Diff: scanner.Diff{}}
	m.Notify(context.Background(), e)
	if len(cap.calls) != 0 {
		t.Fatalf("expected 0 calls for empty diff")
	}
}

func TestMiddleware_PropagatesError(t *testing.T) {
	errNotifier := &errCapture{err: errors.New("boom")}
	m := NewMiddleware(errNotifier, 10*time.Second, 1)
	err := m.Notify(context.Background(), makeWindowEvent("host1"))
	if err == nil {
		t.Fatal("expected error to be propagated")
	}
}

type errCapture struct{ err error }

func (e *errCapture) Notify(_ context.Context, _ alert.Event) error { return e.err }
