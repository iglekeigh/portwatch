package report

import (
	"bytes"
	"fmt"
	"strings"
)

// CSVFormatter formats a scan summary as CSV output.
type CSVFormatter struct{}

// NewCSVFormatter returns a new CSVFormatter.
func NewCSVFormatter() *CSVFormatter {
	return &CSVFormatter{}
}

// Format returns a CSV representation of the summary.
// Columns: host, status, ports
func (f *CSVFormatter) Format(s Summary) (string, error) {
	var buf bytes.Buffer

	buf.WriteString("host,status,ports\n")

	if len(s.OpenPorts) == 0 && len(s.NewPorts) == 0 && len(s.ClosedPorts) == 0 {
		buf.WriteString(fmt.Sprintf("%s,no-change,\n", s.Host))
		return buf.String(), nil
	}

	if len(s.NewPorts) > 0 {
		ports := intsToStrings(s.NewPorts)
		buf.WriteString(fmt.Sprintf("%s,new,%s\n", s.Host, strings.Join(ports, ";")))
	}

	if len(s.ClosedPorts) > 0 {
		ports := intsToStrings(s.ClosedPorts)
		buf.WriteString(fmt.Sprintf("%s,closed,%s\n", s.Host, strings.Join(ports, ";")))
	}

	if len(s.OpenPorts) > 0 {
		ports := intsToStrings(s.OpenPorts)
		buf.WriteString(fmt.Sprintf("%s,open,%s\n", s.Host, strings.Join(ports, ";")))
	}

	return buf.String(), nil
}

func intsToStrings(ports []int) []string {
	out := make([]string, len(ports))
	for i, p := range ports {
		out[i] = fmt.Sprintf("%d", p)
	}
	return out
}
