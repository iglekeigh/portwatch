package scanner

// Diff represents changes between two port scans.
type Diff struct {
	Opened []int
	Closed []int
}

// HasChanges returns true if there are any port changes.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// Compare computes the difference between a previous and current set of open ports.
func Compare(previous, current []int) Diff {
	prevSet := toSet(previous)
	currSet := toSet(current)

	var opened, closed []int

	for p := range currSet {
		if !prevSet[p] {
			opened = append(opened, p)
		}
	}
	for p := range prevSet {
		if !currSet[p] {
			closed = append(closed, p)
		}
	}

	sortInts(opened)
	sortInts(closed)

	return Diff{Opened: opened, Closed: closed}
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}

func sortInts(s []int) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
