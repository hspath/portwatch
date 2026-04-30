package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func makeEnrichedListener(ip string, port uint16, proto string) ports.EnrichedListener {
	l := ports.Listener{IP: ip, Port: port, Protocol: proto}
	return ports.Enrich(l, true, ports.DefaultEnricherOptions())
}

func makeEnrichedEvent(added bool) EnrichedEvent {
	return EnrichedEvent{
		Added:    added,
		Listener: makeEnrichedListener("0.0.0.0", 22, "tcp"),
		At:       time.Now(),
	}
}

func TestEnrichedReport_EmptyLen(t *testing.T) {
	r := NewEnrichedReport()
	if r.Len() != 0 {
		t.Errorf("expected 0, got %d", r.Len())
	}
}

func TestEnrichedReport_Add_IncreasesLen(t *testing.T) {
	r := NewEnrichedReport()
	r.Add(makeEnrichedEvent(true))
	r.Add(makeEnrichedEvent(false))
	if r.Len() != 2 {
		t.Errorf("expected 2, got %d", r.Len())
	}
}

func TestEnrichedReport_WriteTo_TextContainsPlus(t *testing.T) {
	r := NewEnrichedReport()
	r.Add(makeEnrichedEvent(true))
	var buf bytes.Buffer
	if err := r.WriteTo(&buf, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "+") {
		t.Errorf("expected '+' in output, got: %s", buf.String())
	}
}

func TestEnrichedReport_WriteTo_TextContainsMinus(t *testing.T) {
	r := NewEnrichedReport()
	r.Add(makeEnrichedEvent(false))
	var buf bytes.Buffer
	if err := r.WriteTo(&buf, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "-") {
		t.Errorf("expected '-' in output, got: %s", buf.String())
	}
}

func TestEnrichedReport_WriteTo_JSONValid(t *testing.T) {
	r := NewEnrichedReport()
	r.Add(makeEnrichedEvent(true))
	var buf bytes.Buffer
	if err := r.WriteTo(&buf, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []EnrichedEvent
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Errorf("invalid JSON: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 event, got %d", len(out))
	}
}

func TestEnrichedReport_WriteTo_ServiceNameInText(t *testing.T) {
	r := NewEnrichedReport()
	r.Add(makeEnrichedEvent(true))
	var buf bytes.Buffer
	_ = r.WriteTo(&buf, "text")
	if !strings.Contains(buf.String(), "ssh") {
		t.Errorf("expected service name 'ssh' in output, got: %s", buf.String())
	}
}
