package report

import (
	"bytes"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func makeBatch(added, removed []ports.Listener) ports.AggregatedBatch {
	return ports.AggregatedBatch{
		Added:   added,
		Removed: removed,
		At:      time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
	}
}

func batchListener(port uint16) ports.Listener {
	return ports.Listener{IP: net.ParseIP("0.0.0.0"), Port: port, Protocol: "tcp"}
}

func TestAggregatedReport_Summary_ContainsCounts(t *testing.T) {
	batch := makeBatch(
		[]ports.Listener{batchListener(80), batchListener(443)},
		[]ports.Listener{batchListener(22)},
	)
	ar := NewAggregatedReport(batch, "testhost")
	s := ar.Summary()

	if !strings.Contains(s, "+2") {
		t.Errorf("summary missing added count: %s", s)
	}
	if !strings.Contains(s, "-1") {
		t.Errorf("summary missing removed count: %s", s)
	}
	if !strings.Contains(s, "testhost") {
		t.Errorf("summary missing hostname: %s", s)
	}
}

func TestAggregatedReport_WriteTo_ContainsPlusAndMinus(t *testing.T) {
	batch := makeBatch(
		[]ports.Listener{batchListener(8080)},
		[]ports.Listener{batchListener(9090)},
	)
	ar := NewAggregatedReport(batch, "myhost")

	var buf bytes.Buffer
	_, err := ar.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "+ ") {
		t.Errorf("output missing added marker: %s", out)
	}
	if !strings.Contains(out, "- ") {
		t.Errorf("output missing removed marker: %s", out)
	}
}

func TestAggregatedReport_WriteTo_EmptyBatch(t *testing.T) {
	batch := makeBatch(nil, nil)
	ar := NewAggregatedReport(batch, "host")

	var buf bytes.Buffer
	_, err := ar.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo error: %v", err)
	}
	if !strings.Contains(buf.String(), "portwatch batch report") {
		t.Error("expected header line even for empty batch")
	}
}

func TestAggregatedReport_Summary_ContainsTimestamp(t *testing.T) {
	batch := makeBatch(nil, nil)
	ar := NewAggregatedReport(batch, "host")
	if !strings.Contains(ar.Summary(), "2024-06-01") {
		t.Errorf("summary missing date: %s", ar.Summary())
	}
}
