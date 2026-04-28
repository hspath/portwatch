package ports

import (
	"testing"
	"time"
)

func TestDefaultWatchOptions(t *testing.T) {
	opts := DefaultWatchOptions()
	if opts.Interval != 15*time.Second {
		t.Errorf("expected 15s interval, got %v", opts.Interval)
	}
}

func TestDiff_AddedAndRemoved(t *testing.T) {
	a := makeListener("tcp", "0.0.0.0", 80)
	b := makeListener("tcp", "0.0.0.0", 443)
	c := makeListener("tcp", "0.0.0.0", 8080)

	prev := &Snapshot{Listeners: []Listener{a, b}}
	curr := &Snapshot{Listeners: []Listener{b, c}}

	added, removed := diff(prev, curr)

	if len(added) != 1 || added[0].Port != 8080 {
		t.Errorf("expected added=[8080], got %v", added)
	}
	if len(removed) != 1 || removed[0].Port != 80 {
		t.Errorf("expected removed=[80], got %v", removed)
	}
}

func TestDiff_NoChange(t *testing.T) {
	a := makeListener("tcp", "0.0.0.0", 80)
	prev := &Snapshot{Listeners: []Listener{a}}
	curr := &Snapshot{Listeners: []Listener{a}}

	added, removed := diff(prev, curr)
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", added, removed)
	}
}

func TestDiff_EmptyPrev(t *testing.T) {
	a := makeListener("tcp", "127.0.0.1", 9090)
	prev := &Snapshot{}
	curr := &Snapshot{Listeners: []Listener{a}}

	added, removed := diff(prev, curr)
	if len(added) != 1 {
		t.Errorf("expected 1 added, got %d", len(added))
	}
	if len(removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(removed))
	}
}

func TestDiff_EmptyCurr(t *testing.T) {
	a := makeListener("udp", "0.0.0.0", 53)
	prev := &Snapshot{Listeners: []Listener{a}}
	curr := &Snapshot{}

	added, removed := diff(prev, curr)
	if len(added) != 0 {
		t.Errorf("expected 0 added, got %d", len(added))
	}
	if len(removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(removed))
	}
}
