package ports

import (
	"fmt"
	"time"
)

// Snapshot holds a point-in-time capture of active listeners.
type Snapshot struct {
	CapturedAt time.Time
	Listeners  []Listener
}

// NewSnapshot scans current listeners and returns a Snapshot.
func NewSnapshot(opts FilterOptions) (*Snapshot, error) {
	listeners, err := ScanListeners()
	if err != nil {
		return nil, fmt.Errorf("snapshot: scan failed: %w", err)
	}
	filtered := ApplyFilter(listeners, opts)
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		Listeners:  filtered,
	}, nil
}

// Summary returns a human-readable one-line description of the snapshot.
func (s *Snapshot) Summary() string {
	return fmt.Sprintf("snapshot at %s: %d listener(s)",
		s.CapturedAt.Format(time.RFC3339), len(s.Listeners))
}

// Contains reports whether the snapshot includes a listener matching
// the given protocol, address, and port.
func (s *Snapshot) Contains(proto, addr string, port uint16) bool {
	for _, l := range s.Listeners {
		if l.Proto == proto && l.Addr == addr && l.Port == port {
			return true
		}
	}
	return false
}
