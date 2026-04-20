package history

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry represents a single scan result recorded at a point in time.
type Entry struct {
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
	ScannedAt time.Time `json:"scanned_at"`
}

// Record stores a new entry for the given host.
func Record(store Store, host string, result scanner.Result) error {
	entries, err := store.Load(host)
	if err != nil {
		return err
	}
	entries = append(entries, Entry{
		Host:      host,
		Ports:     result.Open,
		ScannedAt: time.Now().UTC(),
	})
	if len(entries) > MaxEntries {
		entries = entries[len(entries)-MaxEntries:]
	}
	return store.Save(host, entries)
}

// Latest returns the most recent entry for a host, or nil if none.
func Latest(store Store, host string) (*Entry, error) {
	entries, err := store.Load(host)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
	}
	e := entries[len(entries)-1]
	return &e, nil
}

// MaxEntries is the maximum number of history entries retained per host.
const MaxEntries = 100
