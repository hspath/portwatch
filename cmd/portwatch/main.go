// portwatch is a lightweight CLI daemon that monitors open ports
// and alerts on unexpected listeners by comparing against a saved baseline.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/ports"
)

const defaultBaselineFile = "/var/lib/portwatch/baseline.json"
const defaultInterval = 30 * time.Second

func main() {
	var (
		baselineFile string
		learn        bool
		once         bool
		interval     time.Duration
	)

	flag.StringVar(&baselineFile, "baseline", defaultBaselineFile, "path to baseline file")
	flag.BoolVar(&learn, "learn", false, "scan current listeners and save as new baseline, then exit")
	flag.BoolVar(&once, "once", false, "run a single scan against the baseline and exit")
	flag.DurationVar(&interval, "interval", defaultInterval, "how often to scan when running as a daemon")
	flag.Parse()

	alerter := alert.New(os.Stdout)

	if learn {
		if err := runLearn(baselineFile); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if once {
		exitCode := runScan(baselineFile, alerter)
		os.Exit(exitCode)
	}

	// Daemon mode: scan on a ticker until interrupted.
	fmt.Fprintf(os.Stdout, "portwatch starting (interval=%s, baseline=%s)\n", interval, baselineFile)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run an immediate scan before waiting for the first tick.
	runScan(baselineFile, alerter)

	for range ticker.C {
		runScan(baselineFile, alerter)
	}
}

// runLearn scans the current listeners and writes them to the baseline file.
func runLearn(baselineFile string) error {
	listeners, err := ports.ScanListeners()
	if err != nil {
		return fmt.Errorf("scanning listeners: %w", err)
	}

	if err := baseline.Save(baselineFile, listeners); err != nil {
		return fmt.Errorf("saving baseline to %s: %w", baselineFile, err)
	}

	fmt.Printf("baseline saved: %d listeners recorded to %s\n", len(listeners), baselineFile)
	return nil
}

// runScan compares current listeners against the baseline and emits alerts.
// It returns 0 if no unexpected listeners are found, 1 otherwise.
func runScan(baselineFile string, alerter *alert.Alerter) int {
	current, err := ports.ScanListeners()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error scanning listeners: %v\n", err)
		return 1
	}

	saved, err := baseline.Load(baselineFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading baseline from %s: %v\n", baselineFile, err)
		fmt.Fprintf(os.Stderr, "hint: run with -learn to create an initial baseline\n")
		return 1
	}

	added, removed := baseline.Diff(saved, current)

	for _, l := range added {
		alerter.Unexpected(l)
	}
	for _, l := range removed {
		alerter.Gone(l)
	}

	if len(added) > 0 {
		return 1
	}
	return 0
}
