package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/ports"
)

func listener(proto, ip string, port uint16) ports.Listener {
	return ports.Listener{Proto: proto, IP: ip, Port: port}
}

func TestUnexpected_OutputContainsWARN(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	a := n.Unexpected(listener("tcp", "0.0.0.0", 9999))

	if a.Level != alert.LevelWarn {
		t.Errorf("expected level WARN, got %s", a.Level)
	}
	if !strings.Contains(buf.String(), "WARN") {
		t.Errorf("output missing WARN: %q", buf.String())
	}
	if !strings.Contains(buf.String(), "9999") {
		t.Errorf("output missing port 9999: %q", buf.String())
	}
}

func TestGone_OutputContainsINFO(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	a := n.Gone(listener("tcp", "127.0.0.1", 22))

	if a.Level != alert.LevelInfo {
		t.Errorf("expected level INFO, got %s", a.Level)
	}
	if !strings.Contains(buf.String(), "INFO") {
		t.Errorf("output missing INFO: %q", buf.String())
	}
}

func TestAlertString_Format(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	n.Unexpected(listener("udp", "0.0.0.0", 53))
	out := buf.String()

	for _, want := range []string{"udp", "0.0.0.0", "53", "unexpected listener"} {
		if !strings.Contains(out, want) {
			t.Errorf("alert string missing %q in: %q", want, out)
		}
	}
}

func TestNew_NilWriterUsesStderr(t *testing.T) {
	// Should not panic when w is nil.
	n := alert.New(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}
