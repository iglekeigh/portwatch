package report

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestMarkdownFormatter_ContainsHost(t *testing.T) {
	f := NewMarkdownFormatter()
	s := fixedSummary()
	out, err := f.Format(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected host in output, got:\n%s", out)
	}
}

func TestMarkdownFormatter_ShowsNewAndClosedPorts(t *testing.T) {
	f := NewMarkdownFormatter()
	s := fixedSummary()
	s.Diff = scanner.Diff{
		New:    []int{8080, 9090},
		Closed: []int{22},
	}
	out, err := f.Format(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"8080", "9090", "22", "## New Ports", "## Closed Ports"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestMarkdownFormatter_EmptyPorts(t *testing.T) {
	f := NewMarkdownFormatter()
	s := fixedSummary()
	s.OpenPorts = nil
	s.Diff = scanner.Diff{}
	out, err := f.Format(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "_None_") {
		t.Errorf("expected '_None_' placeholder in output:\n%s", out)
	}
}

func TestMarkdownFormatter_HeadingStructure(t *testing.T) {
	f := NewMarkdownFormatter()
	s := fixedSummary()
	out, err := f.Format(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, heading := range []string{"# Port Scan Report", "## Open Ports", "## New Ports", "## Closed Ports"} {
		if !strings.Contains(out, heading) {
			t.Errorf("missing heading %q in:\n%s", heading, out)
		}
	}
}
