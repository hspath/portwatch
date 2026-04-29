package ports

import "testing"

func makeSnapshot(listeners ...Listener) *Snapshot {
	return NewSnapshot(listeners)
}

func TestSnapshot_Summary_ContainsCount(t *testing.T) {
	s := makeSnapshot(
		Listener{Proto: "tcp", IP: "0.0.0.0", Port: 80},
		Listener{Proto: "tcp", IP: "0.0.0.0", Port: 443},
	)
	sum := s.Summary()
	if sum == "" {
		t.Fatal("expected non-empty summary")
	}
	if s.Len() != 2 {
		t.Errorf("expected Len 2, got %d", s.Len())
	}
}

func TestSnapshot_Contains_Match(t *testing.T) {
	l := Listener{Proto: "tcp", IP: "127.0.0.1", Port: 8080}
	s := makeSnapshot(l)
	if !s.Contains(l) {
		t.Error("expected snapshot to contain listener")
	}
}

func TestSnapshot_Contains_NoMatch(t *testing.T) {
	s := makeSnapshot(Listener{Proto: "tcp", IP: "0.0.0.0", Port: 80})
	other := Listener{Proto: "udp", IP: "0.0.0.0", Port: 80}
	if s.Contains(other) {
		t.Error("expected snapshot NOT to contain listener with different proto")
	}
}

func TestSnapshot_Contains_Empty(t *testing.T) {
	s := makeSnapshot()
	if s.Contains(Listener{Proto: "tcp", IP: "0.0.0.0", Port: 22}) {
		t.Error("empty snapshot should not contain any listener")
	}
}

func TestSnapshot_Listeners_IsCopy(t *testing.T) {
	original := Listener{Proto: "tcp", IP: "0.0.0.0", Port: 22}
	s := makeSnapshot(original)
	ls := s.Listeners()
	ls[0].Port = 9999
	if s.listeners[0].Port == 9999 {
		t.Error("mutating returned slice should not affect snapshot")
	}
}

func TestSnapshot_Len_Empty(t *testing.T) {
	s := makeSnapshot()
	if s.Len() != 0 {
		t.Errorf("expected Len 0, got %d", s.Len())
	}
}
