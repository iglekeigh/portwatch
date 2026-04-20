package notify

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// NtfyNotifier sends notifications to an ntfy.sh topic.
type NtfyNotifier struct {
	serverURL string
	topic     string
	client    *http.Client
}

// NewNtfyNotifier creates a notifier that publishes to ntfy.sh or a self-hosted ntfy server.
func NewNtfyNotifier(serverURL, topic string) *NtfyNotifier {
	if serverURL == "" {
		serverURL = "https://ntfy.sh"
	}
	return &NtfyNotifier{
		serverURL: serverURL,
		topic:     topic,
		client:    &http.Client{},
	}
}

// Notify sends a port change event to the configured ntfy topic.
func (n *NtfyNotifier) Notify(event alert.Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}

	body := formatNtfyMessage(event)
	url := fmt.Sprintf("%s/%s", n.serverURL, n.topic)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return fmt.Errorf("ntfy: create request: %w", err)
	}
	req.Header.Set("Title", fmt.Sprintf("Port change on %s", event.Host))
	req.Header.Set("Content-Type", "text/plain")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("ntfy: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatNtfyMessage(event alert.Event) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Host: %s\n", event.Host)
	if len(event.Diff.New) > 0 {
		fmt.Fprintf(&buf, "New ports: %v\n", event.Diff.New)
	}
	if len(event.Diff.Closed) > 0 {
		fmt.Fprintf(&buf, "Closed ports: %v\n", event.Diff.Closed)
	}
	return buf.String()
}
