package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// WebhookNotifier sends port change events to a generic HTTP webhook endpoint.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

// NewWebhookNotifier creates a new WebhookNotifier that posts JSON payloads to url.
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		url:    url,
		client: &http.Client{},
	}
}

type webhookPayload struct {
	Host       string `json:"host"`
	NewPorts   []int  `json:"new_ports"`
	ClosedPorts []int `json:"closed_ports"`
	OpenPorts  []int  `json:"open_ports"`
}

// Notify sends a JSON POST request to the configured webhook URL when there are changes.
func (w *WebhookNotifier) Notify(event alert.Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}

	payload := webhookPayload{
		Host:        event.Host,
		NewPorts:    event.Diff.New,
		ClosedPorts: event.Diff.Closed,
		OpenPorts:   event.OpenPorts,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}

	return nil
}
