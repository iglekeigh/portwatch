package rollup_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/rollup"
	"github.com/user/portwatch/internal/scanner"
)

func TestMiddleware_ForwardsAfterWindow(t *testing.T) {
	cap := &captureNotifier{}
	mw := rollup.NewMiddleware(50*time.Millisecond, cap)

	event := alert.Event{
		Host: "192.168.1.1",
		Diff: scanner.Diff{New: []int{22, 80}},
	}
	_ = mw.Notify(event)

	if len(cap.received()) != 0 {
		t.Fatal("expected no immediate notification before window expires")
	}

	time.Sleep(100 * time.Millisecond)

	events := cap.received()
	if len(events) != 1 {
		t.Fatalf("expected 1 event after window, got %d", len(events))
	}
	if events[0].Host != "192.168.1.1" {
		t.Errorf("unexpected host: %s", events[0].Host)
	}
}

func TestMiddleware_CoalescesWithinWindow(t *testing.T) {
	cap := &captureNotifier{}
	mw := rollup.NewMiddleware(80*time.Millisecond, cap)

	_ = mw.Notify(alert.Event{
		Host: "10.0.0.1",
		Diff: scanner.Diff{New: []int{80}},
	})
	time.Sleep(20 * time.Millisecond)
	_ = mw.Notify(alert.Event{
		Host: "10.0.0.1",
		Diff: scanner.Diff{New: []int{443}},
	})

	time.Sleep(150 * time.Millisecond)

	events := cap.received()
	if len(events) != 1 {
		t.Fatalf("expected 1 coalesced event, got %d", len(events))
	}
	if len(events[0].Diff.New) != 2 {
		t.Errorf("expected ports [80 443], got %v", events[0].Diff.New)
	}
}

func TestMiddleware_ImplementsNotifier(t *testing.T) {
	cap := &captureNotifier{}
	var _ interface {
		Notify(alert.Event) error
	} = rollup.NewMiddleware(time.Second, cap)
}
