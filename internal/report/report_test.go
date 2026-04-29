package report_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/report"
)

func makeListener(proto, ip string, port uint16) ports.Listener {
	return ports.Listener{Proto: proto, IP: ip, Port: port}
}

func TestNew_EmptyReport(t *testing.T) {
	r := report.New()
	if r.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", r.Len())
	}
	if r.GeneratedAt.IsZero() {
		t.Fatal("expected GeneratedAt to be set")
	}
}

func TestAdd_IncreasesLen(t *testing.T) {
	r := report.New()
	r.Add("ADDED", makeListener("tcp", "0.0.0.0", 8080))
	r.Add("REMOVED", makeListener("tcp", "0.0.0.0", 9090))
	if r.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", r.Len())
	}
}

func TestWrite_TextFormat_ContainsEvent(t *testing.T) {
	r := report.New()
	r.Add("ADDED", makeListener("tcp", "127.0.0.1", 3000))

	var buf bytes.Buffer
	if err := r.Write(&buf, report.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ADDED") {
		t.Errorf("expected output to contain ADDED, got:\n%s", out)
	}
	if !strings.Contains(out, "3000") {
		t.Errorf("expected output to contain port 3000, got:\n%s", out)
	}
}

func TestWrite_JSONFormat_ValidJSON(t *testing.T) {
	r := report.New()
	r.Add("REMOVED", makeListener("udp", "0.0.0.0", 53))

	var buf bytes.Buffer
	if err := r.Write(&buf, report.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded report.Report
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(decoded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(decoded.Entries))
	}
	if decoded.Entries[0].Event != "REMOVED" {
		t.Errorf("expected event REMOVED, got %s", decoded.Entries[0].Event)
	}
}

func TestWrite_UnknownFormat_ReturnsError(t *testing.T) {
	r := report.New()
	var buf bytes.Buffer
	if err := r.Write(&buf, report.Format("xml")); err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestWrite_NilWriter_UsesStdout(t *testing.T) {
	r := report.New()
	r.Add("ADDED", makeListener("tcp", "0.0.0.0", 80))
	// Should not panic with nil writer
	if err := r.Write(nil, report.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
