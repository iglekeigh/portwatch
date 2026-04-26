package digest

import (
	"testing"
)

func TestCompute_DeterministicAndSorted(t *testing.T) {
	a := Compute("localhost", []int{443, 80, 8080})
	b := Compute("localhost", []int{8080, 443, 80})

	if a.Digest != b.Digest {
		t.Errorf("expected same digest for same ports in different order, got %s vs %s", a.Digest, b.Digest)
	}
	if len(a.Ports) != 3 || a.Ports[0] != 80 {
		t.Errorf("expected sorted ports [80 443 8080], got %v", a.Ports)
	}
}

func TestCompute_DifferentPortsDifferentDigest(t *testing.T) {
	a := Compute("host", []int{80})
	b := Compute("host", []int{443})
	if a.Digest == b.Digest {
		t.Error("expected different digests for different ports")
	}
}

func TestCompute_DifferentHostsDifferentDigest(t *testing.T) {
	a := Compute("host-a", []int{80})
	b := Compute("host-b", []int{80})
	if a.Digest == b.Digest {
		t.Error("expected different digests for different hosts")
	}
}

func TestEqual(t *testing.T) {
	a := Compute("h", []int{22, 80})
	b := Compute("h", []int{80, 22})
	if !Equal(a, b) {
		t.Error("expected Equal to return true for identical port sets")
	}
}

func TestChanged_EmptyPrevious(t *testing.T) {
	r := Compute("h", []int{80})
	if !Changed(r, "") {
		t.Error("expected Changed to return true when previous is empty")
	}
}

func TestChanged_SameDigest(t *testing.T) {
	r := Compute("h", []int{80})
	if Changed(r, r.Digest) {
		t.Error("expected Changed to return false when digest matches")
	}
}

func TestChanged_DifferentDigest(t *testing.T) {
	old := Compute("h", []int{80})
	new := Compute("h", []int{443})
	if !Changed(new, old.Digest) {
		t.Error("expected Changed to return true when digest differs")
	}
}
