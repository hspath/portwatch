package ports

import (
	"net"
	"testing"
	"time"
)

func fixedAggClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func aggListener(port uint16) Listener {
	return Listener{IP: net.ParseIP("0.0.0.0"), Port: port, Protocol: "tcp"}
}

func aggOpts(window time.Duration, maxBatch int, now time.Time) AggregatorOptions {
	return AggregatorOptions{
		Window:   window,
		MaxBatch: maxBatch,
		Clock:    fixedAggClock(now),
	}
}

func TestAggregator_NoFlushBeforeWindow(t *testing.T) {
	now := time.Now()
	a := NewAggregator(aggOpts(10*time.Second, 100, now))
	_, flushed := a.Record([]Listener{aggListener(80)}, nil)
	if flushed {
		t.Fatal("expected no flush before window expires")
	}
}

func TestAggregator_FlushOnWindowExpiry(t *testing.T) {
	base := time.Now()
	opts := AggregatorOptions{
		Window:   5 * time.Second,
		MaxBatch: 100,
		Clock:    fixedAggClock(base),
	}
	a := NewAggregator(opts)
	a.Record([]Listener{aggListener(80)}, nil)

	// Advance clock past window
	opts.Clock = fixedAggClock(base.Add(6 * time.Second))
	a.opts = opts

	batch, flushed := a.Record([]Listener{aggListener(443)}, nil)
	if !flushed {
		t.Fatal("expected flush after window expired")
	}
	if len(batch.Added) != 2 {
		t.Fatalf("expected 2 added, got %d", len(batch.Added))
	}
}

func TestAggregator_FlushOnMaxBatch(t *testing.T) {
	now := time.Now()
	a := NewAggregator(aggOpts(60*time.Second, 2, now))

	a.Record([]Listener{aggListener(80)}, nil)
	batch, flushed := a.Record([]Listener{aggListener(443)}, nil)
	if !flushed {
		t.Fatal("expected flush on max batch")
	}
	if len(batch.Added) != 2 {
		t.Fatalf("expected 2 added, got %d", len(batch.Added))
	}
}

func TestAggregator_HasChanges_False(t *testing.T) {
	b := AggregatedBatch{}
	if b.HasChanges() {
		t.Fatal("empty batch should have no changes")
	}
}

func TestAggregator_HasChanges_True(t *testing.T) {
	b := AggregatedBatch{Added: []Listener{aggListener(22)}}
	if !b.HasChanges() {
		t.Fatal("batch with added listeners should have changes")
	}
}

func TestAggregator_ForceFlush(t *testing.T) {
	now := time.Now()
	a := NewAggregator(aggOpts(60*time.Second, 100, now))
	a.Record([]Listener{aggListener(8080)}, nil)

	batch := a.Flush()
	if len(batch.Added) != 1 {
		t.Fatalf("expected 1 added after force flush, got %d", len(batch.Added))
	}

	// Second flush should be empty
	batch2 := a.Flush()
	if batch2.HasChanges() {
		t.Fatal("second flush should be empty")
	}
}

func TestDefaultAggregatorOptions_NonZero(t *testing.T) {
	opts := DefaultAggregatorOptions()
	if opts.Window == 0 {
		t.Error("default window must be non-zero")
	}
	if opts.MaxBatch == 0 {
		t.Error("default max batch must be non-zero")
	}
	if opts.Clock == nil {
		t.Error("default clock must not be nil")
	}
}
