// Package notify provides additional notification backends for portwatch.
//
// Each notifier implements the alert.Notifier interface and sends port-change
// events to an external service. Available notifiers:
//
//   - EmailNotifier   – SMTP email alerts
//   - PagerDutyNotifier – PagerDuty incidents
//   - OpsGenieNotifier  – OpsGenie alerts
//   - TeamsNotifier    – Microsoft Teams webhook
//   - DiscordNotifier  – Discord webhook
//   - TelegramNotifier – Telegram Bot API
//   - GotifyNotifier   – Gotify push notifications
//   - NtfyNotifier     – ntfy.sh push notifications
//   - MatrixNotifier   – Matrix room messages via Client-Server API
package notify
