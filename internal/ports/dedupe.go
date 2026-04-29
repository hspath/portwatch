package ports

import (
	"sync"
	"time"
)

// DedupeOptions configures the deduplication window.
type DedupeOptions struct {
	// WindowDuration is how long a seen event is suppressed.
	WindowDuration time.Duration
	// ClockFn allows injecting a clock for testing.
	ClockFn func() time.Time
}

// DefaultDedupeOptions returns sensible defaults.
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{
		WindowDuration: 5 * time.Minute,
		ClockFn:        time.Now,
	}
}

// Deduplicator suppresses repeated change events within a time window.
// It is safe for concurrent use.
type Deduplicator struct {
	opts DedupeOptions
	mu   sync.Mutex
	seen map[string]time.Time
}

// NewDeduplicator creates a Deduplicator with the given options.
func NewDeduplicator(opts DedupeOptions) *Deduplicator {
	return &Deduplicator{
		opts: opts,
		seen: make(map[string]time.Time),
	}
}

// IsDuplicate returns true if the key was seen within the configured window.
// If it is not a duplicate, the key is recorded with the current time.
func (d *Deduplicator) IsDuplicate(key string) bool {
	now := d.opts.ClockFn()
	d.mu.Lock()
	defer d.mu.Unlock()

	if last, ok := d.seen[key]; ok {
		if now.Sub(last) < d.opts.WindowDuration {
			return true
		}
	}
	d.seen[key] = now
	return false
}

// Evict removes all entries older than the configured window.
// Call periodically to prevent unbounded memory growth.
func (d *Deduplicator) Evict() {
	now := d.opts.ClockFn()
	d.mu.Lock()
	defer d.mu.Unlock()

	for k, last := range d.seen {
		if now.Sub(last) >= d.opts.WindowDuration {
			delete(d.seen, k)
		}
	}
}

// Len returns the number of tracked keys.
func (d *Deduplicator) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}
