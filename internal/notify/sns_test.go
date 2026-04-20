package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func TestSNSNotifier_SendsOnChanges(t *testing.T) {
	var received snsMessage
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	notifier := NewSNSNotifier(ts.URL)
	event := alert.Event{
		Host: "testhost",
		Diff: scanner.Diff{New: []int{8080}, Closed: []int{22}},
	}

	if err := notifier.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Subject == "" {
		t.Error("expected Subject to be set")
	}
	if received.Message == "" {
		t.Error("expected Message to be set")
	}
}

func TestSNSNotifier_NoRequestWhenNoDiff(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	notifier := NewSNSNotifier(ts.URL)
	event := alert.Event{
		Host: "testhost",
		Diff: scanner.Diff{},
	}

	if err := notifier.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request when there are no changes")
	}
}

func TestSNSNotifier_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	notifier := NewSNSNotifier(ts.URL)
	event := alert.Event{
		Host: "testhost",
		Diff: scanner.Diff{New: []int{443}},
	}

	if err := notifier.Notify(event); err == nil {
		t.Error("expected error on non-2xx status")
	}
}

func TestFormatSNSMessage_NewAndClosed(t *testing.T) {
	event := alert.Event{
		Host: "myhost",
		Diff: scanner.Diff{New: []int{80, 443}, Closed: []int{8080}},
	}
	msg := formatSNSMessage(event)
	if msg == "" {
		t.Error("expected non-empty message")
	}
	for _, want := range []string{"myhost", "80", "443", "8080"} {
		if !containsStr(msg, want) {
			t.Errorf("expected message to contain %q, got: %s", want, msg)
		}
	}
}
