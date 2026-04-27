package rollup

import (
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// ScanFunc is the signature of the next scan step in the pipeline.
type ScanFunc func(host string, ports []int) (scanner.Result, error)

// Middleware wraps a Notifier with rollup windowing so that rapid successive
// scan cycles for the same host produce at most one notification per window.
type Middleware struct {
	window *Window
}

// NewMiddleware creates a Middleware backed by a rollup Window of the given
// duration forwarding to next.
func NewMiddleware(duration time.Duration, next Notifier) *Middleware {
	return &Middleware{
		window: New(duration, next),
	}
}

// Notify forwards the event into the rollup window.
func (m *Middleware) Notify(event alert.Event) error {
	return m.window.Notify(event)
}
