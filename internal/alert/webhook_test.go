package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
	diff := scanner.Diff{New: []int{8080}, Closed: []int{}}

	if err := n.Notify("localhost", diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(received.NewPorts) != 1 || received.NewPorts[0] != 8080 {
		t.Errorf("expected new port 8080, got %v", received.NewPorts)
	}
	if received.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", received.Host)
	}
}

func TestWebhookNotifier_NoRequestWhenNoDiff(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewWebhookNotifier(ts.URL)
	diff := scanner.Diff{New: []int{}, Closed: []int{}}

	if err := n.Notify("localhost", diff); err != nil {
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
	diff := scanner.Diff{New: []int{443}, Closed: []int{}}

	if err := n.Notify("host", diff); err == nil {
		t.Error("expected error on non-2xx status")
	}
}
