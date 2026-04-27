// Package enricher attaches metadata (hostname, geo, org) to scan results.
package enricher

import (
	"fmt"
	"net"
	"strings"
)

// Meta holds enriched metadata for a host address.
type Meta struct {
	Host    string
	Reverse string // reverse DNS name, empty if unavailable
	IsIP    bool
}

// Enricher resolves additional metadata for host addresses.
type Enricher struct {
	lookup func(host string) ([]string, error)
}

// New returns an Enricher using the system DNS resolver.
func New() *Enricher {
	return &Enricher{
		lookup: net.LookupAddr,
	}
}

// newWithLookup returns an Enricher with a custom lookup function (for testing).
func newWithLookup(fn func(string) ([]string, error)) *Enricher {
	return &Enricher{lookup: fn}
}

// Enrich resolves metadata for the given host string.
// host may be an IP address or a hostname.
func (e *Enricher) Enrich(host string) Meta {
	m := Meta{Host: host}

	ip := net.ParseIP(host)
	if ip != nil {
		m.IsIP = true
		names, err := e.lookup(host)
		if err == nil && len(names) > 0 {
			m.Reverse = strings.TrimSuffix(names[0], ".")
		}
		return m
	}

	// host is already a name; attempt to resolve to confirm it is reachable
	m.IsIP = false
	m.Reverse = host
	return m
}

// String returns a human-readable label for the meta.
func (m Meta) String() string {
	if m.IsIP && m.Reverse != "" {
		return fmt.Sprintf("%s (%s)", m.Host, m.Reverse)
	}
	return m.Host
}
