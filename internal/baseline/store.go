package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/portwatch/internal/ports"
)

// Baseline holds the set of expected (approved) listeners.
type Baseline struct {
	Listeners []ports.Listener `json:"listeners"`
}

// Load reads a baseline file from disk. Returns an empty Baseline if the
// file does not exist yet.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Baseline{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}

	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: parse %s: %w", path, err)
	}
	return &b, nil
}

// Save writes the baseline to disk, creating parent directories as needed.
func Save(path string, b *Baseline) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("baseline: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", path, err)
	}
	return nil
}

// Diff compares a live snapshot against the baseline and returns two slices:
// added  – listeners present in current but not in the baseline.
// removed – listeners present in the baseline but not in current.
func Diff(base *Baseline, current []ports.Listener) (added, removed []ports.Listener) {
	baseSet := make(map[string]struct{}, len(base.Listeners))
	for _, l := range base.Listeners {
		baseSet[l.String()] = struct{}{}
	}

	currentSet := make(map[string]struct{}, len(current))
	for _, l := range current {
		currentSet[l.String()] = struct{}{}
		if _, ok := baseSet[l.String()]; !ok {
			added = append(added, l)
		}
	}

	for _, l := range base.Listeners {
		if _, ok := currentSet[l.String()]; !ok {
			removed = append(removed, l)
		}
	}
	return
}
