package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func TestDiscordNotifier_SendsOnChanges(t *testing.T) {
	var got discordPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	n := NewDiscordNotifier(ts.URL, ts.Client())
	event := alert.Event{
		Host: "myhost",
		Diff: scanner.Diff{New: []int{22, 80}},
	}
	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got.Content, "myhost") {
		t.Errorf("expected host in message, got: %s", got.Content)
	}
}

func TestDiscordNotifier_NoRequestWhenNoDiff(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewDiscordNotifier(ts.URL, ts.Client())
	event := alert.Event{Host: "localhost", Diff: scanner.Diff{}}
	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request when no diff")
	}
}

func TestDiscordNotifier_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewDiscordNotifier(ts.URL, ts.Client())
	event := alert.Event{
		Host: "localhost",
		Diff: scanner.Diff{Closed: []int{3306}},
	}
	if err := n.Notify(event); err == nil {
		t.Error("expected error on bad status")
	}
}
