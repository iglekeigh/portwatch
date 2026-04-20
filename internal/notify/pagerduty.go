package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const pagerDutyEventURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDutyNotifier sends alerts to PagerDuty via Events API v2.
type PagerDutyNotifier struct {
	integrationKey string
	client         *http.Client
	endpoint       string
}

type pdPayload struct {
	RoutingKey  string    `json:"routing_key"`
	EventAction string    `json:"event_action"`
	Payload     pdDetails `json:"payload"`
}

type pdDetails struct {
	Summary  string `json:"summary"`
	Source   string `json:"source"`
	Severity string `json:"severity"`
}

// NewPagerDutyNotifier creates a new PagerDutyNotifier.
func NewPagerDutyNotifier(integrationKey string) *PagerDutyNotifier {
	return &PagerDutyNotifier{
		integrationKey: integrationKey,
		client:         &http.Client{},
		endpoint:       pagerDutyEventURL,
	}
}

// Notify sends a PagerDuty event if there are port changes.
func (p *PagerDutyNotifier) Notify(event alert.Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}

	summary := fmt.Sprintf("Port change detected on %s: %d new, %d closed",
		event.Host, len(event.Diff.New), len(event.Diff.Closed))

	body := pdPayload{
		RoutingKey:  p.integrationKey,
		EventAction: "trigger",
		Payload: pdDetails{
			Summary:  summary,
			Source:   event.Host,
			Severity: "warning",
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("pagerduty: marshal error: %w", err)
	}

	resp, err := p.client.Post(p.endpoint, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty: request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}
