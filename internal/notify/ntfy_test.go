package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com	/user/portwatch/internal/scanner"
)

func TestNtfyNotifier_SendsOnChanges(t *testing.T) {
	received := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = true
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewNtfyNotifier(ts.URL, "portwatch")
	event := alert.BuildEvent("host1", scanner.Diff{New: []int{80}, Closed: []int{}})
	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !received {
		t.Error("expected HTTP request to be sent")
	}
}

func TestNtfyNotifier_NoRequestWhenNoDiff(t *testing.T) {
	sent := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sent = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewNtfyNotifier(ts.URL, "portwatch")
	event := alert.BuildEvent("host1", scanner.Diff{})
	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sent {
		t.Error("expected no HTTP request when there is no diff")
	}
}

func TestNtfyNotifier_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewNtfyNotifier(ts.URL, "portwatch")
	event := alert.BuildEvent("host1", scanner.Diff{New: []int{443}})
	if err := n.Notify(event); err == nil {
		t.Error("expected error on non-2xx status")
	}
}

func TestFormatNtfyMessage_NewAndClosed(t *testing.T) {
	event := alert.BuildEvent("myhost", scanner.Diff{New: []int{22, 80}, Closed: []int{8080}})
	msg := formatNtfyMessage(event)
	if !containsStr(msg, "myhost") {
		t.Error("expected host in message")
	}
	if !containsStr(msg, "New ports") {
		t.Error("expected new ports section")
	}
	if !containsStr(msg, "Closed ports") {
		t.Error("expected closed ports section")
	}
}
