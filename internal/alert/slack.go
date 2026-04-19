package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SlackNotifier sends port change alerts to a Slack webhook URL.
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlackNotifier creates a SlackNotifier targeting the given Slack webhook URL.
func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends a Slack message if the event contains port changes.
func (s *SlackNotifier) Notify(event Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}

	msg := formatSlackMessage(event)
	payload, err := json.Marshal(slackPayload{Text: msg})
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("slack: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatSlackMessage(event Event) string {
	msg := fmt.Sprintf("*Port change detected on %s*\n", event.Host)
	if len(event.Diff.NewPorts) > 0 {
		msg += fmt.Sprintf(":large_green_circle: New ports: %v\n", event.Diff.NewPorts)
	}
	if len(event.Diff.ClosedPorts) > 0 {
		msg += fmt.Sprintf(":red_circle: Closed ports: %v\n", event.Diff.ClosedPorts)
	}
	return msg
}
