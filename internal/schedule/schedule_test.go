package schedule_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/schedule"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/store"
)

func tempStore(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.New(t.TempDir() + "/state.json")
	if err != nil {
		t.Fatal(err)
	}
	return st
}

func TestRunner_RunCancels(t *testing.T) {
	sc := scanner.New(500 * time.Millisecond)
	st := tempStore(t)
	notifier := alert.NewConsoleNotifier()
	cfg := schedule.Config{
		Hosts:    []string{"127.0.0.1"},
		Ports:    []int{},
		Interval: 50 * time.Millisecond,
	}
	r := schedule.New(cfg, sc, st, notifier)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	err := r.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestRunner_StoresResults(t *testing.T) {
	sc := scanner.New(500 * time.Millisecond)
	st := tempStore(t)
	notifier := alert.NewConsoleNotifier()
	cfg := schedule.Config{
		Hosts:    []string{"127.0.0.1"},
		Ports:    []int{},
		Interval: 1 * time.Hour,
	}
	r := schedule.New(cfg, sc, st, notifier)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	r.Run(ctx) //nolint
	ports, err := st.Get("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if ports == nil {
		t.Error("expected stored ports slice, got nil")
	}
}
