package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestSlackNotifier_SendsOnChanges(t *testing.T) {
	var received slackPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	notifier := NewSlackNotifier(ts.URL)
	event := Event{
		Host: "localhost",
		Diff: scanner.Diff{
			NewPorts:    []int{8080},
			ClosedPorts: []int{},
		},
	}

	if err := notifier.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Text == "" {
		t.Error("expected non-empty slack message")
	}
}

func TestSlackNotifier_NoRequestWhenNoDiff(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	notifier := NewSlackNotifier(ts.URL)
	event := Event{
		Host: "localhost",
		Diff: scanner.Diff{},
	}

	if err := notifier.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request when there are no changes")
	}
}

func TestSlackNotifier_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	notifier := NewSlackNotifier(ts.URL)
	event := Event{
		Host: "localhost",
		Diff: scanner.Diff{
			NewPorts: []int{443},
		},
	}

	if err := notifier.Notify(event); err == nil {
		t.Error("expected error on non-2xx status")
	}
}
