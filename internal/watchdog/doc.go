// Package watchdog provides a self-monitoring component for portwatch.
//
// A Watchdog tracks the last time each configured host produced a scan
// result. If a host goes silent for longer than the configured deadline,
// the registered AlertFunc is invoked so operators can investigate
// whether the scanner itself has stalled or lost connectivity.
//
// Typical usage:
//
//	w := watchdog.New(5*time.Minute, func(host string, d time.Duration) {
//	    log.Printf("no scan result for %s in %v", host, d)
//	})
//	w.Checkin("192.168.1.1") // call after each successful scan
//	w.Run(ctx, 30*time.Second) // background checker
package watchdog
