package rollup_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/rollup"
	"github.com/user/portwatch/internal/scanner"
)

type captureNotifier struct {
	mu     sync.Mutex
	events []alert.Event
}

func (c *captureNotifier) Notify(e alert.Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = append(c.events, e)
	return nil
}

func (c *captureNotifier) received() []alert.Event {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]alert.Event(nil), c.events...)
}

func makeEvent(host string, newPorts, closed []int) alert.Event {
	return alert.Event{
		Host: host,
		Diff: scanner.Diff{
			New:    newPorts,
			Closed: closed,
		},
	}
}

func TestRollup_SingleEventForwarded(t *testing.T) {
	cap := &captureNotifier{}
	w := rollup.New(50*time.Millisecond, cap)

	_ = w.Notify(makeEvent("host1", []int{80}, nil))
	time.Sleep(100 * time.Millisecond)

	events := cap.received()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if len(events[0].Diff.New) != 1 || events[0].Diff.New[0] != 80 {
		t.Errorf("unexpected new ports: %v", events[0].Diff.New)
	}
}

func TestRollup_MultipleEventsCoalesced(t *testing.T) {
	cap := &captureNotifier{}
	w := rollup.New(80*time.Millisecond, cap)

	_ = w.Notify(makeEvent("host1", []int{80}, nil))
	_ = w.Notify(makeEvent("host1", []int{443}, []int{22}))
	time.Sleep(150 * time.Millisecond)

	events := cap.received()
	if len(events) != 1 {
		t.Fatalf("expected 1 coalesced event, got %d", len(events))
	}
	if len(events[0].Diff.New) != 2 {
		t.Errorf("expected 2 new ports, got %v", events[0].Diff.New)
	}
	if len(events[0].Diff.Closed) != 1 || events[0].Diff.Closed[0] != 22 {
		t.Errorf("unexpected closed ports: %v", events[0].Diff.Closed)
	}
}

func TestRollup_DifferentHostsAreIndependent(t *testing.T) {
	cap := &captureNotifier{}
	w := rollup.New(50*time.Millisecond, cap)

	_ = w.Notify(makeEvent("host1", []int{80}, nil))
	_ = w.Notify(makeEvent("host2", []int{443}, nil))
	time.Sleep(120 * time.Millisecond)

	events := cap.received()
	if len(events) != 2 {
		t.Fatalf("expected 2 events (one per host), got %d", len(events))
	}
}

func TestRollup_DeduplicatesPorts(t *testing.T) {
	cap := &captureNotifier{}
	w := rollup.New(60*time.Millisecond, cap)

	_ = w.Notify(makeEvent("host1", []int{80, 443}, nil))
	_ = w.Notify(makeEvent("host1", []int{80, 8080}, nil))
	time.Sleep(130 * time.Millisecond)

	events := cap.received()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if len(events[0].Diff.New) != 3 {
		t.Errorf("expected 3 deduplicated new ports, got %v", events[0].Diff.New)
	}
}
