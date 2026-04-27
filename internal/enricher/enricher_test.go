package enricher

import (
	"errors"
	"testing"
)

func TestEnrich_IPWithReverse(t *testing.T) {
	e := newWithLookup(func(host string) ([]string, error) {
		if host == "93.184.216.34" {
			return []string{"example.com."}, nil
		}
		return nil, errors.New("not found")
	})

	m := e.Enrich("93.184.216.34")

	if !m.IsIP {
		t.Error("expected IsIP to be true")
	}
	if m.Reverse != "example.com" {
		t.Errorf("expected reverse 'example.com', got %q", m.Reverse)
	}
	if m.Host != "93.184.216.34" {
		t.Errorf("unexpected host %q", m.Host)
	}
}

func TestEnrich_IPNoReverse(t *testing.T) {
	e := newWithLookup(func(host string) ([]string, error) {
		return nil, errors.New("no PTR")
	})

	m := e.Enrich("10.0.0.1")

	if !m.IsIP {
		t.Error("expected IsIP to be true")
	}
	if m.Reverse != "" {
		t.Errorf("expected empty reverse, got %q", m.Reverse)
	}
}

func TestEnrich_Hostname(t *testing.T) {
	e := New()

	m := e.Enrich("localhost")

	if m.IsIP {
		t.Error("expected IsIP to be false for hostname")
	}
	if m.Host != "localhost" {
		t.Errorf("unexpected host %q", m.Host)
	}
	if m.Reverse != "localhost" {
		t.Errorf("expected reverse to equal host, got %q", m.Reverse)
	}
}

func TestMeta_String_WithReverse(t *testing.T) {
	m := Meta{Host: "1.2.3.4", Reverse: "host.example.com", IsIP: true}
	got := m.String()
	want := "1.2.3.4 (host.example.com)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestMeta_String_NoReverse(t *testing.T) {
	m := Meta{Host: "1.2.3.4", IsIP: true}
	if m.String() != "1.2.3.4" {
		t.Errorf("unexpected string %q", m.String())
	}
}
