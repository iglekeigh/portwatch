package report

import (
	"encoding/json"
	"io"
)

// JSONFormatter renders a Summary as a JSON object.
type JSONFormatter struct {
	Pretty bool
}

func NewJSONFormatter(pretty bool) *JSONFormatter {
	return &JSONFormatter{Pretty: pretty}
}

func (f *JSONFormatter) Format(w io.Writer, s Summary) error {
	enc := json.NewEncoder(w)
	if f.Pretty {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(struct {
		Host      string `json:"host"`
		ScannedAt string `json:"scanned_at"`
		OpenPorts []int  `json:"open_ports"`
		NewPorts  []int  `json:"new_ports"`
		Closed    []int  `json:"closed_ports"`
	}{
		Host:      s.Host,
		ScannedAt: s.ScannedAt.Format("2006-01-02T15:04:05Z"),
		OpenPorts: s.OpenPorts,
		NewPorts:  s.NewPorts,
		Closed:    s.Closed,
	})
}
