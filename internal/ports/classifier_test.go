package ports

import (
	"testing"
)

func TestClassify_WellKnownPort(t *testing.T) {
	l := Listener{Port: 22, Proto: "tcp"}
	c := Classify(l)
	if !c.WellKnown {
		t.Error("expected WellKnown=true for port 22")
	}
	if c.ServiceName != "ssh" {
		t.Errorf("expected ServiceName=ssh, got %s", c.ServiceName)
	}
	if !c.Privileged {
		t.Error("expected Privileged=true for port 22")
	}
}

func TestClassify_UnknownPort(t *testing.T) {
	l := Listener{Port: 9999, Proto: "tcp"}
	c := Classify(l)
	if c.WellKnown {
		t.Error("expected WellKnown=false for port 9999")
	}
	if c.ServiceName != "unknown" {
		t.Errorf("expected ServiceName=unknown, got %s", c.ServiceName)
	}
	if c.Privileged {
		t.Error("expected Privileged=false for port 9999")
	}
}

func TestClassify_PrivilegedUnknownPort(t *testing.T) {
	l := Listener{Port: 999, Proto: "tcp"}
	c := Classify(l)
	if !c.Privileged {
		t.Error("expected Privileged=true for port 999")
	}
	if c.WellKnown {
		t.Error("expected WellKnown=false for port 999")
	}
}

func TestClassify_HighPort_NotPrivileged(t *testing.T) {
	l := Listener{Port: 8080, Proto: "tcp"}
	c := Classify(l)
	if c.Privileged {
		t.Error("expected Privileged=false for port 8080")
	}
	if c.ServiceName != "http-alt" {
		t.Errorf("expected ServiceName=http-alt, got %s", c.ServiceName)
	}
}

func TestClassifyAll_ReturnsAllListeners(t *testing.T) {
	listeners := []Listener{
		{Port: 80, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
		{Port: 12345, Proto: "tcp"},
	}
	result := ClassifyAll(listeners)
	if len(result) != 3 {
		t.Errorf("expected 3 classifications, got %d", len(result))
	}
	if result[listeners[0]].ServiceName != "http" {
		t.Errorf("expected http for port 80, got %s", result[listeners[0]].ServiceName)
	}
	if result[listeners[2]].WellKnown {
		t.Error("expected WellKnown=false for port 12345")
	}
}
