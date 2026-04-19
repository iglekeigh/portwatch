package scanner

import (
	"fmt"
	"net"
	"sort"
	"time"
)

// PortResult holds the result of a single port scan.
type PortResult struct {
	Port  int
	Open  bool
}

// ScanOptions configures a port scan.
type ScanOptions struct {
	Host    string
	Ports   []int
	Timeout time.Duration
}

// Scanner performs port scanning.
type Scanner struct {
	opts ScanOptions
}

// New creates a new Scanner with the given options.
func New(opts ScanOptions) *Scanner {
	if opts.Timeout == 0 {
		opts.Timeout = 2 * time.Second
	}
	return &Scanner{opts: opts}
}

// Scan checks all configured ports and returns open ones.
func (s *Scanner) Scan() ([]int, error) {
	results := make(chan PortResult, len(s.opts.Ports))

	for _, port := range s.opts.Ports {
		go func(p int) {
			addr := fmt.Sprintf("%s:%d", s.opts.Host, p)
			conn, err := net.DialTimeout("tcp", addr, s.opts.Timeout)
			if err != nil {
				results <- PortResult{Port: p, Open: false}
				return
			}
			conn.Close()
			results <- PortResult{Port: p, Open: true}
		}(port)
	}

	var open []int
	for range s.opts.Ports {
		r := <-results
		if r.Open {
			open = append(open, r.Port)
		}
	}
	sort.Ints(open)
	return open, nil
}
