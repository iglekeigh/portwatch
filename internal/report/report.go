package report

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Summary holds a scan report for a single host.
type Summary struct {
	Host      string
	ScannedAt time.Time
	OpenPorts []int
	NewPorts  []int
	Closed    []int
}

// Formatter writes a Summary to an io.Writer.
type Formatter interface {
	Format(w io.Writer, s Summary) error
}

// TextFormatter renders a human-readable text report.
type TextFormatter struct{}

func NewTextFormatter() *TextFormatter { return &TextFormatter{} }

func (f *TextFormatter) Format(w io.Writer, s Summary) error {
	_, err := fmt.Fprintf(w,
		"[%s] Host: %s\n  Open: %s\n  New: %s\n  Closed: %s\n",
		s.ScannedAt.Format(time.RFC3339),
		s.Host,
		joinInts(s.OpenPorts),
		joinInts(s.NewPorts),
		joinInts(s.Closed),
	)
	return err
}

// Build constructs a Summary from scan results and a diff.
func Build(host string, open []int, diff scanner.Diff) Summary {
	return Summary{
		Host:      host,
		ScannedAt: time.Now().UTC(),
		OpenPorts: open,
		NewPorts:  diff.New,
		Closed:    diff.Closed,
	}
}

// Print writes a Summary to stdout using the given formatter.
func Print(s Summary, f Formatter) error {
	return f.Format(os.Stdout, s)
}

func joinInts(vals []int) string {
	if len(vals) == 0 {
		return "(none)"
	}
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, ", ")
}
