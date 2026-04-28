// Package cooldown enforces a minimum interval between successive port scans
// of the same host. It prevents redundant or overly frequent scanning when
// multiple schedule ticks fire in quick succession or when a host is listed
// in several target groups.
//
// Use New to create a Cooldown with a configured interval, then call Allow
// before initiating each scan. Reset can be used to force an immediate
// re-scan, for example after a manual trigger or configuration reload.
package cooldown
