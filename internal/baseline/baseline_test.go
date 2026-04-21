package baseline_test

import (
	"testing"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/scanner"
)

// memStore is an in-memory Store for testing.
type memStore struct {
	data map[string]*baseline.Snapshot
}

func newMem() *memStore { return &memStore{data: map[string]*baseline.Snapshot{}} }

func (m *memStore) Save(host string, snap *baseline.Snapshot) error {
	m.data[host] = snap
	return nil
}

func (m *memStore) Load(host string) (*baseline.Snapshot, error) {
	return m.data[host], nil
}

func TestCapture_StoresSnapshot(t *testing.T) {
	mgr := baseline.New(newMem())
	if err := mgr.Capture("host1", []int{80, 443, 22}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snap, err := mgr.Get("host1")
	if err != nil || snap == nil {
		t.Fatalf("expected snapshot, got err=%v snap=%v", err, snap)
	}
	if len(snap.Ports) != 3 {
		t.Errorf("expected 3 ports, got %d", len(snap.Ports))
	}
}

func TestDiff_NewPorts(t *testing.T) {
	mgr := baseline.New(newMem())
	_ = mgr.Capture("h", []int{80, 443})
	newP, missing, err := mgr.Diff("h", []int{80, 443, 8080})
	if err != nil {
		t.Fatal(err)
	}
	if len(newP) != 1 || newP[0] != 8080 {
		t.Errorf("expected [8080] new, got %v", newP)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing, got %v", missing)
	}
}

func TestDiff_MissingPorts(t *testing.T) {
	mgr := baseline.New(newMem())
	_ = mgr.Capture("h", []int{80, 443, 22})
	_, missing, err := mgr.Diff("h", []int{80, 443})
	if err != nil {
		t.Fatal(err)
	}
	if len(missing) != 1 || missing[0] != 22 {
		t.Errorf("expected [22] missing, got %v", missing)
	}
}

func TestDiff_NoBaseline_ReturnsNil(t *testing.T) {
	mgr := baseline.New(newMem())
	newP, missing, err := mgr.Diff("unknown", []int{80})
	if err != nil {
		t.Fatal(err)
	}
	if newP != nil || missing != nil {
		t.Errorf("expected nil slices, got new=%v missing=%v", newP, missing)
	}
}

func TestCaptureFromResult(t *testing.T) {
	mgr := baseline.New(newMem())
	r := scanner.Result{Host: "srv", Ports: []int{22, 80}}
	if err := mgr.CaptureFromResult(r); err != nil {
		t.Fatal(err)
	}
	snap, _ := mgr.Get("srv")
	if snap == nil || snap.Host != "srv" {
		t.Errorf("unexpected snapshot: %+v", snap)
	}
}
