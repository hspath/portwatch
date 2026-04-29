package ports

import (
	"sync"
	"time"
)

// RateLimiter prevents alert flooding by suppressing repeated events
// for the same port within a configurable window.
type RateLimiter struct {
	mu      sync.Mutex
	window  time.Duration
	seen    map[string]time.Time
	nowFunc func() time.Time
}

// NewRateLimiter creates a RateLimiter with the given suppression window.
func NewRateLimiter(window time.Duration) *RateLimiter {
	return &RateLimiter{
		window:  window,
		seen:    make(map[string]time.Time),
		nowFunc: time.Now,
	}
}

// Allow returns true if the event for the given key should be allowed
// through (i.e. not suppressed). It records the event time on first
// occurrence or after the window has elapsed.
func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.nowFunc()
	if last, ok := r.seen[key]; ok {
		if now.Sub(last) < r.window {
			return false
		}
	}
	r.seen[key] = now
	return true
}

// Purge removes all entries older than the window, freeing memory for
// ports that are no longer active.
func (r *RateLimiter) Purge() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.nowFunc()
	for k, t := range r.seen {
		if now.Sub(t) >= r.window {
			delete(r.seen, k)
		}
	}
}

// Len returns the number of tracked keys.
func (r *RateLimiter) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.seen)
}
