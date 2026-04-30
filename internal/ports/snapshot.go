package ports

import "fmt"

// Snapshot is an immutable set of listeners captured at one moment.
type Snapshot struct {
	listeners []Listener
	index     map[string]Listener
}

// NewSnapshot creates a Snapshot from a slice of Listeners.
func NewSnapshot(ls []Listener) *Snapshot {
	idx := make(map[string]Listener, len(ls))
	for _, l := range ls {
		idx[listenerKey(l)] = l
	}
	out := make([]Listener, len(ls))
	copy(out, ls)
	return &Snapshot{listeners: out, index: idx}
}

// Listeners returns a copy of all listeners in the snapshot.
func (s *Snapshot) Listeners() []Listener {
	out := make([]Listener, len(s.listeners))
	copy(out, s.listeners)
	return out
}

// Contains reports whether l is present in the snapshot.
func (s *Snapshot) Contains(l Listener) bool {
	_, ok := s.index[listenerKey(l)]
	return ok
}

// Len returns the number of listeners.
func (s *Snapshot) Len() int { return len(s.listeners) }

// Summary returns a human-readable one-liner.
func (s *Snapshot) Summary() string {
	return fmt.Sprintf("snapshot: %d listener(s)", len(s.listeners))
}

// Diff returns the listeners added and removed between s and a newer snapshot.
// added contains listeners present in next but not in s.
// removed contains listeners present in s but not in next.
func (s *Snapshot) Diff(next *Snapshot) (added, removed []Listener) {
	for _, l := range next.listeners {
		if !s.Contains(l) {
			added = append(added, l)
		}
	}
	for _, l := range s.listeners {
		if !next.Contains(l) {
			removed = append(removed, l)
		}
	}
	return added, removed
}

func listenerKey(l Listener) string {
	return fmt.Sprintf("%s|%s|%d", l.Proto, l.IP, l.Port)
}
