// Package suppress reduces alert noise by tracking the last-alerted port
// state for each host. An alert is only forwarded to notifiers when the
// open-port fingerprint has changed or the configured TTL has elapsed,
// preventing repeated notifications for a stable, unchanged host.
package suppress
