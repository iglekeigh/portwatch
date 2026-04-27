package deadman

import (
	"context"
	"time"
)

// Runner periodically calls Watcher.Check on a fixed interval.
type Runner struct {
	watcher  *Watcher
	interval time.Duration
}

// NewRunner creates a Runner that checks the watcher every interval.
func NewRunner(w *Watcher, interval time.Duration) *Runner {
	return &Runner{watcher: w, interval: interval}
}

// Run starts the periodic check loop. It blocks until ctx is cancelled.
func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.watcher.Check()
		case <-ctx.Done():
			return
		}
	}
}
