package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func TestTeamsNotifier_SendsOnChanges(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewTeamsNotifier(ts.URL, ts.Client())
	event := alert.Event{
		Host: "localhost",
		Diff: scanner.Diff{New: []int{8080}, Closed: []int{}},
	}
	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTeamsNotifier_NoRequestWhenNoDiff(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewTeamsNotifier(ts.URL, ts.Client())
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

func TestTeamsNotifier_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewTeamsNotifier(ts.URL, ts.Client())
	event := alert.Event{
		Host: "localhost",
		Diff: scanner.Diff{New: []int{443}},
	}
	if err := n.Notify(event); err == nil {
		t.Error("expected error on bad status")
	}
}

func TestFormatTeamsMessage_NewAndClosed(t *testing.T) {
	event := alert.Event{
		Host: "host1",
		Diff: scanner.Diff{New: []int{80, 443}, Closed: []int{8080}},
	}
	msg := formatTeamsMessage(event)
	if msg == "" {
		t.Error("expected non-empty message")
	}
}
