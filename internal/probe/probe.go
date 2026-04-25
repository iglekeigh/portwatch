// Package probe provides host reachability checks before port scanning.
package probe

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a host probe.
type Result struct {
	Host        string
	Reachable   bool
	Latency     time.Duration
	Error       error
}

// Prober checks whether a host is reachable.
type Prober struct {
	timeout time.Duration
	port    string
}

// New returns a Prober with the given timeout. If timeout is zero, 3 seconds is used.
func New(timeout time.Duration, port string) *Prober {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	if port == "" {
		port = "80"
	}
	return &Prober{timeout: timeout, port: port}
}

// Probe attempts a TCP connection to host:port and returns a Result.
func (p *Prober) Probe(ctx context.Context, host string) Result {
	addr := net.JoinHostPort(host, p.port)
	start := time.Now()

	dialer := &net.Dialer{Timeout: p.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	latency := time.Since(start)

	if err != nil {
		return Result{
			Host:      host,
			Reachable: false,
			Latency:   latency,
			Error:     fmt.Errorf("probe %s: %w", addr, err),
		}
	}
	_ = conn.Close()
	return Result{
		Host:      host,
		Reachable: true,
		Latency:   latency,
	}
}

// ProbeAll probes multiple hosts concurrently and returns all results.
func (p *Prober) ProbeAll(ctx context.Context, hosts []string) []Result {
	results := make([]Result, len(hosts))
	type indexed struct {
		i int
		r Result
	}
	ch := make(chan indexed, len(hosts))
	for i, h := range hosts {
		go func(idx int, host string) {
			ch <- indexed{i: idx, r: p.Probe(ctx, host)}
		}(i, h)
	}
	for range hosts {
		v := <-ch
		results[v.i] = v.r
	}
	return results
}
