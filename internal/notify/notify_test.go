package notify

import (
	"bytes"
	"strings"
	"testing"
)

func TestSend_Stdout(t *testing.T) {
	var buf bytes.Buffer
	n := New(Config{Method: MethodStdout}, &buf)
	if err := n.Send("hello portwatch"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "hello portwatch") {
		t.Errorf("expected output to contain message, got: %q", buf.String())
	}
}

func TestSend_Log(t *testing.T) {
	var buf bytes.Buffer
	n := New(Config{Method: MethodLog}, &buf)
	if err := n.Send("test message"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "test message") {
		t.Errorf("expected output to contain message, got: %q", buf.String())
	}
}

func TestSend_EmptyMethod_DefaultsToStdout(t *testing.T) {
	var buf bytes.Buffer
	n := New(Config{}, &buf)
	if err := n.Send("default"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "default") {
		t.Errorf("expected output, got: %q", buf.String())
	}
}

func TestSend_UnknownMethod_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	n := New(Config{Method: "sms"}, &buf)
	err := n.Send("msg")
	if err == nil {
		t.Fatal("expected error for unknown method, got nil")
	}
	if !strings.Contains(err.Error(), "unknown method") {
		t.Errorf("unexpected error text: %v", err)
	}
}

func TestSend_Exec_MissingCommand_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	n := New(Config{Method: MethodExec, Command: ""}, &buf)
	err := n.Send("msg")
	if err == nil {
		t.Fatal("expected error when command is empty")
	}
}

func TestSend_Exec_ValidCommand(t *testing.T) {
	var buf bytes.Buffer
	n := New(Config{Method: MethodExec, Command: "echo"}, &buf)
	if err := n.Send("ping"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "ping") {
		t.Errorf("expected echo output to contain 'ping', got: %q", buf.String())
	}
}

func TestNew_NilWriter_UsesStdout(t *testing.T) {
	n := New(Config{Method: MethodLog}, nil)
	if n.out == nil {
		t.Error("expected non-nil writer when nil passed to New")
	}
}
