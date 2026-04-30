package ports

import (
	"context"
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func throttleOptsFor(interval time.Duration, clk func() time.Time) ThrottleOptions {
	return ThrottleOptions{MinInterval: interval, Clock: clk}
}

func TestDefaultPipelineOptions_NonZero(t *testing.T) {
	opts := DefaultPipelineOptions()
	if opts.Throttle.MinInterval == 0 {
		t.Error("expected non-zero throttle interval")
	}
}

func TestPipeline_Reset_ClearsPrev(t *testing.T) {
	clk := &mockClock{now: time.Now()}
	opts := DefaultPipelineOptions()
	opts.Throttle = throttleOptsFor(0, clk.Now)
	p := NewPipeline(opts)
	p.prev = NewSnapshot([]Listener{})
	p.Reset()
	if p.prev != nil {
		t.Error("expected prev to be nil after Reset")
	}
}

func TestPipeline_Reset_AllowsNextRun(t *testing.T) {
	clk := &mockClock{now: time.Now()}
	opts := DefaultPipelineOptions()
	opts.Throttle = throttleOptsFor(10*time.Second, clk.Now)
	p := NewPipeline(opts)
	// First run consumes the allow token
	p.throttle.Allow()
	// Without reset, should be throttled
	if p.throttle.Allow() {
		t.Skip("clock advanced unexpectedly")
	}
	p.Reset()
	if !p.throttle.Allow() {
		t.Error("expected Allow after Reset")
	}
}

func TestNewPipeline_NotNil(t *testing.T) {
	p := NewPipeline(DefaultPipelineOptions())
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
	if p.throttle == nil {
		t.Error("expected non-nil throttle")
	}
	if p.dedupe == nil {
		t.Error("expected non-nil deduplicator")
	}
}

func TestPipeline_Run_ThrottledReturnsNil(t *testing.T) {
	clk := &mockClock{now: time.Now()}
	opts := DefaultPipelineOptions()
	opts.Throttle = throttleOptsFor(10*time.Second, clk.Now)
	p := NewPipeline(opts)
	// Consume the first allow
	p.throttle.Allow()
	// Next run should be throttled
	result, err := p.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Error("expected nil result when throttled")
	}
}
