package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestBuildEvent_NewPorts(t *testing.T) {
	d := scanner.Diff{NewPorts: []int{80, 443}, ClosedPorts: nil}
	e := BuildEvent("localhost", d)
	if e.Level != LevelAlert {
		t.Errorf("expected ALERT, got %s", e.Level)
	}
	if e.Host != "localhost" {
		t.Errorf("unexpected host %s", e.Host)
	}
}

func TestBuildEvent_ClosedPorts(t *testing.T) {
	d := scanner.Diff{NewPorts: nil, ClosedPorts: []int{22}}
	e := BuildEvent("10.0.0.1", d)
	if e.Level != LevelWarn {
		t.Errorf("expected WARN, got %s", e.Level)
	}
}

func TestBuildEvent_NoChanges(t *testing.T) {
	d := scanner.Diff{}
	e := BuildEvent("host", d)
	if e.Level != LevelInfo {
		t.Errorf("expected INFO, got %s", e.Level)
	}
}

func TestConsoleNotifier_Notify(t *testing.T) {
	var buf bytes.Buffer
	n := &ConsoleNotifier{Out: &buf}

	d := scanner.Diff{NewPorts: []int{8080}, ClosedPorts: []int{3000}}
	e := BuildEvent("myhost", d)
	if err := n.Notify(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "NEW ports opened") {
		t.Errorf("expected new ports line, got: %s", out)
	}
	if !strings.Contains(out, "ports CLOSED") {
		t.Errorf("expected closed ports line, got: %s", out)
	}
	if !strings.Contains(out, "myhost") {
		t.Errorf("expected host in output, got: %s", out)
	}
}

func TestConsoleNotifier_NoOutput_WhenNoDiff(t *testing.T) {
	var buf bytes.Buffer
	n := &ConsoleNotifier{Out: &buf}
	e := BuildEvent("host", scanner.Diff{})
	_ = n.Notify(e)
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}
