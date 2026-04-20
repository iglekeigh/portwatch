package config

// NotifierConfig holds the type and settings for a single notification backend.
// The Settings map contains backend-specific key/value pairs (e.g. API keys,
// webhook URLs, recipient addresses).
//
// Example YAML:
//
//	notifiers:
//	  - type: slack
//	    settings:
//	      webhook_url: https://hooks.slack.com/...
//	  - type: sms
//	    settings:
//	      gateway_url: https://sms.example.com/send
//	      api_key: secret
//	      from: "+10000000000"
//	      to: "+19999999999"
type NotifierConfig struct {
	// Type identifies the notification backend (e.g. "slack", "email", "sms").
	Type string `yaml:"type"`

	// Settings contains backend-specific configuration values.
	Settings map[string]string `yaml:"settings"`
}
