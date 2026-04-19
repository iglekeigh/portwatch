package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func writeTempConfig(t *testing.T, v any) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(v); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func TestLoad_BasicConfig(t *testing.T) {
	raw := map[string]any{
		"interval_seconds": 30,
		"store_path":       "/tmp/pw.db",
		"hosts": []map[string]any{
			{"address": "localhost", "port_range": "80-443", "label": "local"},
		},
	}
	path := writeTempConfig(t, raw)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.IntervalSeconds != 30 {
		t.Errorf("expected interval 30, got %d", cfg.IntervalSeconds)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s duration, got %v", cfg.Interval)
	}
	if cfg.StorePath != "/tmp/pw.db" {
		t.Errorf("unexpected store path: %s", cfg.StorePath)
	}
	if len(cfg.Hosts) != 1 || cfg.Hosts[0].Address != "localhost" {
		t.Errorf("unexpected hosts: %+v", cfg.Hosts)
	}
}

func TestLoad_Defaults(t *testing.T) {
	path := writeTempConfig(t, map[string]any{})
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.IntervalSeconds != 60 {
		t.Errorf("expected default interval 60, got %d", cfg.IntervalSeconds)
	}
	if cfg.StorePath != "portwatch.db" {
		t.Errorf("expected default store path, got %s", cfg.StorePath)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.Interval != 60*time.Second {
		t.Errorf("unexpected default interval: %v", cfg.Interval)
	}
	if cfg.Hosts == nil {
		t.Error("hosts slice should not be nil")
	}
}
