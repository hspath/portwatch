package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func TestDefaults(t *testing.T) {
	cfg := config.Defaults()
	if cfg.BaselineFile != "baseline.json" {
		t.Errorf("expected baseline.json, got %q", cfg.BaselineFile)
	}
	if cfg.ScanInterval.Duration != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.ScanInterval.Duration)
	}
}

func TestLoad_MissingFile_ReturnsDefaults(t *testing.T) {
	cfg, err := config.Load("/nonexistent/portwatch_config.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BaselineFile != "baseline.json" {
		t.Errorf("expected default baseline, got %q", cfg.BaselineFile)
	}
}

func TestLoad_CorruptFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "cfg.json")
	os.WriteFile(p, []byte("{not valid json"), 0600)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "cfg.json")

	orig := config.Config{
		BaselineFile: "/var/lib/portwatch/baseline.json",
		ScanInterval: config.Duration{Duration: 2 * time.Minute},
		AlertLogFile: "/var/log/portwatch.log",
		IgnorePorts:  []uint16{22, 80, 443},
	}

	if err := config.Save(p, orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := config.Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.BaselineFile != orig.BaselineFile {
		t.Errorf("BaselineFile mismatch: %q vs %q", loaded.BaselineFile, orig.BaselineFile)
	}
	if loaded.ScanInterval.Duration != orig.ScanInterval.Duration {
		t.Errorf("ScanInterval mismatch: %v vs %v", loaded.ScanInterval, orig.ScanInterval)
	}
	if loaded.AlertLogFile != orig.AlertLogFile {
		t.Errorf("AlertLogFile mismatch: %q vs %q", loaded.AlertLogFile, orig.AlertLogFile)
	}
	if len(loaded.IgnorePorts) != len(orig.IgnorePorts) {
		t.Errorf("IgnorePorts length mismatch: %v vs %v", loaded.IgnorePorts, orig.IgnorePorts)
	}
}

func TestDuration_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "cfg.json")
	os.WriteFile(p, []byte(`{"scan_interval": "not-a-duration"}`), 0600)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid duration")
	}
}
