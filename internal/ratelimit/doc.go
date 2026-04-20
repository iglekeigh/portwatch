// Package ratelimit implements a per-host cooldown mechanism to suppress
// repeated alert notifications within a configurable time window.
//
// This prevents alert storms when a host experiences rapid, repeated
// port changes during a monitoring interval.
package ratelimit
