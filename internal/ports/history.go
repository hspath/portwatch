package ports

import (
	"sync"
	"time"
)

// HistoryEntry records a snapshot at a point in time.
type HistoryEntry struct {
	Timestamp time.Time
	Snapshot  *Snapshot
}

// History maintains a rolling window of snapshots.
type History struct {
	mu      sync.RWMutex
	entries []HistoryEntry
	maxSize int
}

// NewHistory creates a History that retains at most maxSize entries.
func NewHistory(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 10
	}
	return &History{maxSize: maxSize}
}

// Push appends a new snapshot, evicting the oldest if at capacity.
func (h *History) Push(s *Snapshot) {
	h.mu.Lock()
	defer h.mu.Unlock()
	entry := HistoryEntry{Timestamp: time.Now(), Snapshot: s}
	if len(h.entries) >= h.maxSize {
		h.entries = h.entries[1:]
	}
	h.entries = append(h.entries, entry)
}

// Len returns the number of stored entries.
func (h *History) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.entries)
}

// Latest returns the most recent entry, and false if empty.
func (h *History) Latest() (HistoryEntry, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.entries) == 0 {
		return HistoryEntry{}, false
	}
	return h.entries[len(h.entries)-1], true
}

// All returns a copy of all stored entries, oldest first.
func (h *History) All() []HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]HistoryEntry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Clear removes all entries.
func (h *History) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = nil
}
