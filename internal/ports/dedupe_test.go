package ports

import (
	"fmt"
	"testing"
	"time"
)

type mockClock struct {
	now time.Time
}

func (c *mockClock) Now() time.Time { return c.now }

func newDedupe(window time.Duration, clk *mockClock) *Deduplicator {
	return NewDeduplicator(DedupeOptions{
		WindowDuration: window,
		ClockFn:        clk.Now,
	})
}

func TestDeduplicator_FirstCallNotDuplicate(t *testing.T) {
	clk := &mockClock{now: time.Now()}
	d := newDedupe(time.Minute, clk)

	if d.IsDuplicate("tcp:0.0.0.0:8080") {
		t.Fatal("expected first call to not be a duplicate")
	}
}

func TestDeduplicator_SecondCallWithinWindowIsDuplicate(t *testing.T) {
	clk := &mockClock{now: time.Now()}
	d := newDedupe(time.Minute, clk)

	d.IsDuplicate("tcp:0.0.0.0:8080")
	clk.now = clk.now.Add(30 * time.Second)

	if !d.IsDuplicate("tcp:0.0.0.0:8080") {
		t.Fatal("expected second call within window to be a duplicate")
	}
}

func TestDeduplicator_CallAfterWindowNotDuplicate(t *testing.T) {
	clk := &mockClock{now: time.Now()}
	d := newDedupe(time.Minute, clk)

	d.IsDuplicate("tcp:0.0.0.0:9090")
	clk.now = clk.now.Add(2 * time.Minute)

	if d.IsDuplicate("tcp:0.0.0.0:9090") {
		t.Fatal("expected call after window expiry to not be a duplicate")
	}
}

func TestDeduplicator_DifferentKeysAreIndependent(t *testing.T) {
	clk := &mockClock{now: time.Now()}
	d := newDedupe(time.Minute, clk)

	d.IsDuplicate("tcp:0.0.0.0:80")

	if d.IsDuplicate("tcp:0.0.0.0:443") {
		t.Fatal("expected different key to not be a duplicate")
	}
}

func TestDeduplicator_Evict_RemovesOldEntries(t *testing.T) {
	clk := &mockClock{now: time.Now()}
	d := newDedupe(time.Minute, clk)

	for i := 0; i < 5; i++ {
		d.IsDuplicate(fmt.Sprintf("tcp:0.0.0.0:%d", 8000+i))
	}

	if d.Len() != 5 {
		t.Fatalf("expected 5 entries, got %d", d.Len())
	}

	clk.now = clk.now.Add(2 * time.Minute)
	d.Evict()

	if d.Len() != 0 {
		t.Fatalf("expected 0 entries after eviction, got %d", d.Len())
	}
}

func TestDeduplicator_Evict_PreservesRecentEntries(t *testing.T) {
	clk := &mockClock{now: time.Now()}
	d := newDedupe(time.Minute, clk)

	d.IsDuplicate("tcp:0.0.0.0:80")  // old
	clk.now = clk.now.Add(90 * time.Second)
	d.IsDuplicate("tcp:0.0.0.0:443") // recent

	d.Evict()

	if d.Len() != 1 {
		t.Fatalf("expected 1 entry after eviction, got %d", d.Len())
	}
}
