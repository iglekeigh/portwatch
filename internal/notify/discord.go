package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// DiscordNotifier sends alerts to a Discord channel via webhook.
type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewDiscordNotifier creates a new DiscordNotifier.
func NewDiscordNotifier(webhookURL string, client *http.Client) *DiscordNotifier {
	if client == nil {
		client = &http.Client{}
	}
	return &DiscordNotifier{webhookURL: webhookURL, client: client}
}

type discordPayload struct {
	Content string `json:"content"`
}

// Notify sends a Discord message if there are port changes.
func (d *DiscordNotifier) Notify(event alert.Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}

	content := fmt.Sprintf("**Port change on %s**\n", event.Host)
	if len(event.Diff.New) > 0 {
		content += fmt.Sprintf("🟢 New ports: %v\n", event.Diff.New)
	}
	if len(event.Diff.Closed) > 0 {
		content += fmt.Sprintf("🔴 Closed ports: %v\n", event.Diff.Closed)
	}

	payload := discordPayload{Content: content}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord: marshal payload: %w", err)
	}

	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("discord: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord: unexpected status %d", resp.StatusCode)
	}
	return nil
}
