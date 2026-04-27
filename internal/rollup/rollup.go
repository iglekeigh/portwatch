// Package rollup coalesces multiple scan results for the same host
// within a time window into a single aggregated diff event.
package rollup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Notifier is the downstream notification interface.
type Notifier interface {
	Notify(event alert.Event) error
}

// Window holds pending events for a host until the window expires.
type Window struct {
	mu       sync.Mutex
	window   time.Duration
	next     Notifier
	pending  map[string]*pendingEntry
	timers   map[string]*time.Timer
}

type pendingEntry struct {
	event alert.Event
}

// New creates a Window that batches events per host for the given duration
// before forwarding a merged event to next.
func New(window time.Duration, next Notifier) *Window {
	return &Window{
		window:  window,
		next:    next,
		pending: make(map[string]*pendingEntry),
		timers:  make(map[string]*time.Timer),
	}
}

// Notify buffers the event for the host. If an event for that host is already
// pending, the new and closed port lists are merged. The merged event is
// forwarded after the window expires.
func (w *Window) Notify(event alert.Event) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	host := event.Host
	if existing, ok := w.pending[host]; ok {
		existing.event = merge(existing.event, event)
		return nil
	}

	w.pending[host] = &pendingEntry{event: event}
	w.timers[host] = time.AfterFunc(w.window, func() {
		w.flush(host)
	})
	return nil
}

func (w *Window) flush(host string) {
	w.mu.Lock()
	entry, ok := w.pending[host]
	if !ok {
		w.mu.Unlock()
		return
	}
	delete(w.pending, host)
	delete(w.timers, host)
	w.mu.Unlock()

	_ = w.next.Notify(entry.event)
}

func merge(a, b alert.Event) alert.Event {
	seen := make(map[int]struct{})
	var newPorts []int
	for _, p := range append(a.Diff.New, b.Diff.New...) {
		if _, ok := seen[p]; !ok {
			seen[p] = struct{}{}
			newPorts = append(newPorts, p)
		}
	}
	seen = make(map[int]struct{})
	var closedPorts []int
	for _, p := range append(a.Diff.Closed, b.Diff.Closed...) {
		if _, ok := seen[p]; !ok {
			seen[p] = struct{}{}
			closedPorts = append(closedPorts, p)
		}
	}
	a.Diff.New = newPorts
	a.Diff.Closed = closedPorts
	return a
}
