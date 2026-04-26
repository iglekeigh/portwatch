package tag_test

import (
	"testing"

	"github.com/user/portwatch/internal/tag"
)

func TestNew_DeduplicatesTags(t *testing.T) {
	s := tag.New("prod", "prod", "web")
	if s.Len() != 2 {
		t.Fatalf("expected 2 tags, got %d", s.Len())
	}
}

func TestNew_IgnoresEmptyTags(t *testing.T) {
	s := tag.New("", "prod", "")
	if s.Len() != 1 {
		t.Fatalf("expected 1 tag, got %d", s.Len())
	}
}

func TestAdd_InsertsTag(t *testing.T) {
	s := tag.New()
	s.Add("staging")
	if !s.Has("staging") {
		t.Fatal("expected tag 'staging' to be present")
	}
}

func TestHas_MissingTag(t *testing.T) {
	s := tag.New("prod")
	if s.Has("staging") {
		t.Fatal("expected tag 'staging' to be absent")
	}
}

func TestAll_ReturnsSorted(t *testing.T) {
	s := tag.New("web", "prod", "eu")
	all := s.All()
	expected := []string{"eu", "prod", "web"}
	for i, v := range expected {
		if all[i] != v {
			t.Fatalf("expected %s at index %d, got %s", v, i, all[i])
		}
	}
}

func TestMatchesAny_True(t *testing.T) {
	s := tag.New("prod", "web")
	if !s.MatchesAny("db", "web") {
		t.Fatal("expected MatchesAny to return true")
	}
}

func TestMatchesAny_False(t *testing.T) {
	s := tag.New("prod")
	if s.MatchesAny("staging", "dev") {
		t.Fatal("expected MatchesAny to return false")
	}
}

func TestMatchesAll_True(t *testing.T) {
	s := tag.New("prod", "web", "eu")
	if !s.MatchesAll("prod", "eu") {
		t.Fatal("expected MatchesAll to return true")
	}
}

func TestMatchesAll_False(t *testing.T) {
	s := tag.New("prod", "web")
	if s.MatchesAll("prod", "db") {
		t.Fatal("expected MatchesAll to return false")
	}
}

func TestMatchesAll_EmptyInput(t *testing.T) {
	s := tag.New("prod")
	if !s.MatchesAll() {
		t.Fatal("expected MatchesAll with no args to return true")
	}
}
