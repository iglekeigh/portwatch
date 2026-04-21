package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func TestPushoverNotifier_SendsOnChanges(t *testing.T) {
	var received map[string]string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewPushoverNotifier("tok123", "user456")
	n.apiURL = ts.URL

	event := alert.Event{
		Host: "host1",
		Diff: scanner.Diff{
			New:    []int{8080, 9090},
			Closed: []int{},
		},
	}

	if err := n.Notify(event); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if received["token"] != "tok123" {
		t.Errorf("expected token tok123, got %q", received["token"])
	}
	if received["user"] != "user456" {
		t.Errorf("expected user user456, got %q", received["user"])
	}
	if received["title"] == "" {
		t.Error("expected non-empty title")
	}
	if received["message"] == "" {
		t.Error("expected non-empty message")
	}
}

func TestPushoverNotifier_NoRequestWhenNoDiff(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewPushoverNotifier("tok", "usr")
	n.apiURL = ts.URL

	event := alert.Event{
		Host: "host1",
		Diff: scanner.Diff{},
	}

	if err := n.Notify(event); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if called {
		t.Error("expected no HTTP request when there are no changes")
	}
}

func TestPushoverNotifier_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	n := NewPushoverNotifier("tok", "usr")
	n.apiURL = ts.URL

	event := alert.Event{
		Host: "host1",
		Diff: scanner.Diff{
			New: []int{443},
		},
	}

	if err := n.Notify(event); err == nil {
		t.Error("expected error on bad status, got nil")
	}
}

func TestFormatPushoverMessage_NewAndClosed(t *testing.T) {
	event := alert.Event{
		Host: "myhost",
		Diff: scanner.Diff{
			New:    []int{22, 80},
			Closed: []int{8080},
		},
	}
	msg := formatPushoverMessage(event)
	if msg == "" {
		t.Error("expected non-empty message")
	}
}
