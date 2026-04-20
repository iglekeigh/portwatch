package report

import "io"

// Formatter writes a Summary to the given writer in a specific format.
type Formatter interface {
	Format(s Summary, w io.Writer) error
}

// Summary holds the data for a single scan report.
type Summary struct {
	Host        string
	ScannedAt   string
	OpenPorts   []int
	NewPorts    []int
	ClosedPorts []int
}

// HasChanges reports whether the summary contains any port changes,
// i.e. ports that were newly opened or recently closed.
func (s Summary) HasChanges() bool {
	return len(s.NewPorts) > 0 || len(s.ClosedPorts) > 0
}
