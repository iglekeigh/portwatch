// Package store provides a lightweight JSON-backed persistence layer for
// portwatch snapshots. Each snapshot records the set of open ports observed
// on a given host at the time of the last scan. Snapshots are keyed by host
// string and written atomically to a single file so that the watcher can
// detect changes between successive runs.
package store
