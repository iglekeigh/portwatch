package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func TestGotifyNotifier_SendsOnChanges(t *testing.T) {
	var received gotifyPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewGotifyNotifier(ts.URL, "apptoken", 5)
	event := alert.Event{
		Host: "remotehost",
		Diff: scanner.Diff{New: []int{3000}, Closed: []int{8080}},
	}

	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Title == "" {
		t.Error("expected non-empty title")
	}
	if received.Priority != 5 {
		t.Errorf("expected priority 5, got %d", received.Priority)
	}
}

func TestGotifyNotifier_NoRequestWhenNoDiff(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewGotifyNotifier(ts.URL, "token", 3)
	event := alert.Event{
		Host: "host",
		Diff: scanner.Diff{},
	}
	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request when no diff")
	}
}

func TestGotifyNotifier_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewGotifyNotifier(ts.URL, "badtoken", 1)
	event := alert.Event{
		Host: "host",
		Diff: scanner.Diff{New: []int{443}},
	}
	if err := n.Notify(event); err == nil {
		t.Error("expected error on non-2xx status")
	}
}
