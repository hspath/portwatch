package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// EnrichedEvent pairs a change direction with an enriched listener.
type EnrichedEvent struct {
	Added    bool                   `json:"added"`
	Listener ports.EnrichedListener `json:"listener"`
	At       time.Time              `json:"at"`
}

// EnrichedReport holds enriched events for reporting.
type EnrichedReport struct {
	events []EnrichedEvent
}

// NewEnrichedReport creates an empty EnrichedReport.
func NewEnrichedReport() *EnrichedReport {
	return &EnrichedReport{}
}

// Add appends an enriched event.
func (r *EnrichedReport) Add(e EnrichedEvent) {
	r.events = append(r.events, e)
}

// Len returns the number of events.
func (r *EnrichedReport) Len() int {
	return len(r.events)
}

// WriteTo writes the report to w. format is "text" or "json".
func (r *EnrichedReport) WriteTo(w io.Writer, format string) error {
	switch format {
	case "json":
		return json.NewEncoder(w).Encode(r.events)
	default:
		for _, e := range r.events {
			direction := "+"
			if !e.Added {
				direction = "-"
			}
			_, err := fmt.Fprintf(w, "[%s] %s %s:%d (%s) sev=%s class=%s\n",
				e.At.Format(time.RFC3339),
				direction,
				e.Listener.IP,
				e.Listener.Port,
				e.Listener.ServiceName,
				e.Listener.Severity,
				e.Listener.Classification,
			)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
