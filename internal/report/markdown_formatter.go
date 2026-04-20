package report

import (
	"fmt"
	"strings"
)

// MarkdownFormatter formats a scan summary as a Markdown document.
type MarkdownFormatter struct{}

// NewMarkdownFormatter returns a new MarkdownFormatter.
func NewMarkdownFormatter() *MarkdownFormatter {
	return &MarkdownFormatter{}
}

// Format renders the Summary as a Markdown string.
func (f *MarkdownFormatter) Format(s Summary) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Port Scan Report\n\n"))
	sb.WriteString(fmt.Sprintf("**Host:** %s\n\n", s.Host))
	sb.WriteString(fmt.Sprintf("**Scanned At:** %s\n\n", s.ScannedAt.Format("2006-01-02 15:04:05 UTC")))

	sb.WriteString("## Open Ports\n\n")
	if len(s.OpenPorts) == 0 {
		sb.WriteString("_None_\n\n")
	} else {
		for _, p := range s.OpenPorts {
			sb.WriteString(fmt.Sprintf("- %d\n", p))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## New Ports\n\n")
	if len(s.Diff.New) == 0 {
		sb.WriteString("_None_\n\n")
	} else {
		for _, p := range s.Diff.New {
			sb.WriteString(fmt.Sprintf("- %d\n", p))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Closed Ports\n\n")
	if len(s.Diff.Closed) == 0 {
		sb.WriteString("_None_\n\n")
	} else {
		for _, p := range s.Diff.Closed {
			sb.WriteString(fmt.Sprintf("- %d\n", p))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}
