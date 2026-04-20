package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func TestTelegramNotifier_SendsOnChanges(t *testing.T) {
	var received telegramPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewTelegramNotifier("testtoken", "123456")
	n.baseURL = ts.URL + "/bot"

	event := alert.Event{
		Host: "localhost",
		Diff: scanner.Diff{New: []int{80, 443}, Closed: []int{}},
	}

	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.ChatID != "123456" {
		t.Errorf("expected chat_id 123456, got %s", received.ChatID)
	}
	if received.Text == "" {
		t.Error("expected non-empty message text")
	}
}

func TestTelegramNotifier_NoRequestWhenNoDiff(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	n := NewTelegramNotifier("token", "chat")
	n.baseURL = ts.URL + "/bot"

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

func TestTelegramNotifier_ErrorOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewTelegramNotifier("badtoken", "chat")
	n.baseURL = ts.URL + "/bot"

	event := alert.Event{
		Host: "host",
		Diff: scanner.Diff{New: []int{22}},
	}
	if err := n.Notify(event); err == nil {
		t.Error("expected error on non-2xx status")
	}
}

func TestFormatTelegramMessage(t *testing.T) {
	event := alert.Event{
		Host: "myhost",
		Diff: scanner.Diff{New: []int{8080}, Closed: []int{22}},
	}
	msg := formatTelegramMessage(event)
	if msg == "" {
		t.Error("expected non-empty message")
	}
	for _, want := range []string{"myhost", "8080", "22"} {
		if !containsStr(msg, want) {
			t.Errorf("expected message to contain %q", want)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
