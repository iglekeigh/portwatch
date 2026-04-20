package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// GotifyNotifier sends alerts to a self-hosted Gotify server.
type GotifyNotifier struct {
	serverURL string
	token     string
	priority  int
	client    *http.Client
}

type gotifyPayload struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

// NewGotifyNotifier creates a GotifyNotifier targeting the given server URL.
func NewGotifyNotifier(serverURL, token string, priority int) *GotifyNotifier {
	return &GotifyNotifier{
		serverURL: serverURL,
		token:     token,
		priority:  priority,
		client:    &http.Client{},
	}
}

// Notify sends a Gotify push notification if there are port changes.
func (g *GotifyNotifier) Notify(event alert.Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}

	var buf bytes.Buffer
	if len(event.Diff.New) > 0 {
		fmt.Fprintf(&buf, "New ports: %v\n", event.Diff.New)
	}
	if len(event.Diff.Closed) > 0 {
		fmt.Fprintf(&buf, "Closed ports: %v\n", event.Diff.Closed)
	}

	payload := gotifyPayload{
		Title:    fmt.Sprintf("Port change on %s", event.Host),
		Message:  buf.String(),
		Priority: g.priority,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gotify: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/message?token=%s", g.serverURL, g.token)
	resp, err := g.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("gotify: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gotify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
