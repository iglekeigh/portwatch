package filter_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/filter"
)

func TestFromConfig_IncludeOnly(t *testing.T) {
	cfg := filter.Config{IncludePorts: "80,443"}
	f, err := filter.FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := f.Apply([]int{22, 80, 443, 8080})
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %v", got)
	}
}

func TestFromConfig_ExcludeOnly(t *testing.T) {
	cfg := filter.Config{ExcludePorts: "22"}
	f, err := filter.FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := f.Apply([]int{22, 80, 443})
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %v", got)
	}
}

func TestFromConfig_BothRules(t *testing.T) {
	cfg := filter.Config{IncludePorts: "80,443,22", ExcludePorts: "22"}
	f, err := filter.FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := f.Apply([]int{22, 80, 443, 8080})
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %v", got)
	}
}

func TestFromConfig_InvalidRange(t *testing.T) {
	cfg := filter.Config{IncludePorts: "notaport"}
	_, err := filter.FromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid port range")
	}
}

func TestFromConfig_Empty(t *testing.T) {
	cfg := filter.Config{}
	f, err := filter.FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := f.Apply([]int{80, 443})
	if len(got) != 2 {
		t.Fatalf("expected all ports returned, got %v", got)
	}
}
