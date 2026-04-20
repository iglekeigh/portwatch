package ratelimit

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Notifier is the interface satisfied by alert notifiers.
type Notifier interface {
	Notify(event alert.Event) error
}

// Middleware wraps a Notifier and suppresses duplicate notifications
// for the same host within a configurable cooldown window.
type Middleware struct {
	inner   Notifier
	limiter *RateLimiter
}

// NewMiddleware returns a Middleware that rate-limits calls to inner
// using the given cooldown duration.
func NewMiddleware(inner Notifier, cooldown time.Duration) *Middleware {
	return &Middleware{
		inner:   inner,
		limiter: New(cooldown),
	}
}

// Notify forwards the event to the inner Notifier only when the rate
// limiter permits it for the event's host. It returns an error if the
// inner Notifier fails, or a sentinel error when the event is suppressed.
func (m *Middleware) Notify(event alert.Event) error {
	if !m.limiter.Allow(event.Host) {
		return fmt.Errorf("ratelimit: notification suppressed for host %q (cooldown active)", event.Host)
	}
	return m.inner.Notify(event)
}

// Reset clears the rate-limit state for the given host, allowing the
// next notification to pass through immediately.
func (m *Middleware) Reset(host string) {
	m.limiter.Reset(host)
}
