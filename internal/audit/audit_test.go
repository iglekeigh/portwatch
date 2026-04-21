package audit_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/scanner"
)

func makeResult(newPorts, closed []int) scanner.DiffResult {
	return scanner.DiffResult{New: newPorts, Closed: closed}
}

func TestRecord_WritesJSONEntry(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	err := l.Record("localhost", makeResult([]int{80, 443}, nil), []int{80, 443, 8080})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if entry.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", entry.Host)
	}
	if len(entry.NewPorts) != 2 {
		t.Errorf("expected 2 new ports, got %d", len(entry.NewPorts))
	}
	if !entry.Changed {
		t.Error("expected Changed=true")
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestRecord_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	err := l.Record("10.0.0.1", makeResult(nil, nil), []int{22})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if entry.Changed {
		t.Error("expected Changed=false when no diff")
	}
}

func TestNewFileLogger_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	l, err := audit.NewFileLogger(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_ = l.Record("host1", makeResult([]int{9090}, nil), []int{9090})

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read log file: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty log file")
	}
}

func TestRecord_TimestampIsUTC(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)
	before := time.Now().UTC()
	_ = l.Record("h", makeResult(nil, nil), nil)
	after := time.Now().UTC()

	var entry audit.Entry
	_ = json.Unmarshal(buf.Bytes(), &entry)

	if entry.Timestamp.Before(before) || entry.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", entry.Timestamp, before, after)
	}
}
