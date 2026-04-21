package suppress

import (
	"context"

	"github.com/user/portwatch/internal/alert"
)

// Notifier is the interface for sending alert notifications.
type Notifier interface {
	Notify(ctx context.Context, event alert.Event) error
}

// Middleware wraps a Notifier and suppresses duplicate alerts
// within a configurable TTL window using a Suppressor.
type Middleware struct {
	next       Notifier
	suppressor *Suppressor
}

// NewMiddleware creates a new suppression middleware wrapping the given notifier.
// Alerts with the same host+diff fingerprint will be suppressed within the TTL.
func NewMiddleware(next Notifier, s *Suppressor) *Middleware {
	return &Middleware{
		next:       next,
		suppressor: s,
	}
}

// Notify checks whether the event has already been sent recently.
// If it has, the notification is silently suppressed. Otherwise,
// it delegates to the wrapped Notifier.
func (m *Middleware) Notify(ctx context.Context, event alert.Event) error {
	if !m.suppressor.ShouldAlert(event.Host, event.Diff) {
		return nil
	}
	return m.next.Notify(ctx, event)
}
