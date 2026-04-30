package batch_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/batch"
)

// mockScanner is a test double for batch.Scanner.
type mockScanner struct {
	ports []int
	err   error
	calls atomic.Int32
}

func (m *mockScanner) Scan(_ context.Context, _ string) ([]int, error) {
	m.calls.Add(1)
	return m.ports, m.err
}

func TestRun_ReturnsResultsForAllHosts(t *testing.T) {
	ms := &mockScanner{ports: []int{80, 443}}
	r := batch.New(ms, 2)

	hosts := []string{"host1", "host2", "host3"}
	results := r.Run(context.Background(), hosts)

	if len(results) != len(hosts) {
		t.Fatalf("expected %d results, got %d", len(hosts), len(results))
	}
	for i, res := range results {
		if res.Host != hosts[i] {
			t.Errorf("result[%d].Host = %q, want %q", i, res.Host, hosts[i])
		}
		if res.Err != nil {
			t.Errorf("result[%d].Err = %v, want nil", i, res.Err)
		}
		if len(res.Ports) != 2 {
			t.Errorf("result[%d].Ports len = %d, want 2", i, len(res.Ports))
		}
	}
}

func TestRun_PropagatesErrors(t *testing.T) {
	sentinel := errors.New("scan failed")
	ms := &mockScanner{err: sentinel}
	r := batch.New(ms, 2)

	results := r.Run(context.Background(), []string{"host1"})

	if !errors.Is(results[0].Err, sentinel) {
		t.Errorf("expected sentinel error, got %v", results[0].Err)
	}
}

func TestRun_RespectsContextCancellation(t *testing.T) {
	blocking := &blockingScanner{}
	r := batch.New(blocking, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	results := r.Run(ctx, []string{"host1"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestHasErrors_True(t *testing.T) {
	results := []batch.Result{{Host: "h1", Err: errors.New("oops")}}
	if !batch.HasErrors(results) {
		t.Error("expected HasErrors to return true")
	}
}

func TestHasErrors_False(t *testing.T) {
	results := []batch.Result{{Host: "h1", Ports: []int{22}}}
	if batch.HasErrors(results) {
		t.Error("expected HasErrors to return false")
	}
}

func TestSuccessful_FiltersErrors(t *testing.T) {
	results := []batch.Result{
		{Host: "h1", Ports: []int{80}},
		{Host: "h2", Err: errors.New("fail")},
		{Host: "h3", Ports: []int{443}},
	}
	ok := batch.Successful(results)
	if len(ok) != 2 {
		t.Fatalf("expected 2 successful results, got %d", len(ok))
	}
}

func TestNew_DefaultConcurrency(t *testing.T) {
	ms := &mockScanner{ports: []int{22}}
	r := batch.New(ms, 0) // 0 should default to 4
	results := r.Run(context.Background(), []string{"h1", "h2"})
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

// blockingScanner blocks until context is done.
type blockingScanner struct{}

func (b *blockingScanner) Scan(ctx context.Context, _ string) ([]int, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}
