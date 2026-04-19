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
