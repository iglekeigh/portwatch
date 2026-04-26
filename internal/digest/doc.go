// Package digest provides content-addressable fingerprinting for port scan
// results. Each scan result is reduced to a deterministic SHA-256 hex string
// so that downstream components can skip processing when nothing has changed.
//
// The FileStore persists digests across restarts, enabling portwatch to avoid
// re-alerting on an unchanged open-port set even after a process restart.
package digest
