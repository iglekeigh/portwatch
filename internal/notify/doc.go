// Package notify provides additional notification backends for portwatch.
//
// Supported backends:
//   - Email (SMTP)
//   - PagerDuty
//   - OpsGenie
//   - Microsoft Teams
//   - Discord
//
// Each notifier implements the alert.Notifier interface and sends
// a message only when the port diff contains changes.
package notify
