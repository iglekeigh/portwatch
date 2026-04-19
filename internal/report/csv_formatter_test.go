package report

import (
	"strings"
	"testing"
)

func TestCSVFormatter_NoChanges(t *testing.T) {
	f := NewCSVFormatter()
	s := Summary{Host: "localhost"}

	out, err := f.Format(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "host,status,ports") {
		t.Error("expected CSV header")
	}
	if !strings.Contains(out, "localhost,no-change") {
		t.Error("expected no-change row")
	}
}

func TestCSVFormatter_NewAndClosedPorts(t *testing.T) {
	f := NewCSVFormatter()
	s := Summary{
		Host:        "192.168.1.1",
		NewPorts:    []int{80, 443},
		ClosedPorts: []int{8080},
		OpenPorts:   []int{22},
	}

	out, err := f.Format(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "192.168.1.1,new,80;443") {
		t.Errorf("expected new ports row, got:\n%s", out)
	}
	if !strings.Contains(out, "192.168.1.1,closed,8080") {
		t.Errorf("expected closed ports row, got:\n%s", out)
	}
	if !strings.Contains(out, "192.168.1.1,open,22") {
		t.Errorf("expected open ports row, got:\n%s", out)
	}
}

func TestCSVFormatter_OnlyOpenPorts(t *testing.T) {
	f := NewCSVFormatter()
	s := Summary{
		Host:      "10.0.0.1",
		OpenPorts: []int{22, 80},
	}

	out, err := f.Format(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "10.0.0.1,open,22;80") {
		t.Errorf("expected open ports row, got:\n%s", out)
	}
}
