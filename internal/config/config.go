package config

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/user/portwatch/internal/notify"
)

// Config holds all portwatch runtime configuration.
type Config struct {
	BaselineFile string        `json:"baseline_file"`
	Interval     Duration      `json:"interval"`
	IgnorePorts  []int         `json:"ignore_ports"`
	Notify       notify.Config `json:"notify"`
}

// Duration wraps time.Duration for JSON marshal/unmarshal as a string.
type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = v
	return nil
}

// Defaults returns a Config populated with sensible defaults.
func Defaults() Config {
	return Config{
		BaselineFile: "baseline.json",
		Interval:     Duration{30 * time.Second},
		IgnorePorts:  []int{},
		Notify:       notify.Config{Method: notify.MethodStdout},
	}
}

// Load reads a Config from path. Returns Defaults if the file does not exist.
func Load(path string) (Config, error) {
	cfg := Defaults()
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// Save writes cfg as JSON to path, creating or truncating the file.
func Save(path string, cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// FilterIgnored removes ports listed in cfg.IgnorePorts from ports.
func FilterIgnored(ports []int, ignored []int) []int {
	set := make(map[int]struct{}, len(ignored))
	for _, p := range ignored {
		set[p] = struct{}{}
	}
	out := ports[:0]
	for _, p := range ports {
		if _, skip := set[p]; !skip {
			out = append(out, p)
		}
	}
	return out
}
