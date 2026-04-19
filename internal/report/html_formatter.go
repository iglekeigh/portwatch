package report

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
)

const htmlTmpl = `<!DOCTYPE html>
<html>
<head><title>PortWatch Report</title>
<style>body{font-family:sans-serif;margin:2em}table{border-collapse:collapse}td,th{border:1px solid #ccc;padding:6px 12px}th{background:#f0f0f0}.new{color:green}.closed{color:red}</style>
</head>
<body>
<h1>PortWatch Report</h1>
<p><strong>Host:</strong> {{.Host}}</p>
<p><strong>Scanned At:</strong> {{.ScannedAt}}</p>
<h2>Open Ports</h2>
<table><tr><th>Port</th></tr>
{{range .OpenPorts}}<tr><td>{{.}}</td></tr>{{end}}
</table>
<h2>Changes</h2>
<table><tr><th>Type</th><th>Ports</th></tr>
<tr class="new"><td>New</td><td>{{.NewPorts}}</td></tr>
<tr class="closed"><td>Closed</td><td>{{.ClosedPorts}}</td></tr>
</table>
</body></html>
`

type htmlFormatter struct {
	pretty bool
}

// NewHTMLFormatter returns a Formatter that renders an HTML report.
func NewHTMLFormatter() Formatter {
	return &htmlFormatter{}
}

func (h *htmlFormatter) Format(s Summary, w io.Writer) error {
	tmpl, err := template.New("report").Parse(htmlTmpl)
	if err != nil {
		return fmt.Errorf("html formatter: parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, s); err != nil {
		return fmt.Errorf("html formatter: execute template: %w", err)
	}
	_, err = w.Write(buf.Bytes())
	return err
}
