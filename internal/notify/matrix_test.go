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

func TestMatrixNotifier_SendsOnChanges(t *testing.T) {
	var received matrixTextMessage
	var authHeader string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"event_id":"$abc123"}`))
	}))
	defer ts.Close()

	n := NewMatrixNotifier(ts.URL, "secret-token", "!room:example.org")
	event := alert.Event{
		Host: "host1",
		Diff: scanner.Diff{New: []int{8080}, Closed: []int{22}},
	}

	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.MsgType != "m.text" {
		t.Errorf("expected msgtype m.text, got %q", received.MsgType)
	}
	if !strings.Contains(received.Body, "host1") {
		t.Errorf("expected body to contain host, got %q", received.Body)
	}
	if authHeader != "Bearer secret-token" {
		t.Errorf("expected bearer token header, got %q", authHeader)
	}
}

func TestMatrixNotifier_NoRequestWhenNoDiff(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewMatrixNotifier(ts.URL, "token", "!room:example.org")
	event := alert.Event{Host: "host1", Diff: scanner.Diff{}}

	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request when diff is empty")
	}
}

func TestMatrixNotifier_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := NewMatrixNotifier(ts.URL, "bad-token", "!room:example.org")
	event := alert.Event{
		Host: "host1",
		Diff: scanner.Diff{New: []int{443}},
	}

	if err := n.Notify(event); err == nil {
		t.Error("expected error on non-2xx status")
	}
}

func TestFormatMatrixMessage_NewAndClosed(t *testing.T) {
	event := alert.Event{
		Host: "myhost",
		Diff: scanner.Diff{New: []int{80, 443}, Closed: []int{8080}},
	}
	msg := formatMatrixMessage(event)
	if !strings.Contains(msg, "myhost") {
		t.Errorf("expected host in message, got: %q", msg)
	}
	if !strings.Contains(msg, "New ports") {
		t.Errorf("expected new ports label, got: %q", msg)
	}
	if !strings.Contains(msg, "Closed ports") {
		t.Errorf("expected closed ports label, got: %q", msg)
	}
}
