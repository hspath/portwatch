package ports

import (
	"context"
	"fmt"
	"time"
)

// PipelineOptions configures the scan pipeline.
type PipelineOptions struct {
	Throttle   ThrottleOptions
	Filter     FilterOptions
	Dedupe     DedupeOptions
	RateLimit  RateLimiterOptions
}

// DefaultPipelineOptions returns sensible pipeline defaults.
func DefaultPipelineOptions() PipelineOptions {
	return PipelineOptions{
		Throttle:  DefaultThrottleOptions(),
		Filter:    DefaultFilterOptions(),
		Dedupe:    DefaultDedupeOptions(),
		RateLimit: DefaultRateLimiterOptions(),
	}
}

// ScanResult holds the output of a single pipeline run.
type ScanResult struct {
	Snapshot  *Snapshot
	Added     []Listener
	Removed   []Listener
	Timestamp time.Time
}

// Pipeline combines throttling, filtering, deduplication, and diffing
// into a single reusable scan pipeline.
type Pipeline struct {
	opts      PipelineOptions
	throttle  *Throttle
	dedupe    *Deduplicator
	prev      *Snapshot
}

// NewPipeline constructs a Pipeline with the provided options.
func NewPipeline(opts PipelineOptions) *Pipeline {
	return &Pipeline{
		opts:     opts,
		throttle: NewThrottle(opts.Throttle),
		dedupe:   NewDeduplicator(opts.Dedupe),
	}
}

// Run executes one scan cycle. It returns (nil, nil) if throttled.
func (p *Pipeline) Run(ctx context.Context) (*ScanResult, error) {
	if !p.throttle.Allow() {
		return nil, nil
	}

	listeners, err := ScanListeners()
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	filtered := ApplyFilter(listeners, p.opts.Filter)
	snap := NewSnapshot(filtered)

	result := &ScanResult{
		Snapshot:  snap,
		Timestamp: p.opts.Throttle.Clock(),
	}

	if p.prev != nil {
		added, removed := diff(p.prev.Listeners(), snap.Listeners())
		result.Added = added
		result.Removed = removed
	}

	p.prev = snap
	return result, nil
}

// Reset clears pipeline state, forcing the next Run to be treated as a fresh start.
func (p *Pipeline) Reset() {
	p.throttle.Reset()
	p.dedupe = NewDeduplicator(p.opts.Dedupe)
	p.prev = nil
}
