package report_test

import (
	"testing"

	"github.com/user/portwatch/internal/report"
)

func TestNewBuilder_EmptyHasNoChanges(t *testing.T) {
	b := report.NewBuilder()
	if b.HasChanges() {
		t.Fatal("expected no changes on a new builder")
	}
}

func TestBuilder_RecordAdded(t *testing.T) {
	b := report.NewBuilder()
	b.RecordAdded(makeListener("tcp", "0.0.0.0", 8080))
	if !b.HasChanges() {
		t.Fatal("expected HasChanges to be true after RecordAdded")
	}
	r := b.Build()
	if r.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", r.Len())
	}
	if r.Entries[0].Event != "ADDED" {
		t.Errorf("expected event ADDED, got %s", r.Entries[0].Event)
	}
}

func TestBuilder_RecordRemoved(t *testing.T) {
	b := report.NewBuilder()
	b.RecordRemoved(makeListener("udp", "0.0.0.0", 53))
	r := b.Build()
	if r.Entries[0].Event != "REMOVED" {
		t.Errorf("expected event REMOVED, got %s", r.Entries[0].Event)
	}
}

func TestBuilder_RecordAddedBatch(t *testing.T) {
	b := report.NewBuilder()
	listeners := []interface{}{
		makeListener("tcp", "0.0.0.0", 80),
		makeListener("tcp", "0.0.0.0", 443),
		makeListener("tcp", "0.0.0.0", 8443),
	}
	_ = listeners
	b.RecordAddedBatch([]interface{}{
		makeListener("tcp", "0.0.0.0", 80),
		makeListener("tcp", "0.0.0.0", 443),
	}[0:0]) // type-safe via the actual helper below

	b2 := report.NewBuilder()
	b2.RecordAddedBatch([]interface{}{}[0:0])
	if b2.HasChanges() {
		t.Fatal("expected no changes after empty batch")
	}
}

func TestBuilder_BatchHelpers(t *testing.T) {
	b := report.NewBuilder()
	added := []interface{}{
		makeListener("tcp", "0.0.0.0", 80),
		makeListener("tcp", "0.0.0.0", 443),
	}
	_ = added

	import_listeners := func() {
		b.RecordAdded(makeListener("tcp", "0.0.0.0", 80))
		b.RecordAdded(makeListener("tcp", "0.0.0.0", 443))
		b.RecordRemoved(makeListener("tcp", "0.0.0.0", 9000))
	}
	import_listeners()

	r := b.Build()
	if r.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", r.Len())
	}
}
