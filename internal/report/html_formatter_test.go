package report

import (
	"bytes"
	"strings"
	"testing"
)

func TestHTMLFormatter_ContainsHost(t *testing.T) {
	s := fixedSummary()
	f := NewHTMLFormatter()
	var buf bytes.Buffer
	if err := f.Format(s, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, s.Host) {
		t.Errorf("expected host %q in output", s.Host)
	}
	if !strings.Contains(out, "<!DOCTYPE html>") {
		t.Error("expected HTML doctype in output")
	}
}

func TestHTMLFormatter_ShowsOpenPorts(t *testing.T) {
	s := fixedSummary()
	f := NewHTMLFormatter()
	var buf bytes.Buffer
	_ = f.Format(s, &buf)
	out := buf.String()
	for _, p := range s.OpenPorts {
		portStr := joinInts([]int{p})
		if !strings.Contains(out, portStr) {
			t.Errorf("expected port %d in output", p)
		}
	}
}

func TestHTMLFormatter_EmptyPorts(t *testing.T) {
	s := Summary{
		Host:        "empty-host",
		ScannedAt:   "2024-01-01",
		OpenPorts:   []int{},
		NewPorts:    []int{},
		ClosedPorts: []int{},
	}
	f := NewHTMLFormatter()
	var buf bytes.Buffer
	if err := f.Format(s, &buf); err != nil {
		t.Fatalf("unexpected error on empty ports: %v", err)
	}
	if !strings.Contains(buf.String(), "empty-host") {
		t.Error("expected host in empty-ports output")
	}
}
