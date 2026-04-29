package report

import (
	"github.com/user/portwatch/internal/ports"
)

// Builder accumulates diff results and produces a Report.
type Builder struct {
	report *Report
}

// NewBuilder returns a Builder backed by a fresh Report.
func NewBuilder() *Builder {
	return &Builder{report: New()}
}

// RecordAdded records a newly detected listener.
func (b *Builder) RecordAdded(l ports.Listener) {
	b.report.Add("ADDED", l)
}

// RecordRemoved records a listener that has disappeared.
func (b *Builder) RecordRemoved(l ports.Listener) {
	b.report.Add("REMOVED", l)
}

// RecordAddedBatch records multiple added listeners.
func (b *Builder) RecordAddedBatch(listeners []ports.Listener) {
	for _, l := range listeners {
		b.RecordAdded(l)
	}
}

// RecordRemovedBatch records multiple removed listeners.
func (b *Builder) RecordRemovedBatch(listeners []ports.Listener) {
	for _, l := range listeners {
		b.RecordRemoved(l)
	}
}

// HasChanges returns true when at least one entry has been recorded.
func (b *Builder) HasChanges() bool {
	return b.report.Len() > 0
}

// Build returns the completed Report.
func (b *Builder) Build() *Report {
	return b.report
}
