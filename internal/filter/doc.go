// Package filter provides port filtering utilities for portwatch.
//
// It supports include and exclude rules that can be composed via Config or
// Rule slices. Rules are evaluated in order, with exclude rules taking
// precedence over include rules when both match a given port.
//
// Basic usage:
//
//	cfg := filter.Config{
//		Include: []uint16{80, 443},
//		Exclude: []uint16{8080},
//	}
//	f := filter.New(cfg)
//	if f.Match(80) {
//		// port 80 is included
//	}
//
// If Include is empty, all ports are considered included by default unless
// explicitly listed in Exclude. If Exclude is empty, no ports are suppressed.
//
// Example with only exclusions (block a single port, allow everything else):
//
//	cfg := filter.Config{
//		Exclude: []uint16{22},
//	}
//	f := filter.New(cfg)
//	if f.Match(80) {
//		// port 80 is allowed (not in exclude list)
//	}
package filter
