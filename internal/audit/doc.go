// Package audit provides append-only structured audit logging for portwatch.
//
// Each scan cycle that produces a DiffResult can be recorded as a JSON
// audit entry, capturing the host, timestamp, open ports, and any changes
// detected. Entries are written one per line (NDJSON) to any io.Writer,
// making them easy to ingest into log aggregation systems.
package audit
