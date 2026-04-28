package notify

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Method represents a notification delivery method.
type Method string

const (
	MethodLog    Method = "log"
	MethodStdout Method = "stdout"
	MethodExec   Method = "exec"
)

// Config holds configuration for a notifier.
type Config struct {
	Method  Method `json:"method"`
	Command string `json:"command,omitempty"` // used when Method == MethodExec
}

// Notifier sends notification messages via a configured method.
type Notifier struct {
	cfg Config
	out io.Writer
}

// New creates a Notifier. If out is nil, os.Stdout is used.
func New(cfg Config, out io.Writer) *Notifier {
	if out == nil {
		out = os.Stdout
	}
	return &Notifier{cfg: cfg, out: out}
}

// Send delivers msg via the configured method.
func (n *Notifier) Send(msg string) error {
	switch n.cfg.Method {
	case MethodExec:
		return n.execCommand(msg)
	case MethodStdout, MethodLog, "":
		_, err := fmt.Fprintln(n.out, msg)
		return err
	default:
		return fmt.Errorf("notify: unknown method %q", n.cfg.Method)
	}
}

// execCommand runs cfg.Command with msg appended as the final argument.
func (n *Notifier) execCommand(msg string) error {
	if n.cfg.Command == "" {
		return fmt.Errorf("notify: exec method requires a command")
	}
	parts := strings.Fields(n.cfg.Command)
	args := append(parts[1:], msg)
	cmd := exec.Command(parts[0], args...)
	cmd.Stdout = n.out
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
