package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func TestSMSNotifier_SendsOnChanges(t *testing.T) {
	var received smsPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewSMSNotifier(ts.URL, "key", "+10000000000", "+19999999999")
	event := alert.Event{
		Host: "host1",
		Diff: scanner.Diff{New: []int{22, 80}, Closed: []int{}},
	}

	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.To != "+19999999999" {
		t.Errorf("expected to +19999999999, got %s", received.To)
	}
	if received.From != "+10000000000" {
		t.Errorf("expected from +10000000000, got %s", received.From)
	}
	if received.Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestSMSNotifier_NoRequestWhenNoDiff(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewSMSNotifier(ts.URL, "key", "+10000000000", "+19999999999")
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

func TestSMSNotifier_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewSMSNotifier(ts.URL, "key", "+10000000000", "+19999999999")
	event := alert.Event{
		Host: "host1",
		Diff: scanner.Diff{New: []int{443}},
	}

	if err := n.Notify(event); err == nil {
		t.Error("expected error on bad status")
	}
}

func TestFormatSMSMessage_NewAndClosed(t *testing.T) {
	event := alert.Event{
		Host: "myhost",
		Diff: scanner.Diff{New: []int{80}, Closed: []int{22}},
	}
	msg := formatSMSMessage(event)
	if msg == "" {
		t.Error("expected non-empty message")
	}
	for _, sub := range []string{"myhost", "NEW", "CLOSED"} {
		if !containsStr(msg, sub) {
			t.Errorf("expected message to contain %q, got: %s", sub, msg)
		}
	}
}
