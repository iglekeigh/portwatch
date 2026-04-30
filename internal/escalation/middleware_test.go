package escalation

import (
	"context"
	"strings"
	"testing"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/scanner"
)

type captureNotifier struct {
	events []alert.Event
}

func (c *captureNotifier) Notify(_ context.Context, e alert.Event) error {
	c.events = append(c.events, e)
	return nil
}

func makeEscalationEvent(host string, newPorts []int) alert.Event {
	return alert.Event{
		Host:    host,
		Summary: "port change detected",
		Diff: scanner.Diff{
			New: newPorts,
		},
	}
}

func TestMiddleware_SuppressesBeforeWarning(t *testing.T) {
	now := time.Now()
	tr := New(DefaultConfig())
	tr.nowFn = func() time.Time { return now }

	cap := &captureNotifier{}
	mw := NewMiddleware(cap, tr)

	_ = mw.Notify(context.Background(), makeEscalationEvent("host1", []int{8080}))
	if len(cap.events) != 0 {
		t.Fatal("expected event suppressed before warning threshold")
	}
}

func TestMiddleware_ForwardsAtWarningLevel(t *testing.T) {
	now := time.Now()
	tr := New(DefaultConfig())
	tr.nowFn = func() time.Time { return now }

	// seed first-seen
	tr.Evaluate("host1")

	// advance past warning
	tr.nowFn = func() time.Time { return now.Add(6 * time.Minute) }

	cap := &captureNotifier{}
	mw := NewMiddleware(cap, tr)

	_ = mw.Notify(context.Background(), makeEscalationEvent("host1", []int{8080}))
	if len(cap.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(cap.events))
	}
	if !strings.Contains(cap.events[0].Summary, "warning") {
		t.Errorf("expected summary to contain 'warning', got %q", cap.events[0].Summary)
	}
}

func TestMiddleware_ResolvesOnNoDiff(t *testing.T) {
	now := time.Now()
	tr := New(DefaultConfig())
	tr.nowFn = func() time.Time { return now }
	tr.Evaluate("host1") // register

	cap := &captureNotifier{}
	mw := NewMiddleware(cap, tr)

	// send an event with no changes
	_ = mw.Notify(context.Background(), alert.Event{Host: "host1", Diff: scanner.Diff{}})

	if len(tr.Hosts()) != 0 {
		t.Fatal("expected host resolved after no-diff event")
	}
	if len(cap.events) != 0 {
		t.Fatal("expected no events forwarded for no-diff")
	}
}

func TestMiddleware_EmergencyLabelInSummary(t *testing.T) {
	now := time.Now()
	tr := New(DefaultConfig())
	tr.nowFn = func() time.Time { return now }
	tr.Evaluate("host1")

	tr.nowFn = func() time.Time { return now.Add(90 * time.Minute) }

	cap := &captureNotifier{}
	mw := NewMiddleware(cap, tr)

	_ = mw.Notify(context.Background(), makeEscalationEvent("host1", []int{443}))
	if len(cap.events) == 0 {
		t.Fatal("expected event forwarded")
	}
	if !strings.Contains(cap.events[0].Summary, "emergency") {
		t.Errorf("expected 'emergency' in summary, got %q", cap.events[0].Summary)
	}
}
