package ports

import (
	"testing"
	"time"
)

type mockClock struct {
	now time.Time
}

func (c *mockClock) Now() time.Time { return c.now }
func (c *mockClock) Advance(d time.Duration) { c.now = c.now.Add(d) }

func newTestThrottle(interval time.Duration, clk *mockClock) *Throttle {
	return NewThrottle(ThrottleOptions{
		MinInterval: interval,
		Clock:       clk.Now,
	})
}

func TestThrottle_FirstCallAllowed(t *testing.T) {
	clk := &mockClock{now: time.Now()}
	th := newTestThrottle(5*time.Second, clk)
	if !th.Allow() {
		t.Fatal("expected first call to be allowed")
	}
}

func TestThrottle_SecondCallWithinIntervalBlocked(t *testing.T) {
	clk := &mockClock{now: time.Now()}
	th := newTestThrottle(5*time.Second, clk)
	th.Allow()
	clk.Advance(1 * time.Second)
	if th.Allow() {
		t.Fatal("expected second call within interval to be blocked")
	}
}

func TestThrottle_AllowedAfterIntervalExpires(t *testing.T) {
	clk := &mockClock{now: time.Now()}
	th := newTestThrottle(5*time.Second, clk)
	th.Allow()
	clk.Advance(6 * time.Second)
	if !th.Allow() {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestThrottle_Reset_AllowsImmediately(t *testing.T) {
	clk := &mockClock{now: time.Now()}
	th := newTestThrottle(5*time.Second, clk)
	th.Allow()
	clk.Advance(1 * time.Second)
	th.Reset()
	if !th.Allow() {
		t.Fatal("expected Allow after Reset to succeed")
	}
}

func TestThrottle_LastRun_UpdatedOnAllow(t *testing.T) {
	base := time.Now()
	clk := &mockClock{now: base}
	th := newTestThrottle(5*time.Second, clk)
	th.Allow()
	if !th.LastRun().Equal(base) {
		t.Errorf("expected LastRun %v, got %v", base, th.LastRun())
	}
}

func TestThrottle_DefaultOptions(t *testing.T) {
	opts := DefaultThrottleOptions()
	if opts.MinInterval != 2*time.Second {
		t.Errorf("expected 2s default interval, got %v", opts.MinInterval)
	}
	if opts.Clock == nil {
		t.Error("expected non-nil default clock")
	}
}
