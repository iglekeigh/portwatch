package filter_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/filter"
)

func TestApply_NoRules(t *testing.T) {
	f := filter.New(nil)
	ports := []int{80, 443, 8080}
	got := f.Apply(ports)
	if len(got) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(got))
	}
}

func TestApply_ExcludeRule(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Ports: []int{80}, Exclude: true},
	})
	got := f.Apply([]int{80, 443, 8080})
	for _, p := range got {
		if p == 80 {
			t.Fatal("port 80 should have been excluded")
		}
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
}

func TestApply_IncludeRule(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Ports: []int{443, 8080}, Exclude: false},
	})
	got := f.Apply([]int{80, 443, 8080})
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
}

func TestApply_ExcludeOverridesInclude(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Ports: []int{80, 443}, Exclude: false},
		{Ports: []int{80}, Exclude: true},
	})
	got := f.Apply([]int{80, 443, 8080})
	if len(got) != 1 || got[0] != 443 {
		t.Fatalf("expected [443], got %v", got)
	}
}

func TestMatchesAny(t *testing.T) {
	if !filter.MatchesAny(80, []int{22, 80, 443}) {
		t.Fatal("expected match")
	}
	if filter.MatchesAny(8080, []int{22, 80, 443}) {
		t.Fatal("expected no match")
	}
}
