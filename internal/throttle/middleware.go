package throttle

import (
	"fmt"

	"github.com/user/portwatch/internal/alert"
)

// Notifier is the interface for sending notifications.
type Notifier interface {
	Notify(event alert.Event) error
}

// Middleware wraps a Notifier and suppresses notifications for hosts
// that have been scanned too recently, based on the throttle policy.
type Middleware struct {
	inner    Notifier
	throttle *Throttle
}

// NewMiddleware returns a Middleware that gates notifications through
// the provided Throttle before delegating to inner.
func NewMiddleware(inner Notifier, t *Throttle) *Middleware {
	return &Middleware{inner: inner, throttle: t}
}

// Notify forwards the event to the inner notifier only if the throttle
// permits a notification for the event's host at this time.
func (m *Middleware) Notify(event alert.Event) error {
	if !m.throttle.Allow(event.Host) {
		return fmt.Errorf("throttle: notification suppressed for host %q", event.Host)
	}
	return m.inner.Notify(event)
}
