// Package tag provides host tagging support for grouping and filtering scan targets.
package tag

import "sort"

// Tag represents a label attached to a host.
type Tag = string

// Set holds a deduplicated, sorted collection of tags for a host.
type Set struct {
	tags map[Tag]struct{}
}

// New creates a new Set from the provided tags.
func New(tags ...Tag) *Set {
	s := &Set{tags: make(map[Tag]struct{})}
	for _, t := range tags {
		if t != "" {
			s.tags[t] = struct{}{}
		}
	}
	return s
}

// Add inserts a tag into the set.
func (s *Set) Add(t Tag) {
	if t != "" {
		s.tags[t] = struct{}{}
	}
}

// Has reports whether the set contains the given tag.
func (s *Set) Has(t Tag) bool {
	_, ok := s.tags[t]
	return ok
}

// All returns a sorted slice of all tags in the set.
func (s *Set) All() []Tag {
	out := make([]Tag, 0, len(s.tags))
	for t := range s.tags {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// Len returns the number of tags in the set.
func (s *Set) Len() int {
	return len(s.tags)
}

// MatchesAny reports whether the set contains at least one of the given tags.
func (s *Set) MatchesAny(tags ...Tag) bool {
	for _, t := range tags {
		if s.Has(t) {
			return true
		}
	}
	return false
}

// MatchesAll reports whether the set contains every one of the given tags.
func (s *Set) MatchesAll(tags ...Tag) bool {
	for _, t := range tags {
		if !s.Has(t) {
			return false
		}
	}
	return true
}
