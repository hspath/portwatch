package ports

import (
	"testing"
	"time"
)

func makeTestSnapshot(port uint16) *Snapshot {
	return NewSnapshot([]Listener{
		{Proto: "tcp", IP: "0.0.0.0", Port: port, PID: 1},
	})
}

func TestNewHistory_DefaultMaxSize(t *testing.T) {
	h := NewHistory(0)
	if h.maxSize != 10 {
		t.Errorf("expected default maxSize 10, got %d", h.maxSize)
	}
}

func TestHistory_PushAndLen(t *testing.T) {
	h := NewHistory(5)
	for i := 0; i < 3; i++ {
		h.Push(makeTestSnapshot(uint16(8000 + i)))
	}
	if h.Len() != 3 {
		t.Errorf("expected 3 entries, got %d", h.Len())
	}
}

func TestHistory_Eviction(t *testing.T) {
	h := NewHistory(3)
	for i := 0; i < 5; i++ {
		h.Push(makeTestSnapshot(uint16(9000 + i)))
	}
	if h.Len() != 3 {
		t.Errorf("expected 3 entries after eviction, got %d", h.Len())
	}
}

func TestHistory_Latest_Empty(t *testing.T) {
	h := NewHistory(5)
	_, ok := h.Latest()
	if ok {
		t.Error("expected Latest to return false on empty history")
	}
}

func TestHistory_Latest_ReturnsNewest(t *testing.T) {
	h := NewHistory(5)
	h.Push(makeTestSnapshot(8080))
	time.Sleep(time.Millisecond)
	h.Push(makeTestSnapshot(9090))

	entry, ok := h.Latest()
	if !ok {
		t.Fatal("expected Latest to return true")
	}
	if !entry.Snapshot.Contains(Listener{Proto: "tcp", IP: "0.0.0.0", Port: 9090, PID: 1}) {
		t.Error("expected latest snapshot to contain port 9090")
	}
}

func TestHistory_All_OrderPreserved(t *testing.T) {
	h := NewHistory(5)
	ports := []uint16{8001, 8002, 8003}
	for _, p := range ports {
		h.Push(makeTestSnapshot(p))
	}
	all := h.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(all))
	}
	for i, entry := range all {
		if !entry.Snapshot.Contains(Listener{Proto: "tcp", IP: "0.0.0.0", Port: ports[i], PID: 1}) {
			t.Errorf("entry %d: unexpected snapshot content", i)
		}
	}
}

func TestHistory_Clear(t *testing.T) {
	h := NewHistory(5)
	h.Push(makeTestSnapshot(8080))
	h.Push(makeTestSnapshot(9090))
	h.Clear()
	if h.Len() != 0 {
		t.Errorf("expected 0 entries after Clear, got %d", h.Len())
	}
}
