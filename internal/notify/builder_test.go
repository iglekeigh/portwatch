package notify

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func notifierConfig(typ string, settings map[string]string) config.NotifierConfig {
	return config.NotifierConfig{Type: typ, Settings: settings}
}

func TestFromConfig_SMS(t *testing.T) {
	cfg := config.Config{
		Notifiers: []config.NotifierConfig{
			notifierConfig("sms", map[string]string{
				"gateway_url": "https://sms.example.com",
				"api_key":     "secret",
				"from":        "+10000000000",
				"to":          "+19999999999",
			}),
		},
	}
	notifiers, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) != 1 {
		t.Fatalf("expected 1 notifier, got %d", len(notifiers))
	}
	if _, ok := notifiers[0].(*SMSNotifier); !ok {
		t.Errorf("expected *SMSNotifier, got %T", notifiers[0])
	}
}

func TestFromConfig_MultipleNotifiers(t *testing.T) {
	cfg := config.Config{
		Notifiers: []config.NotifierConfig{
			notifierConfig("discord", map[string]string{"webhook_url": "https://discord.example.com"}),
			notifierConfig("ntfy", map[string]string{"url": "https://ntfy.sh/topic"}),
		},
	}
	notifiers, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) != 2 {
		t.Fatalf("expected 2 notifiers, got %d", len(notifiers))
	}
}

func TestFromConfig_UnknownType(t *testing.T) {
	cfg := config.Config{
		Notifiers: []config.NotifierConfig{
			notifierConfig("carrier_pigeon", map[string]string{}),
		},
	}
	_, err := FromConfig(cfg)
	if err == nil {
		t.Error("expected error for unknown notifier type")
	}
}

func TestFromConfig_Empty(t *testing.T) {
	cfg := config.Config{}
	notifiers, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifiers) != 0 {
		t.Errorf("expected 0 notifiers, got %d", len(notifiers))
	}
}
