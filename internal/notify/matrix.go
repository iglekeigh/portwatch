package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// MatrixNotifier sends port change alerts to a Matrix room via the Client-Server API.
type MatrixNotifier struct {
	homeserver string
	token      string
	roomID     string
	client     *http.Client
}

type matrixTextMessage struct {
	MsgType string `json:"msgtype"`
	Body    string `json:"body"`
}

// NewMatrixNotifier creates a notifier that posts messages to a Matrix room.
// homeserver should be the base URL, e.g. "https://matrix.example.org".
func NewMatrixNotifier(homeserver, token, roomID string) *MatrixNotifier {
	return &MatrixNotifier{
		homeserver: homeserver,
		token:      token,
		roomID:     roomID,
		client:     &http.Client{},
	}
}

// Notify sends a Matrix message if the event contains port changes.
func (m *MatrixNotifier) Notify(event alert.Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}

	body := formatMatrixMessage(event)
	payload := matrixTextMessage{
		MsgType: "m.text",
		Body:    body,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("matrix: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/send/m.room.message", m.homeserver, m.roomID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("matrix: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.token)

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("matrix: send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("matrix: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatMatrixMessage(event alert.Event) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "[portwatch] Port changes detected on %s\n", event.Host)
	if len(event.Diff.New) > 0 {
		fmt.Fprintf(&buf, "  New ports: %v\n", event.Diff.New)
	}
	if len(event.Diff.Closed) > 0 {
		fmt.Fprintf(&buf, "  Closed ports: %v\n", event.Diff.Closed)
	}
	return buf.String()
}
