package ports

import (
	"testing"
)

func TestSeverity_String(t *testing.T) {
	tests := []struct {
		s    Severity
		want string
	}{
		{SeverityInfo, "INFO"},
		{SeverityWarning, "WARNING"},
		{SeverityCritical, "CRITICAL"},
		{Severity(99), "UNKNOWN"},
	}
	for _, tc := range tests {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}

func TestSeverityFor_Removed_AlwaysInfo(t *testing.T) {
	listeners := []Listener{
		{Port: 22, Protocol: "tcp"},
		{Port: 8080, Protocol: "tcp"},
		{Port: 443, Protocol: "tcp"},
	}
	for _, l := range listeners {
		if got := SeverityFor(l, false); got != SeverityInfo {
			t.Errorf("SeverityFor(port=%d, added=false) = %v, want INFO", l.Port, got)
		}
	}
}

func TestSeverityFor_Added_WellKnown_IsCritical(t *testing.T) {
	l := Listener{Port: 22, Protocol: "tcp"} // SSH — well-known
	if got := SeverityFor(l, true); got != SeverityCritical {
		t.Errorf("SeverityFor(port=22, added=true) = %v, want CRITICAL", got)
	}
}

func TestSeverityFor_Added_PrivilegedUnknown_IsCritical(t *testing.T) {
	// Port 999 is privileged (<1024) but not in the well-known list.
	l := Listener{Port: 999, Protocol: "tcp"}
	if got := SeverityFor(l, true); got != SeverityCritical {
		t.Errorf("SeverityFor(port=999, added=true) = %v, want CRITICAL", got)
	}
}

func TestSeverityFor_Added_HighPort_IsWarning(t *testing.T) {
	l := Listener{Port: 49152, Protocol: "tcp"}
	if got := SeverityFor(l, true); got != SeverityWarning {
		t.Errorf("SeverityFor(port=49152, added=true) = %v, want WARNING", got)
	}
}

func TestSeverityFor_Added_UnknownHighPort_IsWarning(t *testing.T) {
	l := Listener{Port: 9999, Protocol: "tcp"}
	if got := SeverityFor(l, true); got != SeverityWarning {
		t.Errorf("SeverityFor(port=9999, added=true) = %v, want WARNING", got)
	}
}
