// Package history records and retrieves per-host port scan history.
// Each scan result is appended as a timestamped Entry and persisted via
// a Store implementation. A capped ring of MaxEntries is maintained so
// disk usage stays bounded over long-running deployments.
package history
