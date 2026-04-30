package ports

import (
	"sync"
	"time"
)

// RateLimiterOptions configures the rate limiter.
type RateLimiterOptions struct {
	// Window is the duration over which calls are counted.
	Window time.Duration
	// MaxCalls is the maximum number of allowed calls per window.
	MaxCalls int
	// Clock is an optional injectable time source.
	Clock func() time.Time
}

// DefaultRateLimiterOptions returns sensible defaults.
func DefaultRateLimiterOptions() RateLimiterOptions {
	return RateLimiterOptions{
		Window:   time.Minute,
		MaxCalls: 60,
		Clock:    time.Now,
	}
}

type entry struct {
	calls []time.Time
}

// RateLimiter tracks call counts per key within a sliding window.
type RateLimiter struct {
	mu      sync.Mutex
	opts    RateLimiterOptions
	entries map[string]*entry
}

// NewRateLimiter creates a RateLimiter with the given options.
func NewRateLimiter(opts RateLimiterOptions) *RateLimiter {
	if opts.Clock == nil {
		opts.Clock = time.Now
	}
	if opts.MaxCalls <= 0 {
		opts.MaxCalls = 60
	}
	if opts.Window <= 0 {
		opts.Window = time.Minute
	}
	return &RateLimiter{
		opts:    opts,
		entries: make(map[string]*entry),
	}
}

// Allow returns true if the key has not exceeded MaxCalls within the Window.
func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := r.opts.Clock()
	cutoff := now.Add(-r.opts.Window)
	e, ok := r.entries[key]
	if !ok {
		e = &entry{}
		r.entries[key] = e
	}
	// Evict old calls outside the window
	valid := e.calls[:0]
	for _, ts := range e.calls {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}
	e.calls = valid
	if len(e.calls) >= r.opts.MaxCalls {
		return false
	}
	e.calls = append(e.calls, now)
	return true
}
