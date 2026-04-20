package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const telegramAPIBase = "https://api.telegram.org/bot"

// TelegramNotifier sends alerts to a Telegram chat via Bot API.
type TelegramNotifier struct {
	token  string
	chatID string
	client *http.Client
	baseURL string
}

type telegramPayload struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// NewTelegramNotifier creates a TelegramNotifier with the given bot token and chat ID.
func NewTelegramNotifier(token, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		token:   token,
		chatID:  chatID,
		client:  &http.Client{},
		baseURL: telegramAPIBase,
	}
}

// Notify sends a Telegram message if there are port changes.
func (t *TelegramNotifier) Notify(event alert.Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}

	msg := formatTelegramMessage(event)
	payload := telegramPayload{
		ChatID:    t.chatID,
		Text:      msg,
		ParseMode: "Markdown",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s%s/sendMessage", t.baseURL, t.token)
	resp, err := t.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatTelegramMessage(event alert.Event) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "*Port change detected on %s*\n", event.Host)
	if len(event.Diff.New) > 0 {
		fmt.Fprintf(&buf, "🟢 New ports: %v\n", event.Diff.New)
	}
	if len(event.Diff.Closed) > 0 {
		fmt.Fprintf(&buf, "🔴 Closed ports: %v\n", event.Diff.Closed)
	}
	return buf.String()
}
