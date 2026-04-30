package ports

import (
	"sync"
	"time"
)

// AggregatorOptions configures the event aggregator.
type AggregatorOptions struct {
	// Window is the duration over which events are batched.
	Window time.Duration
	// MaxBatch is the maximum number of events before an early flush.
	MaxBatch int
	// Clock allows injecting a custom time source for testing.
	Clock func() time.Time
}

// DefaultAggregatorOptions returns sensible defaults.
func DefaultAggregatorOptions() AggregatorOptions {
	return AggregatorOptions{
		Window:   5 * time.Second,
		MaxBatch: 50,
		Clock:    time.Now,
	}
}

// AggregatedBatch holds a set of diff events collected within a window.
type AggregatedBatch struct {
	Added   []Listener
	Removed []Listener
	At      time.Time
}

// HasChanges reports whether the batch contains any events.
func (b AggregatedBatch) HasChanges() bool {
	return len(b.Added) > 0 || len(b.Removed) > 0
}

// Aggregator batches rapid port-change events into fixed windows.
type Aggregator struct {
	opts    AggregatorOptions
	mu      sync.Mutex
	added   []Listener
	removed []Listener
	start   time.Time
}

// NewAggregator creates an Aggregator with the given options.
func NewAggregator(opts AggregatorOptions) *Aggregator {
	return &Aggregator{opts: opts}
}

// Record adds a diff result to the current batch.
// It returns a flushed AggregatedBatch and true when the window
// has elapsed or the batch size limit is reached.
func (a *Aggregator) Record(added, removed []Listener) (AggregatedBatch, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := a.opts.Clock()
	if a.start.IsZero() {
		a.start = now
	}

	a.added = append(a.added, added...)
	a.removed = append(a.removed, removed...)

	total := len(a.added) + len(a.removed)
	windowExpired := now.Sub(a.start) >= a.opts.Window

	if windowExpired || total >= a.opts.MaxBatch {
		return a.flush(now), true
	}
	return AggregatedBatch{}, false
}

// Flush forces emission of whatever has been collected so far.
func (a *Aggregator) Flush() AggregatedBatch {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.flush(a.opts.Clock())
}

func (a *Aggregator) flush(at time.Time) AggregatedBatch {
	batch := AggregatedBatch{
		Added:   a.added,
		Removed: a.removed,
		At:      at,
	}
	a.added = nil
	a.removed = nil
	a.start = time.Time{}
	return batch
}
