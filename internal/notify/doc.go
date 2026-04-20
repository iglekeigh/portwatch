// Package notify provides additional notification backends for portwatch.
//
// Supported notifiers:
//   - EmailNotifier: sends scan change events via SMTP
//   - PagerDutyNotifier: triggers PagerDuty incidents on port changes
//   - OpsGenieNotifier: creates OpsGenie alerts on port changes
//
// All notifiers implement the alert.Notifier interface and are safe
// to compose using alert.NewMultiNotifier.
package notify
