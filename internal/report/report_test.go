package report_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/report"
	"github.com/user/portwatch/internal/scanner"
)

func fixedSummary() report.Summary {
	return report.Summary{
		Host:      "localhost",
		ScannedAt: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
		OpenPorts: []int{80, 443},
		NewPorts:  []int{443},
		Closed:    []int{22},
	}
}

func TestBuild(t *testing.T) {
	diff := scanner.Diff{New: []int{8080}, Closed: []int{22}}
	s := report.Build("host1", []int{80, 8080}, diff)
	if s.Host != "host1" {
		t.Errorf("expected host1, got %s", s.Host)
	}
	if len(s.NewPorts) != 1 || s.NewPorts[0] != 8080 {
		t.Errorf("unexpected new ports: %v", s.NewPorts)
	}
}

func TestTextFormatter(t *testing.T) {
	f := report.NewTextFormatter()
	var buf bytes.Buffer
	if err := f.Format(&buf, fixedSummary()); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "localhost") {
		t.Error("expected host in output")
	}
	if !strings.Contains(out, "80, 443") {
		t.Error("expected open ports in output")
	}
	if !strings.Contains(out, "(none)") == false {
		// closed is 22, so (none) should NOT appear for closed
	}
}

func TestTextFormatter_NoPorts(t *testing.T) {
	f := report.NewTextFormatter()
	s := report.Summary{Host: "h", OpenPorts: nil, NewPorts: nil, Closed: nil}
	var buf bytes.Buffer
	if err := f.Format(&buf, s); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "(none)") {
		t.Error("expected (none) for empty slices")
	}
}
