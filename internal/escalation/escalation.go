// Package escalation provides tiered alert escalation based on how long
// a port change condition has persisted without acknowledgement.
package escalation

import (
	"sync"
	"time"
)

// Level represents the severity tier of an escalation.
type Level int

const (
	LevelNone   Level = iota
	LevelWarning       // first tier
	LevelCritical      // second tier
	LevelEmergency     // third tier
)

// String returns a human-readable name for the level.
func (l Level) String() string {
	switch l {
	case LevelWarning:
		return "warning"
	case LevelCritical:
		return "critical"
	case LevelEmergency:
		return "emergency"
	default:
		return "none"
	}
}

// Config defines the thresholds for each escalation tier.
type Config struct {
	WarningAfter   time.Duration
	CriticalAfter  time.Duration
	EmergencyAfter time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		WarningAfter:   5 * time.Minute,
		CriticalAfter:  15 * time.Minute,
		EmergencyAfter: 60 * time.Minute,
	}
}

type entry struct {
	first time.Time
	last  Level
}

// Tracker tracks how long a host has been in an unresolved changed state
// and returns the appropriate escalation level.
type Tracker struct {
	mu     sync.Mutex
	cfg    Config
	hosts  map[string]*entry
	nowFn  func() time.Time
}

// New creates a Tracker with the given config.
func New(cfg Config) *Tracker {
	return &Tracker{
		cfg:   cfg,
		hosts: make(map[string]*entry),
		nowFn: time.Now,
	}
}

// Evaluate returns the current escalation Level for host and records the
// first-seen time if this is a new incident. Returns LevelNone if the host
// is not being tracked.
func (t *Tracker) Evaluate(host string) Level {
	t.mu.Lock()
	defer t.mu.Unlock()

	e, ok := t.hosts[host]
	if !ok {
		e = &entry{first: t.nowFn()}
		t.hosts[host] = e
	}

	elapsed := t.nowFn().Sub(e.first)
	var lvl Level
	switch {
	case elapsed >= t.cfg.EmergencyAfter:
		lvl = LevelEmergency
	case elapsed >= t.cfg.CriticalAfter:
		lvl = LevelCritical
	case elapsed >= t.cfg.WarningAfter:
		lvl = LevelWarning
	default:
		lvl = LevelNone
	}
	e.last = lvl
	return lvl
}

// Resolve clears the tracked state for host (e.g. after ports return to baseline).
func (t *Tracker) Resolve(host string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.hosts, host)
}

// Hosts returns all currently tracked host names.
func (t *Tracker) Hosts() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]string, 0, len(t.hosts))
	for h := range t.hosts {
		out = append(out, h)
	}
	return out
}
