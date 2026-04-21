// Package throttle enforces a minimum interval between successive port scans
// of the same host. This prevents portwatch from hammering targets when
// running with a short schedule interval or during manual re-runs.
//
// Usage:
//
//	th := throttle.New(30 * time.Second)
//	if ok, wait := th.Allow(host); !ok {
//		log.Printf("skipping %s, next scan in %v", host, wait)
//		continue
//	}
package throttle
