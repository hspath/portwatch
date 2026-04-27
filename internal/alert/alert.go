package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelWarn  Level = "WARN"
	LevelInfo  Level = "INFO"
)

// Alert describes a single port-change event.
type Alert struct {
	Timestamp time.Time    `json:"timestamp"`
	Level     Level        `json:"level"`
	Message   string       `json:"message"`
	Listener  ports.Listener `json:"listener"`
}

func (a Alert) String() string {
	return fmt.Sprintf("%s [%s] %s (%s)",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Message,
		a.Listener.String(),
	)
}

// Notifier writes alerts to an output sink.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to w. Pass nil to use os.Stderr.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stderr
	}
	return &Notifier{out: w}
}

// Unexpected emits a WARN alert for a newly discovered listener.
func (n *Notifier) Unexpected(l ports.Listener) Alert {
	a := Alert{
		Timestamp: time.Now().UTC(),
		Level:     LevelWarn,
		Message:   "unexpected listener detected",
		Listener:  l,
	}
	fmt.Fprintln(n.out, a.String())
	return a
}

// Gone emits an INFO alert for a listener that has disappeared.
func (n *Notifier) Gone(l ports.Listener) Alert {
	a := Alert{
		Timestamp: time.Now().UTC(),
		Level:     LevelInfo,
		Message:   "listener no longer present",
		Listener:  l,
	}
	fmt.Fprintln(n.out, a.String())
	return a
}
