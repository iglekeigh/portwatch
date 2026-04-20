package notify

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(newPorts, closedPorts []int) alert.Event {
	return alert.Event{
		Host:      "localhost",
		ScannedAt: time.Now(),
		Diff: scanner.Diff{
			New:    newPorts,
			Closed: closedPorts,
		},
	}
}

func TestBuildEmailBody_NewPorts(t *testing.T) {
	event := makeEvent([]int{80, 443}, nil)
	body := buildEmailBody(event)
	if body == "" {
		t.Fatal("expected non-empty body")
	}
	if !contains(body, "80") || !contains(body, "443") {
		t.Errorf("expected ports in body, got: %s", body)
	}
}

func TestBuildEmailBody_ClosedPorts(t *testing.T) {
	event := makeEvent(nil, []int{22})
	body := buildEmailBody(event)
	if !contains(body, "Closed ports") {
		t.Errorf("expected closed ports section, got: %s", body)
	}
}

func TestBuildEmailBody_ContainsHost(t *testing.T) {
	event := makeEvent([]int{8080}, nil)
	body := buildEmailBody(event)
	if !contains(body, "localhost") {
		t.Errorf("expected host in body, got: %s", body)
	}
}

func TestJoinPorts(t *testing.T) {
	result := joinPorts([]int{22, 80, 443})
	expected := "22, 80, 443"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestEmailNotifier_NoSendWhenNoDiff(t *testing.T) {
	// Notifier should return nil without attempting to send when there are no changes.
	cfg := EmailConfig{
		Host: "127.0.0.1",
		Port: 9999,
		From: "from@example.com",
		To:   []string{"to@example.com"},
	}
	n := NewEmailNotifier(cfg)
	event := makeEvent(nil, nil)
	if err := n.Notify(event); err != nil {
		t.Errorf("expected no error when no diff, got: %v", err)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
