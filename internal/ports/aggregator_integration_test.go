package ports_test

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func makeAggListener(port uint16) ports.Listener {
	return ports.Listener{IP: net.ParseIP("0.0.0.0"), Port: port, Protocol: "tcp"}
}

// TestAggregator_RemovedTrackedSeparately ensures added and removed
// events are kept in separate slices within a single batch.
func TestAggregator_RemovedTrackedSeparately(t *testing.T) {
	now := time.Now()
	opts := ports.AggregatorOptions{
		Window:   60 * time.Second,
		MaxBatch: 10,
		Clock:    func() time.Time { return now },
	}
	a := ports.NewAggregator(opts)
	a.Record([]ports.Listener{makeAggListener(80)}, []ports.Listener{makeAggListener(22)})

	batch := a.Flush()
	if len(batch.Added) != 1 || batch.Added[0].Port != 80 {
		t.Errorf("unexpected added: %+v", batch.Added)
	}
	if len(batch.Removed) != 1 || batch.Removed[0].Port != 22 {
		t.Errorf("unexpected removed: %+v", batch.Removed)
	}
}

// TestAggregator_BatchTimestampSet verifies the At field is populated on flush.
func TestAggregator_BatchTimestampSet(t *testing.T) {
	expected := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	opts := ports.AggregatorOptions{
		Window:   60 * time.Second,
		MaxBatch: 100,
		Clock:    func() time.Time { return expected },
	}
	a := ports.NewAggregator(opts)
	a.Record([]ports.Listener{makeAggListener(443)}, nil)
	batch := a.Flush()

	if !batch.At.Equal(expected) {
		t.Errorf("expected At=%v, got %v", expected, batch.At)
	}
}
