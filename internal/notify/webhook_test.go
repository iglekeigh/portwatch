package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func TestWebhookNotifier_SendsOnChanges(t *testing.T) {
	var received webhookPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewWebhookNotifier(ts.URL)
	event := alert.Event{
		Host:      "example.com",
		OpenPorts: []int{80, 443, 8080},
		Diff: scanner.Diff{
			New:    []int{8080},
			Closed: []int{22},
		},
	}

	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Host != "example.com" {
		t.Errorf("expected host example.com, got %s", received.Host)
	}
	if len(received.NewPorts) != 1 || received.NewPorts[0] != 8080 {
		t.Errorf("expected new port 8080, got %v", received.NewPorts)
	}
	if len(received.ClosedPorts) != 1 || received.ClosedPorts[0] != 22 {
		t.Errorf("expected closed port 22, got %v", received.ClosedPorts)
	}
}

func TestWebhookNotifier_NoRequestWhenNoDiff(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewWebhookNotifier(ts.URL)
	event := alert.Event{
		Host:      "example.com",
		OpenPorts: []int{80},
		Diff:      scanner.Diff{},
	}

	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request when diff has no changes")
	}
}

func TestWebhookNotifier_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewWebhookNotifier(ts.URL)
	event := alert.Event{
		Host:      "example.com",
		OpenPorts: []int{9000},
		Diff: scanner.Diff{
			New: []int{9000},
		},
	}

	if err := n.Notify(event); err == nil {
		t.Error("expected error on non-2xx status, got nil")
	}
}
