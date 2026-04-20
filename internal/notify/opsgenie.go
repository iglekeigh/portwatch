package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

const defaultOpsGenieURL = "https://api.opsgenie.com/v2/alerts"

// OpsGenieNotifier sends alerts to OpsGenie.
type OpsGenieNotifier struct {
	apiKey  string
	apiURL  string
	client  *http.Client
}

type opsGeniePayload struct {
	Message     string            `json:"message"`
	Description string            `json:"description"`
	Priority    string            `json:"priority"`
	Details     map[string]string `json:"details"`
}

// NewOpsGenieNotifier creates a new OpsGenieNotifier.
func NewOpsGenieNotifier(apiKey string) *OpsGenieNotifier {
	return &OpsGenieNotifier{
		apiKey: apiKey,
		apiURL: defaultOpsGenieURL,
		client: &http.Client{},
	}
}

// Notify sends an OpsGenie alert if there are port changes.
func (o *OpsGenieNotifier) Notify(event alert.Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}

	payload := opsGeniePayload{
		Message:     fmt.Sprintf("Port change detected on %s", event.Host),
		Description: fmt.Sprintf("New: %v | Closed: %v", event.Diff.New, event.Diff.Closed),
		Priority:    "P3",
		Details: map[string]string{
			"host":        event.Host,
			"new_ports":   fmt.Sprintf("%v", event.Diff.New),
			"closed_ports": fmt.Sprintf("%v", event.Diff.Closed),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsgenie: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, o.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("opsgenie: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status %d", resp.StatusCode)
	}
	return nil
}
