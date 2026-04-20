// Package notify provides integrations for sending port-change alerts
// to external notification services.
package notify

import (
	"fmt"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
)

// Notifier is the common interface implemented by all notification backends.
type Notifier interface {
	Notify(event alert.Event) error
}

// FromConfig constructs a slice of Notifiers based on the provided configuration.
// Unknown notifier types are returned as an error.
func FromConfig(cfg config.Config) ([]Notifier, error) {
	var notifiers []Notifier

	for _, n := range cfg.Notifiers {
		switch n.Type {
		case "email":
			notifiers = append(notifiers, NewEmailNotifier(
				n.Settings["smtp_host"],
				n.Settings["from"],
				n.Settings["to"],
			))
		case "pagerduty":
			notifiers = append(notifiers, NewPagerDutyNotifier(
				n.Settings["routing_key"],
			))
		case "opsgenie":
			notifiers = append(notifiers, NewOpsGenieNotifier(
				n.Settings["api_key"],
			))
		case "teams":
			notifiers = append(notifiers, NewTeamsNotifier(
				n.Settings["webhook_url"],
			))
		case "discord":
			notifiers = append(notifiers, NewDiscordNotifier(
				n.Settings["webhook_url"],
			))
		case "telegram":
			notifiers = append(notifiers, NewTelegramNotifier(
				n.Settings["bot_token"],
				n.Settings["chat_id"],
			))
		case "gotify":
			notifiers = append(notifiers, NewGotifyNotifier(
				n.Settings["url"],
				n.Settings["token"],
			))
		case "ntfy":
			notifiers = append(notifiers, NewNtfyNotifier(
				n.Settings["url"],
			))
		case "matrix":
			notifiers = append(notifiers, NewMatrixNotifier(
				n.Settings["homeserver"],
				n.Settings["token"],
				n.Settings["room_id"],
			))
		case "sms":
			notifiers = append(notifiers, NewSMSNotifier(
				n.Settings["gateway_url"],
				n.Settings["api_key"],
				n.Settings["from"],
				n.Settings["to"],
			))
		default:
			return nil, fmt.Errorf("notify: unknown notifier type %q", n.Type)
		}
	}

	return notifiers, nil
}
