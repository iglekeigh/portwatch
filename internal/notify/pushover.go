package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

const pushoverAPIURL = "https://api.pushover.net/1/messages.json"

// PushoverNotifier sends notifications via the Pushover API.
type PushoverNotifier struct {
	token   string
	userKey string
	client  *http.Client
	apiURL  string
}

// NewPushoverNotifier creates a new Pushover notifier.
func NewPushoverNotifier(token, userKey string) *PushoverNotifier {
	return &PushoverNotifier{
		token:   token,
		userKey: userKey,
		client:  &http.Client{},
		apiURL:  pushoverAPIURL,
	}
}

// Notify sends a Pushover notification if the event contains port changes.
func (p *PushoverNotifier) Notify(event alert.Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}

	payload := map[string]string{
		"token":   p.token,
		"user":    p.userKey,
		"title":   fmt.Sprintf("Port change detected on %s", event.Host),
		"message": formatPushoverMessage(event),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pushover: marshal payload: %w", err)
	}

	resp, err := p.client.Post(p.apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("pushover: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pushover: unexpected status %d", resp.StatusCode)
	}

	return nil
}

func formatPushoverMessage(event alert.Event) string {
	var sb strings.Builder
	if len(event.Diff.New) > 0 {
		fmt.Fprintf(&sb, "New ports: %v\n", event.Diff.New)
	}
	if len(event.Diff.Closed) > 0 {
		fmt.Fprintf(&sb, "Closed ports: %v\n", event.Diff.Closed)
	}
	return strings.TrimSpace(sb.String())
}
