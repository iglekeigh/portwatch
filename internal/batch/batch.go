// Package batch provides utilities for scanning multiple hosts concurrently
// and collecting results with bounded parallelism.
package batch

import (
	"context"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Result holds the scan outcome for a single host.
type Result struct {
	Host  string
	Ports []int
	Err   error
}

// Scanner defines the interface for scanning a single host.
type Scanner interface {
	Scan(ctx context.Context, host string) ([]int, error)
}

// Runner executes scans across multiple hosts with a configurable concurrency limit.
type Runner struct {
	scanner     Scanner
	concurrency int
}

// New creates a Runner with the given scanner and concurrency limit.
// If concurrency is <= 0 it defaults to 4.
func New(s Scanner, concurrency int) *Runner {
	if concurrency <= 0 {
		concurrency = 4
	}
	return &Runner{scanner: s, concurrency: concurrency}
}

// Run scans all hosts concurrently (up to the concurrency limit) and returns
// a slice of Result values in the same order as the input hosts.
func (r *Runner) Run(ctx context.Context, hosts []string) []Result {
	results := make([]Result, len(hosts))
	sem := make(chan struct{}, r.concurrency)
	var wg sync.WaitGroup

	for i, host := range hosts {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, h string) {
			defer wg.Done()
			defer func() { <-sem }()

			ports, err := r.scanner.Scan(ctx, h)
			results[idx] = Result{Host: h, Ports: ports, Err: err}
		}(i, host)
	}

	wg.Wait()
	return results
}

// HasErrors returns true if any result contains a non-nil error.
func HasErrors(results []Result) bool {
	for _, r := range results {
		if r.Err != nil {
			return true
		}
	}
	return false
}

// Successful returns only results that completed without error.
func Successful(results []Result) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		if r.Err == nil {
			out = append(out, r)
		}
	}
	return out
}

// ensure scanner.Scanner satisfies our local interface at compile time.
var _ Scanner = (*scanner.Scanner)(nil)
