package notify

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// EmailConfig holds SMTP configuration for email notifications.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

type emailNotifier struct {
	cfg EmailConfig
}

// NewEmailNotifier creates a Notifier that sends alerts via SMTP email.
func NewEmailNotifier(cfg EmailConfig) alert.Notifier {
	return &emailNotifier{cfg: cfg}
}

func (e *emailNotifier) Notify(event alert.Event) error {
	if !event.Diff.HasChanges() {
		return nil
	}

	subject := fmt.Sprintf("[portwatch] Port changes detected on %s", event.Host)
	body := buildEmailBody(event)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		e.cfg.From,
		strings.Join(e.cfg.To, ", "),
		subject,
		body,
	))

	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	var auth smtp.Auth
	if e.cfg.Username != "" {
		auth = smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.Host)
	}

	return smtp.SendMail(addr, auth, e.cfg.From, e.cfg.To, msg)
}

func buildEmailBody(event alert.Event) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Host: %s\nScanned at: %s\n\n", event.Host, event.ScannedAt.Format("2006-01-02 15:04:05")))
	if len(event.Diff.New) > 0 {
		sb.WriteString(fmt.Sprintf("New ports: %s\n", joinPorts(event.Diff.New)))
	}
	if len(event.Diff.Closed) > 0 {
		sb.WriteString(fmt.Sprintf("Closed ports: %s\n", joinPorts(event.Diff.Closed)))
	}
	return sb.String()
}

func joinPorts(ports []int) string {
	parts := make([]string, len(ports))
	for i, p := range ports {
		parts[i] = fmt.Sprintf("%d", p)
	}
	return strings.Join(parts, ", ")
}

// Ensure emailNotifier satisfies Notifier at compile time.
var _ alert.Notifier = (*emailNotifier)(nil)

// Compile-time reference to avoid unused import.
var _ = scanner.Result{}
