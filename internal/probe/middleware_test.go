package probe_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/probe"
)

func TestMiddleware_PassesThroughWhenReachable(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	called := false
	next := func(ctx context.Context, host string) ([]int, error) {
		called = true
		return []int{80, 443}, nil
	}

	mw := probe.NewMiddleware(2*time.Second, port, next)
	ports, err := mw.Scan(context.Background(), "127.0.0.1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected next to be called")
	}
	if len(ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(ports))
	}
}

func TestMiddleware_SkipsScanWhenUnreachable(t *testing.T) {
	called := false
	next := func(ctx context.Context, host string) ([]int, error) {
		called = true
		return nil, nil
	}

	mw := probe.NewMiddleware(200*time.Millisecond, "9", next)
	_, err := mw.Scan(context.Background(), "127.0.0.1")

	if err == nil {
		t.Fatal("expected error for unreachable host")
	}
	if called {
		t.Error("expected next NOT to be called")
	}
}

func TestMiddleware_PropagatesNextError(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	scanErr := errors.New("scan failed")
	next := func(ctx context.Context, host string) ([]int, error) {
		return nil, scanErr
	}

	mw := probe.NewMiddleware(2*time.Second, port, next)
	_, err := mw.Scan(context.Background(), "127.0.0.1")

	if !errors.Is(err, scanErr) {
		t.Errorf("expected scanErr, got %v", err)
	}
}
