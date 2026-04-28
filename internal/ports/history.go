package ports

import (
	"fmt"
	"sync"
)

// History keeps an ordered ring of recent Snapshots for trend analysis.
type History struct {
	mu       sync.Mutex
	max      int
	snapshots []*Snapshot
}

// NewHistory creates a History that retains at most maxEntries snapshots.
func NewHistory(maxEntries int) *History {
	if maxEntries < 1 {
		maxEntries = 1
	}
	return &History{max: maxEntries}
}

// Add appends a snapshot, evicting the oldest if the buffer is full.
func (h *History) Add(s *Snapshot) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.snapshots) >= h.max {
		h.snapshots = h.snapshots[1:]
	}
	h.snapshots = append(h.snapshots, s)
}

// Latest returns the most recently added snapshot, or nil if empty.
func (h *History) Latest() *Snapshot {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.snapshots) == 0 {
		return nil
	}
	return h.snapshots[len(h.snapshots)-1]
}

// Len returns the number of snapshots currently stored.
func (h *History) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.snapshots)
}

// All returns a copy of all stored snapshots in chronological order.
func (h *History) All() []*Snapshot {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]*Snapshot, len(h.snapshots))
	copy(out, h.snapshots)
	return out
}

// String returns a brief description of the history buffer state.
func (h *History) String() string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return fmt.Sprintf("history: %d/%d snapshots", len(h.snapshots), h.max)
}
