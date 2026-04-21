// Package baseline provides functionality to capture and compare
// port scan snapshots as a trusted reference point.
package baseline

import (
	"fmt"
	"sort"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot represents a trusted baseline of open ports for a host.
type Snapshot struct {
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
	CapturedAt time.Time `json:"captured_at"`
}

// Manager manages baseline snapshots for hosts.
type Manager struct {
	store Store
}

// New returns a new Manager backed by the given Store.
func New(s Store) *Manager {
	return &Manager{store: s}
}

// Capture records the current scan result as the baseline for a host.
func (m *Manager) Capture(host string, ports []int) error {
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	snap := &Snapshot{
		Host:       host,
		Ports:      sorted,
		CapturedAt: time.Now().UTC(),
	}
	return m.store.Save(host, snap)
}

// Get retrieves the baseline snapshot for a host.
// Returns nil, nil if no baseline has been captured.
func (m *Manager) Get(host string) (*Snapshot, error) {
	return m.store.Load(host)
}

// Diff compares a set of current ports against the stored baseline.
// Returns new ports (not in baseline) and missing ports (in baseline but not current).
func (m *Manager) Diff(host string, current []int) (newPorts, missingPorts []int, err error) {
	snap, err := m.store.Load(host)
	if err != nil {
		return nil, nil, fmt.Errorf("baseline: load %s: %w", host, err)
	}
	if snap == nil {
		return nil, nil, nil
	}

	baseSet := toSet(snap.Ports)
	currentSet := toSet(current)

	for _, p := range current {
		if !baseSet[p] {
			newPorts = append(newPorts, p)
		}
	}
	for _, p := range snap.Ports {
		if !currentSet[p] {
			missingPorts = append(missingPorts, p)
		}
	}
	sort.Ints(newPorts)
	sort.Ints(missingPorts)
	return newPorts, missingPorts, nil
}

// CaptureFromResult is a convenience wrapper over scanner.Result.
func (m *Manager) CaptureFromResult(r scanner.Result) error {
	return m.Capture(r.Host, r.Ports)
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
