package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func TestOpsGenieNotifier_SendsOnChanges(t *testing.T) {
	var received opsGeniePayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewOpsGenieNotifier("test-key")
	n.apiURL = ts.URL

	event := alert.Event{
		Host: "localhost",
		Diff: scanner.Diff{New: []int{8080}, Closed: []int{}},
	}

	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Message == "" {
		t.Error("expected message to be set")
	}
	if received.Details["host"] != "localhost" {
		t.Errorf("expected host detail, got %q", received.Details["host"])
	}
}

func TestOpsGenieNotifier_NoRequestWhenNoDiff(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewOpsGenieNotifier("test-key")
	n.apiURL = ts.URL

	event := alert.Event{
		Host: "localhost",
		Diff: scanner.Diff{},
	}

	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request when no diff")
	}
}

func TestOpsGenieNotifier_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewOpsGenieNotifier("bad-key")
	n.apiURL = ts.URL

	event := alert.Event{
		Host: "localhost",
		Diff: scanner.Diff{New: []int{443}},
	}

	if err := n.Notify(event); err == nil {
		t.Error("expected error on bad status")
	}
}
