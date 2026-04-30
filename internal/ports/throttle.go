package ports

import (
	"sync"
	"time"
)

// ThrottleOptions controls how the throttle behaves.
type ThrottleOptions struct {
	// MinInterval is the minimum time between successive scans.
	MinInterval time.Duration
	// Clock allows injecting a custom time source for testing.
	Clock func() time.Time
}

// DefaultThrottleOptions returns sensible defaults.
func DefaultThrottleOptions() ThrottleOptions {
	return ThrottleOptions{
		MinInterval: 2 * time.Second,
		Clock:       time.Now,
	}
}

// Throttle prevents scans from running more frequently than MinInterval.
type Throttle struct {
	mu       sync.Mutex
	opts     ThrottleOptions
	lastRun  time.Time
}

// NewThrottle creates a Throttle with the given options.
func NewThrottle(opts ThrottleOptions) *Throttle {
	if opts.Clock == nil {
		opts.Clock = time.Now
	}
	return &Throttle{opts: opts}
}

// Allow returns true if enough time has elapsed since the last allowed call.
// If allowed, it records the current time as the last run.
func (t *Throttle) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.opts.Clock()
	if t.lastRun.IsZero() || now.Sub(t.lastRun) >= t.opts.MinInterval {
		t.lastRun = now
		return true
	}
	return false
}

// Reset clears the last-run timestamp so the next call to Allow will succeed.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastRun = time.Time{}
}

// LastRun returns the time of the most recent allowed call.
func (t *Throttle) LastRun() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastRun
}
