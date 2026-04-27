package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	// BaselineFile is the path to the baseline JSON file.
	BaselineFile string `json:"baseline_file"`

	// ScanInterval is how often the daemon rescans open ports.
	ScanInterval Duration `json:"scan_interval"`

	// AlertLogFile is an optional path to write alerts to (empty = stderr).
	AlertLogFile string `json:"alert_log_file,omitempty"`

	// IgnorePorts is a list of ports to silently ignore during diffing.
	IgnorePorts []uint16 `json:"ignore_ports,omitempty"`
}

// Duration is a time.Duration that marshals/unmarshals as a string (e.g. "30s").
type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = parsed
	return nil
}

// Defaults returns a Config populated with sensible defaults.
func Defaults() Config {
	return Config{
		BaselineFile: "baseline.json",
		ScanInterval: Duration{30 * time.Second},
	}
}

// Load reads a JSON config file from path. Missing file returns Defaults.
func Load(path string) (Config, error) {
	cfg := Defaults()
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return Defaults(), err
	}
	return cfg, nil
}

// Save writes cfg as indented JSON to path.
func Save(path string, cfg Config) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}
