package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// WebhookNotifier sends port change events to an HTTP endpoint.
type WebhookNotifier struct {
	URL    string
	Client *http.Client
}

// NewWebhookNotifier creates a WebhookNotifier with the given URL.
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		URL: url,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

type webhookPayload struct {
	Host       string    `json:"host"`
	NewPorts   []int     `json:"new_ports"`
	ClosedPorts []int    `json:"closed_ports"`
	Timestamp  time.Time `json:"timestamp"`
}

// Notify sends a POST request with the diff payload if there are changes.
func (w *WebhookNotifier) Notify(host string, diff scanner.Diff) error {
	if !diff.HasChanges() {
		return nil
	}

	payload := webhookPayload{
		Host:        host,
		NewPorts:    diff.New,
		ClosedPorts: diff.Closed,
		Timestamp:   time.Now().UTC(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.Client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post to %s: %w", w.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, w.URL)
	}

	return nil
}
