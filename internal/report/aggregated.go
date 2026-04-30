package report

import (
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// AggregatedReport summarises a batch of port-change events.
type AggregatedReport struct {
	Batch ports.AggregatedBatch
	Hostname string
}

// NewAggregatedReport creates an AggregatedReport for the given batch.
func NewAggregatedReport(batch ports.AggregatedBatch, hostname string) AggregatedReport {
	return AggregatedReport{Batch: batch, Hostname: hostname}
}

// WriteTo writes a human-readable summary of the batch to w.
func (ar AggregatedReport) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w,
		"[%s] portwatch batch report — host: %s | added: %d | removed: %d\n",
		ar.Batch.At.UTC().Format(time.RFC3339),
		ar.Hostname,
		len(ar.Batch.Added),
		len(ar.Batch.Removed),
	)
	if err != nil {
		return int64(n), err
	}
	total := int64(n)

	for _, l := range ar.Batch.Added {
		n2, err := fmt.Fprintf(w, "  + %s\n", l)
		total += int64(n2)
		if err != nil {
			return total, err
		}
	}
	for _, l := range ar.Batch.Removed {
		n2, err := fmt.Fprintf(w, "  - %s\n", l)
		total += int64(n2)
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

// Summary returns a one-line description of the batch.
func (ar AggregatedReport) Summary() string {
	return fmt.Sprintf("batch at %s: +%d/-%d listeners on %s",
		ar.Batch.At.UTC().Format(time.RFC3339),
		len(ar.Batch.Added),
		len(ar.Batch.Removed),
		ar.Hostname,
	)
}
