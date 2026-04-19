package report_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/user/portwatch/internal/report"
)

func TestJSONFormatter_ValidJSON(t *testing.T) {
	f := report.NewJSONFormatter(false)
	var buf bytes.Buffer
	if err := f.Format(&buf, fixedSummary()); err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["host"] != "localhost" {
		t.Errorf("expected host localhost, got %v", out["host"])
	}
}

func TestJSONFormatter_Pretty(t *testing.T) {
	f := report.NewJSONFormatter(true)
	var buf bytes.Buffer
	if err := f.Format(&buf, fixedSummary()); err != nil {
		t.Fatal(err)
	}
	if buf.Bytes()[0] != '{' {
		t.Error("expected JSON object")
	}
	// pretty output contains newlines
	if !bytes.Contains(buf.Bytes(), []byte("\n")) {
		t.Error("expected indented output")
	}
}

func TestJSONFormatter_EmptyPorts(t *testing.T) {
	f := report.NewJSONFormatter(false)
	s := report.Summary{Host: "h", OpenPorts: []int{}, NewPorts: []int{}, Closed: []int{}}
	var buf bytes.Buffer
	if err := f.Format(&buf, s); err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
}
