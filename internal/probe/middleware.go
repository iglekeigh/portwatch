package probe

import (
	"context"
	"fmt"
	"time"
)

// ScanFunc is a function that performs a port scan on a host.
type ScanFunc func(ctx context.Context, host string) ([]int, error)

// Middleware wraps a ScanFunc with a reachability probe. If the host is not
// reachable the scan is skipped and an error is returned instead.
type Middleware struct {
	prober  *Prober
	next    ScanFunc
}

// NewMiddleware returns a Middleware that probes host reachability before
// delegating to next. probePort is the TCP port used for the probe check.
func NewMiddleware(timeout time.Duration, probePort string, next ScanFunc) *Middleware {
	return &Middleware{
		prober: New(timeout, probePort),
		next:   next,
	}
}

// Scan probes the host and, if reachable, calls the wrapped ScanFunc.
func (m *Middleware) Scan(ctx context.Context, host string) ([]int, error) {
	res := m.prober.Probe(ctx, host)
	if !res.Reachable {
		return nil, fmt.Errorf("host %s unreachable: %w", host, res.Error)
	}
	return m.next(ctx, host)
}
