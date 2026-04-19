package scanner

import (
	"reflect"
	"testing"
)

func TestCompare_NewPorts(t *testing.T) {
	diff := Compare([]int{80}, []int{80, 443})
	if !reflect.DeepEqual(diff.Opened, []int{443}) {
		t.Errorf("expected opened [443], got %v", diff.Opened)
	}
	if len(diff.Closed) != 0 {
		t.Errorf("expected no closed ports, got %v", diff.Closed)
	}
}

func TestCompare_ClosedPorts(t *testing.T) {
	diff := Compare([]int{80, 443}, []int{80})
	if !reflect.DeepEqual(diff.Closed, []int{443}) {
		t.Errorf("expected closed [443], got %v", diff.Closed)
	}
	if len(diff.Opened) != 0 {
		t.Errorf("expected no opened ports, got %v", diff.Opened)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	diff := Compare([]int{80, 443}, []int{80, 443})
	if diff.HasChanges() {
		t.Errorf("expected no changes, got %+v", diff)
	}
}

func TestCompare_EmptyPrevious(t *testing.T) {
	diff := Compare(nil, []int{22, 80})
	if !reflect.DeepEqual(diff.Opened, []int{22, 80}) {
		t.Errorf("expected opened [22,80], got %v", diff.Opened)
	}
}

func TestHasChanges(t *testing.T) {
	if (Diff{}).HasChanges() {
		t.Error("empty diff should have no changes")
	}
	if !(Diff{Opened: []int{8080}}).HasChanges() {
		t.Error("diff with opened ports should have changes")
	}
}
