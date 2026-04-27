// Package fingerprint generates a stable string identity for a scan result,
// combining host, port list, and optional metadata into a reproducible key.
package fingerprint

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Result holds the computed fingerprint and its components.
type Result struct {
	Host  string
	Ports []int
	Hash  string
}

// Compute returns a fingerprint for the given host and port list.
// The port list is sorted before hashing so order does not affect the result.
func Compute(host string, ports []int) Result {
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	strs := make([]string, len(sorted))
	for i, p := range sorted {
		strs[i] = strconv.Itoa(p)
	}

	raw := fmt.Sprintf("%s|%s", host, strings.Join(strs, ","))
	sum := sha256.Sum256([]byte(raw))

	return Result{
		Host:  host,
		Ports: sorted,
		Hash:  fmt.Sprintf("%x", sum),
	}
}

// Equal reports whether two fingerprints represent the same state.
func Equal(a, b Result) bool {
	return a.Hash == b.Hash
}

// Changed reports whether the port state has changed between two scans
// of the same host. It returns false when hosts differ.
func Changed(prev, next Result) bool {
	if prev.Host != next.Host {
		return false
	}
	return prev.Hash != next.Hash
}

// Short returns the first 12 characters of the hash, useful for display.
func (r Result) Short() string {
	if len(r.Hash) < 12 {
		return r.Hash
	}
	return r.Hash[:12]
}
