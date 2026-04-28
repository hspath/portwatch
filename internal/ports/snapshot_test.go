package ports

import (
	"testing"
	"time"
)

func makeSnapshot(listeners []Listener) *Snapshot {
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		Listeners:  listeners,
	}
}

func TestSnapshot_Summary_ContainsCount(t *testing.T) {
	s := makeSnapshot([]Listener{
		{Proto: "tcp", Addr: "0.0.0.0", Port: 80},
		{Proto: "tcp", Addr: "0.0.0.0", Port: 443},
	})
	summary := s.Summary()
	if summary == "" {
		t.Fatal("expected non-empty summary")
	}
	// Should mention the count
	for _, sub := range []string{"2", "listener"} {
		if !containsStr(summary, sub) {
			t.Errorf("summary %q missing %q", summary, sub)
		}
	}
}

func TestSnapshot_Contains_Match(t *testing.T) {
	s := makeSnapshot([]Listener{
		{Proto: "tcp", Addr: "0.0.0.0", Port: 8080},
	})
	if !s.Contains("tcp", "0.0.0.0", 8080) {
		t.Error("expected Contains to return true for existing listener")
	}
}

func TestSnapshot_Contains_NoMatch(t *testing.T) {
	s := makeSnapshot([]Listener{
		{Proto: "tcp", Addr: "0.0.0.0", Port: 8080},
	})
	if s.Contains("udp", "0.0.0.0", 8080) {
		t.Error("expected Contains to return false for mismatched proto")
	}
	if s.Contains("tcp", "0.0.0.0", 9090) {
		t.Error("expected Contains to return false for mismatched port")
	}
}

func TestSnapshot_Contains_Empty(t *testing.T) {
	s := makeSnapshot(nil)
	if s.Contains("tcp", "0.0.0.0", 80) {
		t.Error("expected Contains to return false on empty snapshot")
	}
}

func TestSnapshot_Summary_EmptyListeners(t *testing.T) {
	s := makeSnapshot(nil)
	summary := s.Summary()
	if !containsStr(summary, "0") {
		t.Errorf("expected '0' in summary, got %q", summary)
	}
}

// containsStr is a simple helper to avoid importing strings in test file.
func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && func() bool {
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	}()
}
