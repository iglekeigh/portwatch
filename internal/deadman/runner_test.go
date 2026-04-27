package deadman_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/deadman"
)

func TestRunner_FiresCallbackOnSilentHost(t *testing.T) {
	var mu sync.Mutex
	var alerts []string

	w := deadman.New(10*time.Millisecond, func(host string, _ time.Duration) {
		mu.Lock()
		alerts = append(alerts, host)
		mu.Unlock()
	})
	w.Checkin("runner-host")

	runner := deadman.NewRunner(w, 15*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	time.Sleep(20 * time.Millisecond) // let checkin age past timeout
	go runner.Run(ctx)
	<-ctx.Done()

	mu.Lock()
	defer mu.Unlock()
	if len(alerts) == 0 {
		t.Fatal("expected at least one dead-man alert from runner")
	}
}

func TestRunner_StopsOnContextCancel(t *testing.T) {
	calls := 0
	w := deadman.New(time.Hour, func(_ string, _ time.Duration) { calls++ })
	w.Checkin("x")

	runner := deadman.NewRunner(w, 5*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		runner.Run(ctx)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("runner did not stop after context cancel")
	}
}
