package ports

import (
	"os"
	"testing"
)

func TestParseHexAddr_Valid(t *testing.T) {
	addr, port, err := parseHexAddr("0100007F:0050")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if port != 80 {
		t.Errorf("expected port 80, got %d", port)
	}
	if addr != "0100007F" {
		t.Errorf("expected addr 0100007F, got %s", addr)
	}
}

func TestParseHexAddr_InvalidFormat(t *testing.T) {
	_, _, err := parseHexAddr("BADFORMAT")
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

func TestParseHexAddr_InvalidPort(t *testing.T) {
	_, _, err := parseHexAddr("0100007F:ZZZZ")
	if err == nil {
		t.Fatal("expected error for invalid port hex, got nil")
	}
}

func TestListenerString(t *testing.T) {
	l := Listener{
		Protocol: "tcp",
		Address:  "127.0.0.1",
		Port:     8080,
		PID:      1234,
		Process:  "myapp",
	}
	got := l.String()
	expected := "tcp 127.0.0.1:8080 (pid=1234, process=myapp)"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestParseProcNet_MissingFile(t *testing.T) {
	_, err := parseProcNet("/nonexistent/path", "tcp")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseProcNet_SyntheticFile(t *testing.T) {
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1 0000000000000000 100 0 0 10 0
   1: 0100007F:0050 00000000:0000 06 00000000:00000000 00:00000000 00000000     0        0 67890 1 0000000000000000 100 0 0 10 0
`
	tmpFile, err := os.CreateTemp("", "proc_net_tcp_*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	listeners, err := parseProcNet(tmpFile.Name(), "tcp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only state 0A (LISTEN) should be included; state 06 should be skipped.
	if len(listeners) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(listeners))
	}
	if listeners[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", listeners[0].Port)
	}
}
