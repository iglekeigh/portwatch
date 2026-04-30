package escalation

import (
	"testing"
	"time"
)

func newTestTracker(now time.Time) *Tracker {
	t := New(Config{
		WarningAfter:   5 * time.Minute,
		CriticalAfter:  15 * time.Minute,
		EmergencyAfter: 60 * time.Minute,
	})
	t.nowFn = func() time.Time { return now }
	return t
}

func TestEvaluate_NoneBeforeWarning(t *testing.T) {
	now := time.Now()
	tr := newTestTracker(now)

	lvl := tr.Evaluate("host1")
	if lvl != LevelNone {
		t.Fatalf("expected LevelNone, got %s", lvl)
	}
}

func TestEvaluate_WarningTier(t *testing.T) {
	now := time.Now()
	tr := newTestTracker(now)
	tr.Evaluate("host1") // register first-seen

	// advance time past warning threshold
	tr.nowFn = func() time.Time { return now.Add(6 * time.Minute) }
	lvl := tr.Evaluate("host1")
	if lvl != LevelWarning {
		t.Fatalf("expected LevelWarning, got %s", lvl)
	}
}

func TestEvaluate_CriticalTier(t *testing.T) {
	now := time.Now()
	tr := newTestTracker(now)
	tr.Evaluate("host1")

	tr.nowFn = func() time.Time { return now.Add(20 * time.Minute) }
	lvl := tr.Evaluate("host1")
	if lvl != LevelCritical {
		t.Fatalf("expected LevelCritical, got %s", lvl)
	}
}

func TestEvaluate_EmergencyTier(t *testing.T) {
	now := time.Now()
	tr := newTestTracker(now)
	tr.Evaluate("host1")

	tr.nowFn = func() time.Time { return now.Add(90 * time.Minute) }
	lvl := tr.Evaluate("host1")
	if lvl != LevelEmergency {
		t.Fatalf("expected LevelEmergency, got %s", lvl)
	}
}

func TestResolve_ClearsState(t *testing.T) {
	now := time.Now()
	tr := newTestTracker(now)
	tr.Evaluate("host1")
	tr.Resolve("host1")

	if len(tr.Hosts()) != 0 {
		t.Fatal("expected no tracked hosts after resolve")
	}

	// re-evaluating after resolve should restart the clock
	tr.nowFn = func() time.Time { return now.Add(90 * time.Minute) }
	lvl := tr.Evaluate("host1")
	// first-seen is now reset to the current (advanced) time, so elapsed = 0
	if lvl != LevelNone {
		t.Fatalf("expected LevelNone after resolve, got %s", lvl)
	}
}

func TestHosts_ReturnTracked(t *testing.T) {
	now := time.Now()
	tr := newTestTracker(now)
	tr.Evaluate("alpha")
	tr.Evaluate("beta")

	hosts := tr.Hosts()
	if len(hosts) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(hosts))
	}
}

func TestLevelString(t *testing.T) {
	cases := map[Level]string{
		LevelNone:      "none",
		LevelWarning:   "warning",
		LevelCritical:  "critical",
		LevelEmergency: "emergency",
	}
	for lvl, want := range cases {
		if lvl.String() != want {
			t.Errorf("Level(%d).String() = %q, want %q", lvl, lvl.String(), want)
		}
	}
}
