package filter

// Rule defines a port filter rule.
type Rule struct {
	Ports    []int
	Exclude  bool
}

// Filter applies inclusion/exclusion rules to a list of ports.
type Filter struct {
	rules []Rule
}

// New creates a Filter with the given rules.
func New(rules []Rule) *Filter {
	return &Filter{rules: rules}
}

// Apply returns ports that pass all filter rules.
func (f *Filter) Apply(ports []int) []int {
	if len(f.rules) == 0 {
		return ports
	}

	excludeSet := make(map[int]bool)
	includeSet := make(map[int]bool)
	hasInclude := false

	for _, rule := range f.rules {
		for _, p := range rule.Ports {
			if rule.Exclude {
				excludeSet[p] = true
			} else {
				includeSet[p] = true
				hasInclude = true
			}
		}
	}

	var result []int
	for _, p := range ports {
		if excludeSet[p] {
			continue
		}
		if hasInclude && !includeSet[p] {
			continue
		}
		result = append(result, p)
	}
	return result
}

// MatchesAny returns true if port is in the given list.
func MatchesAny(port int, ports []int) bool {
	for _, p := range ports {
		if p == port {
			return true
		}
	}
	return false
}
