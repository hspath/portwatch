package ports

import (
	"context"
	"log"
	"time"
)

// WatchOptions configures the polling behaviour of Watch.
type WatchOptions struct {
	Interval      time.Duration
	FilterOptions FilterOptions
}

// DefaultWatchOptions returns sensible defaults for Watch.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{
		Interval:      15 * time.Second,
		FilterOptions: DefaultFilterOptions(),
	}
}

// WatchEvent is emitted whenever the set of active listeners changes.
type WatchEvent struct {
	Added   []Listener
	Removed []Listener
	Current *Snapshot
}

// Watch polls for open ports at the configured interval and sends a
// WatchEvent on the returned channel whenever the listener set changes.
// The channel is closed when ctx is cancelled.
func Watch(ctx context.Context, opts WatchOptions) (<-chan WatchEvent, error) {
	initial, err := NewSnapshot(opts.FilterOptions)
	if err != nil {
		return nil, err
	}

	ch := make(chan WatchEvent, 4)

	go func() {
		defer close(ch)
		prev := initial
		ticker := time.NewTicker(opts.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				curr, err := NewSnapshot(opts.FilterOptions)
				if err != nil {
					log.Printf("portwatch/watcher: scan error: %v", err)
					continue
				}

				added, removed := diff(prev, curr)
				if len(added) > 0 || len(removed) > 0 {
					ch <- WatchEvent{
						Added:   added,
						Removed: removed,
						Current: curr,
					}
				}
				prev = curr
			}
		}
	}()

	return ch, nil
}

// diff returns listeners added and removed between two snapshots.
func diff(prev, curr *Snapshot) (added, removed []Listener) {
	for _, l := range curr.Listeners {
		if !prev.Contains(l) {
			added = append(added, l)
		}
	}
	for _, l := range prev.Listeners {
		if !curr.Contains(l) {
			removed = append(removed, l)
		}
	}
	return
}
