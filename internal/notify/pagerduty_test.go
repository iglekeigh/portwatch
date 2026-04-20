package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func TestPagerDutyNotifier_SendsOnChanges(t *testing.T) {
	var received pdPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(202)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("test-key")
	n.endpoint = ts.URL

	event := alert.Event{
		Host: "host1",
		Diff: scanner.Diff{New: []int{8080}, Closed: []int{}},
	}

	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.RoutingKey != "test-key" {
		t.Errorf("expected routing key 'test-key', got %q", received.RoutingKey)
	}
	if received.EventAction != "trigger" {
		t.Errorf("expected event_action 'trigger', got %q", received.EventAction)
	}
	if received.Payload.Source != "host1" {
		t.Errorf("expected source 'host1', got %q", received.Payload.Source)
	}
}

func TestPagerDutyNotifier_NoRequestWhenNoDiff(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(202)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("test-key")
	n.endpoint = ts.URL

	event := alert.Event{
		Host: "host1",
		Diff: scanner.Diff{},
	}

	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request when no diff")
	}
}

func TestPagerDutyNotifier_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("test-key")
	n.endpoint = ts.URL

	event := alert.Event{
		Host: "host1",
		Diff: scanner.Diff{New: []int{443}},
	}

	if err := n.Notify(event); err == nil {
		t.Error("expected error on bad status")
	}
}
