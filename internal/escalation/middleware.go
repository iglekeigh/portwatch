package escalation

import (
	"context"
	"fmt"

	"portwatch/internal/alert"
)

// Notifier is satisfied by any type that can send an alert event.
type Notifier interface {
	Notify(ctx context.Context, event alert.Event) error
}

// Middleware wraps a Notifier and annotates each event with the current
// escalation level for the host before forwarding it downstream.
// Events that have not yet reached the warning threshold are suppressed.
type Middleware struct {
	next    Notifier
	tracker *Tracker
}

// NewMiddleware returns a new escalation Middleware.
func NewMiddleware(next Notifier, tracker *Tracker) *Middleware {
	return &Middleware{next: next, tracker: tracker}
}

// Notify evaluates the escalation level for the event host. If the level is
// LevelNone the event is silently dropped. Otherwise the level is appended to
// the event summary and the event is forwarded to the next notifier.
func (m *Middleware) Notify(ctx context.Context, event alert.Event) error {
	if !event.Diff.HasChanges() {
		m.tracker.Resolve(event.Host)
		return nil
	}

	lvl := m.tracker.Evaluate(event.Host)
	if lvl == LevelNone {
		return nil
	}

	event.Summary = fmt.Sprintf("[%s] %s", lvl, event.Summary)
	return m.next.Notify(ctx, event)
}
