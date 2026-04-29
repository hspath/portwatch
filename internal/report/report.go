package report

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// Format defines the output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Entry represents a single port event in a report.
type Entry struct {
	Timestamp time.Time      `json:"timestamp"`
	Event     string         `json:"event"`
	Listener  ports.Listener `json:"listener"`
}

// Report holds a collection of port change entries.
type Report struct {
	GeneratedAt time.Time `json:"generated_at"`
	Entries     []Entry   `json:"entries"`
}

// New creates an empty Report with the current timestamp.
func New() *Report {
	return &Report{
		GeneratedAt: time.Now(),
		Entries:     []Entry{},
	}
}

// Add appends a new entry to the report.
func (r *Report) Add(event string, l ports.Listener) {
	r.Entries = append(r.Entries, Entry{
		Timestamp: time.Now(),
		Event:     event,
		Listener:  l,
	})
}

// Write serialises the report to w in the given format.
func (r *Report) Write(w io.Writer, f Format) error {
	if w == nil {
		w = os.Stdout
	}
	switch f {
	case FormatJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(r)
	case FormatText, "":
		fmt.Fprintf(w, "Port Watch Report — %s\n", r.GeneratedAt.Format(time.RFC3339))
		fmt.Fprintf(w, "%-10s %-8s %-22s %s\n", "EVENT", "PROTO", "ADDRESS", "TIMESTAMP")
		for _, e := range r.Entries {
			fmt.Fprintf(w, "%-10s %-8s %-22s %s\n",
				e.Event,
				e.Listener.Proto,
				fmt.Sprintf("%s:%d", e.Listener.IP, e.Listener.Port),
				e.Timestamp.Format(time.RFC3339),
			)
		}
		return nil
	default:
		return fmt.Errorf("unsupported report format: %q", f)
	}
}

// Len returns the number of entries in the report.
func (r *Report) Len() int { return len(r.Entries) }
