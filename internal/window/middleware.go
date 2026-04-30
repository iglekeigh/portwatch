package window

import (
	"context"
	"fmt"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Middleware wraps a Notifier and suppresses notifications until a minimum
// number of change events have been recorded within the rolling window.
type Middleware struct {
	next      alert.Notifier
	counter   *Counter
	threshold int
}

// NewMiddleware creates a Middleware that only forwards events to next when
// the host has accumulated at least threshold events within window.
func NewMiddleware(next alert.Notifier, window time.Duration, threshold int) *Middleware {
	return &Middleware{
		next:      next,
		counter:   New(window),
		threshold: threshold,
	}
}

// Notify records the event and forwards it only when the threshold is met.
func (m *Middleware) Notify(ctx context.Context, event alert.Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}
	m.counter.Record(event.Host)
	count := m.counter.Count(event.Host)
	if count < m.threshold {
		return nil
	}
	if err := m.next.Notify(ctx, event); err != nil {
		return fmt.Errorf("window middleware: %w", err)
	}
	return nil
}
