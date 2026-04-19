package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event holds alert data for a port change.
type Event struct {
	Host      string
	Timestamp time.Time
	Level     Level
	Diff      scanner.Diff
}

// Notifier sends alert events somewhere.
type Notifier interface {
	Notify(e Event) error
}

// ConsoleNotifier writes alerts to an io.Writer (default: os.Stdout).
type ConsoleNotifier struct {
	Out io.Writer
}

// NewConsoleNotifier returns a ConsoleNotifier writing to stdout.
func NewConsoleNotifier() *ConsoleNotifier {
	return &ConsoleNotifier{Out: os.Stdout}
}

// Notify prints a formatted alert to the configured writer.
func (c *ConsoleNotifier) Notify(e Event) error {
	ts := e.Timestamp.Format(time.RFC3339)
	if len(e.Diff.NewPorts) > 0 {
		fmt.Fprintf(c.Out, "[%s] %s [%s] NEW ports opened: %v\n",
			ts, e.Level, e.Host, e.Diff.NewPorts)
	}
	if len(e.Diff.ClosedPorts) > 0 {
		fmt.Fprintf(c.Out, "[%s] %s [%s] ports CLOSED: %v\n",
			ts, e.Level, e.Host, e.Diff.ClosedPorts)
	}
	return nil
}

// BuildEvent constructs an Event from a diff result.
func BuildEvent(host string, d scanner.Diff) Event {
	lvl := LevelInfo
	if len(d.NewPorts) > 0 {
		lvl = LevelAlert
	} else if len(d.ClosedPorts) > 0 {
		lvl = LevelWarn
	}
	return Event{
		Host:      host,
		Timestamp: time.Now(),
		Level:     lvl,
		Diff:      d,
	}
}
