package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// SNSNotifier sends port change alerts to an AWS SNS-compatible HTTP endpoint.
type SNSNotifier struct {
	topicURL string
	client   *http.Client
}

type snsMessage struct {
	TopicArn string `json:"TopicArn"`
	Message  string `json:"Message"`
	Subject  string `json:"Subject"`
}

// NewSNSNotifier creates a new SNSNotifier that posts to the given topic URL.
func NewSNSNotifier(topicURL string) *SNSNotifier {
	return &SNSNotifier{
		topicURL: topicURL,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends an SNS message if there are port changes in the event.
func (s *SNSNotifier) Notify(event alert.Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}

	body, err := json.Marshal(snsMessage{
		TopicArn: s.topicURL,
		Subject:  fmt.Sprintf("portwatch: port changes on %s", event.Host),
		Message:  formatSNSMessage(event),
	})
	if err != nil {
		return fmt.Errorf("sns: marshal error: %w", err)
	}

	resp, err := s.client.Post(s.topicURL, "application/json", strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("sns: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("sns: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatSNSMessage(event alert.Event) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Host: %s\n", event.Host))
	if len(event.Diff.New) > 0 {
		sb.WriteString(fmt.Sprintf("New ports: %v\n", event.Diff.New))
	}
	if len(event.Diff.Closed) > 0 {
		sb.WriteString(fmt.Sprintf("Closed ports: %v\n", event.Diff.Closed))
	}
	return sb.String()
}
