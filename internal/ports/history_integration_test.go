package ports

import (
	"testing"
)

// TestHistory_EvictionPreservesNewest verifies that after overflow,
// the retained entries are the most recently pushed ones.
func TestHistory_EvictionPreservesNewest(t *testing.T) {
	h := NewHistory(2)
	h.Push(makeTestSnapshot(7001)) // will be evicted
	h.Push(makeTestSnapshot(7002)) // retained
	h.Push(makeTestSnapshot(7003)) // retained (latest)

	all := h.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if !all[0].Snapshot.Contains(Listener{Proto: "tcp", IP: "0.0.0.0", Port: 7002, PID: 1}) {
		t.Error("expected first retained entry to be port 7002")
	}
	if !all[1].Snapshot.Contains(Listener{Proto: "tcp", IP: "0.0.0.0", Port: 7003, PID: 1}) {
		t.Error("expected second retained entry to be port 7003")
	}
}

// TestHistory_LatestAfterClear confirms Latest returns false after Clear.
func TestHistory_LatestAfterClear(t *testing.T) {
	h := NewHistory(5)
	h.Push(makeTestSnapshot(8080))
	h.Clear()
	_, ok := h.Latest()
	if ok {
		t.Error("expected Latest to return false after Clear")
	}
}

// TestHistory_AllReturnsCopyNotReference ensures mutations to All() result
// do not affect internal state.
func TestHistory_AllReturnsCopyNotReference(t *testing.T) {
	h := NewHistory(5)
	h.Push(makeTestSnapshot(8080))
	all := h.All()
	all[0].Snapshot = makeTestSnapshot(9999)

	latest, _ := h.Latest()
	if latest.Snapshot.Contains(Listener{Proto: "tcp", IP: "0.0.0.0", Port: 9999, PID: 1}) {
		t.Error("mutating All() result should not affect internal history")
	}
}
