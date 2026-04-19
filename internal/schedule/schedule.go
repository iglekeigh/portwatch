package schedule

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/store"
)

// Runner periodically scans hosts and emits alerts on port changes.
type Runner struct {
	cfg      Config
	scanner  *scanner.Scanner
	store    *store.Store
	notifier alert.Notifier
}

// Config holds runner configuration.
type Config struct {
	Hosts    []string
	Ports    []int
	Interval time.Duration
}

// New creates a new Runner.
func New(cfg Config, sc *scanner.Scanner, st *store.Store, n alert.Notifier) *Runner {
	return &Runner{cfg: cfg, scanner: sc, store: st, notifier: n}
}

// Run starts the scan loop, blocking until ctx is cancelled.
func (r *Runner) Run(ctx context.Context) error {
	if err := r.tick(); err != nil {
		log.Printf("scan error: %v", err)
	}
	ticker := time.NewTicker(r.cfg.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := r.tick(); err != nil {
				log.Printf("scan error: %v", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *Runner) tick() error {
	for _, host := range r.cfg.Hosts {
		open, err := r.scanner.Scan(host, r.cfg.Ports)
		if err != nil {
			return err
		}
		prev, err := r.store.Get(host)
		if err != nil {
			return err
		}
		diff := scanner.Compare(prev, open)
		ev := alert.BuildEvent(host, diff)
		if err := r.notifier.Notify(ev); err != nil {
			log.Printf("notify error: %v", err)
		}
		if err := r.store.Save(host, open); err != nil {
			return err
		}
	}
	return nil
}
