package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
)

func TestCompute_DeterministicAndSorted(t *testing.T) {
	a := fingerprint.Compute("localhost", []int{443, 80, 22})
	b := fingerprint.Compute("localhost", []int{22, 80, 443})

	if a.Hash != b.Hash {
		t.Errorf("expected same hash for same ports in different order, got %s vs %s", a.Hash, b.Hash)
	}
}

func TestCompute_DifferentPortsDifferentHash(t *testing.T) {
	a := fingerprint.Compute("host1", []int{80})
	b := fingerprint.Compute("host1", []int{443})

	if a.Hash == b.Hash {
		t.Error("expected different hashes for different ports")
	}
}

func TestCompute_DifferentHostsDifferentHash(t *testing.T) {
	a := fingerprint.Compute("host1", []int{80})
	b := fingerprint.Compute("host2", []int{80})

	if a.Hash == b.Hash {
		t.Error("expected different hashes for different hosts")
	}
}

func TestEqual_SameState(t *testing.T) {
	a := fingerprint.Compute("host", []int{22, 80})
	b := fingerprint.Compute("host", []int{22, 80})

	if !fingerprint.Equal(a, b) {
		t.Error("expected Equal to return true for identical state")
	}
}

func TestEqual_DifferentState(t *testing.T) {
	a := fingerprint.Compute("host", []int{22})
	b := fingerprint.Compute("host", []int{22, 80})

	if fingerprint.Equal(a, b) {
		t.Error("expected Equal to return false for different state")
	}
}

func TestChanged_DetectsChange(t *testing.T) {
	prev := fingerprint.Compute("host", []int{80})
	next := fingerprint.Compute("host", []int{80, 443})

	if !fingerprint.Changed(prev, next) {
		t.Error("expected Changed to return true")
	}
}

func TestChanged_NoChange(t *testing.T) {
	prev := fingerprint.Compute("host", []int{80, 443})
	next := fingerprint.Compute("host", []int{80, 443})

	if fingerprint.Changed(prev, next) {
		t.Error("expected Changed to return false")
	}
}

func TestChanged_DifferentHostsReturnFalse(t *testing.T) {
	prev := fingerprint.Compute("host1", []int{80})
	next := fingerprint.Compute("host2", []int{443})

	if fingerprint.Changed(prev, next) {
		t.Error("expected Changed to return false for different hosts")
	}
}

func TestShort_TruncatesHash(t *testing.T) {
	r := fingerprint.Compute("host", []int{80})
	short := r.Short()

	if len(short) != 12 {
		t.Errorf("expected Short() to return 12 chars, got %d", len(short))
	}
	if short != r.Hash[:12] {
		t.Errorf("expected Short() to be prefix of Hash")
	}
}
