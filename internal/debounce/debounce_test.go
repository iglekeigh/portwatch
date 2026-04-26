package debounce_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/debounce"
)

// recordingNotifier captures every event it receives.
type recordingNotifier struct {
	mu     sync.Mutex
	events []alert.Event
}

func (r *recordingNotifier) Notify(e alert.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, e)
	return nil
}

func (r *recordingNotifier) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.events)
}

func makeEvent(host string) alert.Event {
	return alert.Event{Host: host, OpenPorts: []int{80, 443}}
}

func TestDebounce_SingleEventForwarded(t *testing.T) {
	rec := &recordingNotifier{}
	d := debounce.New(50*time.Millisecond, rec)

	_ = d.Notify(makeEvent("host-a"))

	time.Sleep(120 * time.Millisecond)

	if rec.count() != 1 {
		t.Fatalf("expected 1 notification, got %d", rec.count())
	}
}

func TestDebounce_RapidEventsCoalesced(t *testing.T) {
	rec := &recordingNotifier{}
	d := debounce.New(80*time.Millisecond, rec)

	for i := 0; i < 5; i++ {
		_ = d.Notify(makeEvent("host-b"))
		time.Sleep(20 * time.Millisecond)
	}

	time.Sleep(160 * time.Millisecond)

	if rec.count() != 1 {
		t.Fatalf("rapid events should coalesce into 1 notification, got %d", rec.count())
	}
}

func TestDebounce_DifferentHostsAreIndependent(t *testing.T) {
	rec := &recordingNotifier{}
	d := debounce.New(50*time.Millisecond, rec)

	_ = d.Notify(makeEvent("host-c"))
	_ = d.Notify(makeEvent("host-d"))

	time.Sleep(120 * time.Millisecond)

	if rec.count() != 2 {
		t.Fatalf("expected 2 notifications for distinct hosts, got %d", rec.count())
	}
}

func TestDebounce_FlushCancelsPendingTimers(t *testing.T) {
	rec := &recordingNotifier{}
	d := debounce.New(200*time.Millisecond, rec)

	_ = d.Notify(makeEvent("host-e"))
	d.Flush()

	time.Sleep(300 * time.Millisecond)

	if rec.count() != 0 {
		t.Fatalf("flush should cancel pending timer, got %d notifications", rec.count())
	}
}
