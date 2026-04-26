// Package digest computes and compares fingerprints of scan results,
// allowing portwatch to detect whether a result set has changed since
// the last run without storing the full port list twice.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Result holds the computed fingerprint for a single host scan.
type Result struct {
	Host   string
	Ports  []int
	Digest string
}

// Compute returns a deterministic SHA-256 hex digest for the given host
// and port list. The port list is sorted before hashing so that ordering
// differences do not produce different digests.
func Compute(host string, ports []int) Result {
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	parts := make([]string, len(sorted))
	for i, p := range sorted {
		parts[i] = fmt.Sprintf("%d", p)
	}

	raw := host + ":" + strings.Join(parts, ",")
	sum := sha256.Sum256([]byte(raw))
	return Result{
		Host:   host,
		Ports:  sorted,
		Digest: hex.EncodeToString(sum[:]),
	}
}

// Equal reports whether two Results share the same digest.
func Equal(a, b Result) bool {
	return a.Digest == b.Digest
}

// Changed reports whether the current digest differs from a previously
// stored digest string. An empty previous digest is always considered
// changed so that the first run always triggers downstream processing.
func Changed(current Result, previous string) bool {
	if previous == "" {
		return true
	}
	return current.Digest != previous
}
