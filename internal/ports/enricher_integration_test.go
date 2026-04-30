package ports_test

import (
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func TestEnrichAll_SeverityConsistentWithClassify(t *testing.T) {
	listeners := []ports.Listener{
		{IP: "0.0.0.0", Port: 80, Protocol: "tcp"},
		{IP: "0.0.0.0", Port: 9999, Protocol: "tcp"},
		{IP: "0.0.0.0", Port: 500, Protocol: "tcp"},
	}
	enriched := ports.EnrichAll(listeners, true, ports.DefaultEnricherOptions())
	for i, el := range enriched {
		wantClass := ports.Classify(listeners[i])
		if el.Classification != wantClass {
			t.Errorf("listener %d: classification mismatch: want %q got %q", i, wantClass, el.Classification)
		}
		wantSev := ports.SeverityFor(listeners[i], true)
		if el.Severity != wantSev {
			t.Errorf("listener %d: severity mismatch: want %v got %v", i, wantSev, el.Severity)
		}
	}
}

func TestEnrich_RemovedAlwaysInfo(t *testing.T) {
	l := ports.Listener{IP: "0.0.0.0", Port: 443, Protocol: "tcp"}
	el := ports.Enrich(l, false, ports.DefaultEnricherOptions())
	if el.Severity != ports.SeverityFor(l, false) {
		t.Errorf("expected removed severity to match SeverityFor, got %v", el.Severity)
	}
}

func TestEnrich_HighPortIsWarning(t *testing.T) {
	l := ports.Listener{IP: "0.0.0.0", Port: 50000, Protocol: "tcp"}
	el := ports.Enrich(l, true, ports.DefaultEnricherOptions())
	if el.Severity != ports.SeverityWarning {
		t.Errorf("expected SeverityWarning for high port, got %v", el.Severity)
	}
}
