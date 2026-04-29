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

func listenerKey(l Listener) string {
	return fmt.Sprintf("%s|%s|%d", l.Proto, l.IP, l.Port)
}
