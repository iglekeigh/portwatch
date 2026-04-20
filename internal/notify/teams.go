package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// TeamsNotifier sends alerts to a Microsoft Teams channel via incoming webhook.
type TeamsNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewTeamsNotifier creates a new TeamsNotifier.
func NewTeamsNotifier(webhookURL string, client *http.Client) *TeamsNotifier {
	if client == nil {
		client = &http.Client{}
	}
	return &TeamsNotifier{webhookURL: webhookURL, client: client}
}

type teamsPayload struct {
	Type    string         `json:"@type"`
	Context string         `json:"@context"`
	Text    string         `json:"text"`
	Sections []teamsSection `json:"sections,omitempty"`
}

type teamsSection struct {
	ActivityTitle string `json:"activityTitle"`
	ActivityText  string `json:"activityText"`
}

// Notify sends a Teams message if there are port changes.
func (t *TeamsNotifier) Notify(event alert.Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}

	body := formatTeamsMessage(event)
	payload := teamsPayload{
		Type:    "MessageCard",
		Context: "http://schema.org/extensions",
		Text:    fmt.Sprintf("Port change detected on **%s**", event.Host),
		Sections: []teamsSection{
			{ActivityTitle: "Details", ActivityText: body},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: marshal payload: %w", err)
	}

	resp, err := t.client.Post(t.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("teams: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatTeamsMessage(event alert.Event) string {
	var buf bytes.Buffer
	if len(event.Diff.New) > 0 {
		fmt.Fprintf(&buf, "New ports: %v  ", event.Diff.New)
	}
	if len(event.Diff.Closed) > 0 {
		fmt.Fprintf(&buf, "Closed ports: %v", event.Diff.Closed)
	}
	return buf.String()
}
