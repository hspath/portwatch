package ports

import (
	"fmt"
	"time"
)

// EventKind describes whether a port was opened or closed.
type EventKind string

const (
	EventAdded   EventKind = "added"
	EventRemoved EventKind = "removed"
)

// Event captures a single port change detected during a scan cycle.
type Event struct {
	Kind      EventKind
	Listener  Listener
	Severity  Severity
	Timestamp time.Time
}

// NewEvent constructs an Event, deriving Severity automatically.
func NewEvent(kind EventKind, l Listener, ts time.Time) Event {
	return Event{
		Kind:      kind,
		Listener:  l,
		Severity:  SeverityFor(l, kind == EventAdded),
		Timestamp: ts,
	}
}

// String returns a compact, human-readable representation of the event.
func (e Event) String() string {
	return fmt.Sprintf(
		"[%s] %s %s/%d (%s)",
		e.Severity,
		e.Kind,
		e.Listener.Protocol,
		e.Listener.Port,
		e.Timestamp.Format(time.RFC3339),
	)
}

// EventsFromDiff converts added/removed listener slices into a slice of Events
// using the provided timestamp.
func EventsFromDiff(added, removed []Listener, ts time.Time) []Event {
	events := make([]Event, 0, len(added)+len(removed))
	for _, l := range added {
		events = append(events, NewEvent(EventAdded, l, ts))
	}
	for _, l := range removed {
		events = append(events, NewEvent(EventRemoved, l, ts))
	}
	return events
}
