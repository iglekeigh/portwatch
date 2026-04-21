// Package audit provides structured audit logging for port scan events,
// recording who triggered a scan, when, and what changed.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	NewPorts  []int     `json:"new_ports,omitempty"`
	Closed    []int     `json:"closed_ports,omitempty"`
	Open      []int     `json:"open_ports"`
	Changed   bool      `json:"changed"`
}

// Logger writes audit entries to an output destination.
type Logger struct {
	w       io.Writer
	encoder *json.Encoder
}

// NewLogger creates a Logger that writes JSON audit entries to w.
// Pass os.Stdout or an *os.File as needed.
func NewLogger(w io.Writer) *Logger {
	enc := json.NewEncoder(w)
	return &Logger{w: w, enc: enc}
}

// NewFileLogger opens (or creates/appends) a file at path and returns a Logger.
func NewFileLogger(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return NewLogger(f), nil
}

// Record writes an audit entry derived from a scan diff result.
func (l *Logger) Record(host string, diff scanner.DiffResult, open []int) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Host:      host,
		NewPorts:  diff.New,
		Closed:    diff.Closed,
		Open:      open,
		Changed:   diff.HasChanges(),
	}
	if err := l.enc.Encode(entry); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}
