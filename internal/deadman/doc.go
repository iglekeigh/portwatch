// Package deadman implements a dead-man's switch for portwatch hosts.
//
// A Watcher tracks the last successful scan time for each monitored host.
// If a host has not checked in within the configured silence timeout,
// the registered callback is invoked with the host name and how long it
// has been silent. This allows operators to detect when a scheduled scan
// has stopped running entirely — not just when ports change.
//
// Usage:
//
//	w := deadman.New(5*time.Minute, func(host string, d time.Duration) {
//	    log.Printf("DEAD: %s silent for %v", host, d)
//	})
//	w.Checkin("192.168.1.1")   // call after each successful scan
//
//	// Periodically check (or use Runner):
//	runner := deadman.NewRunner(w, time.Minute)
//	runner.Run(ctx)
package deadman
