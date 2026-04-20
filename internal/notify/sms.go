package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// SMSNotifier sends port change alerts via an SMS gateway (e.g. Twilio-compatible REST API).
type SMSNotifier struct {
	gatewayURL string
	apiKey     string
	from       string
	to         string
	client     *http.Client
}

type smsPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}

// NewSMSNotifier creates a new SMSNotifier.
func NewSMSNotifier(gatewayURL, apiKey, from, to string) *SMSNotifier {
	return &SMSNotifier{
		gatewayURL: gatewayURL,
		apiKey:     apiKey,
		from:       from,
		to:         to,
		client:     &http.Client{},
	}
}

// Notify sends an SMS if there are port changes in the event.
func (s *SMSNotifier) Notify(event alert.Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}

	msg := formatSMSMessage(event)
	payload := smsPayload{From: s.from, To: s.to, Message: msg}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("sms: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.gatewayURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("sms: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("sms: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("sms: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatSMSMessage(event alert.Event) string {
	msg := fmt.Sprintf("[portwatch] %s port changes:", event.Host)
	if len(event.Diff.New) > 0 {
		msg += fmt.Sprintf(" NEW=%v", event.Diff.New)
	}
	if len(event.Diff.Closed) > 0 {
		msg += fmt.Sprintf(" CLOSED=%v", event.Diff.Closed)
	}
	return msg
}
