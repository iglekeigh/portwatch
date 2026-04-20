// Package retry provides a simple exponential back-off retry helper used by
// portwatch notifiers when delivering alerts over unreliable transports.
//
// Usage:
//
//	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
//		return notifier.Notify(event)
//	})
package retry
