package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/ports"
)

func makeListener(proto, ip string, port uint16) ports.Listener {
	return ports.Listener{Proto: proto, IP: ip, Port: port}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	b := &baseline.Baseline{
		Listeners: []ports.Listener{
			makeListener("tcp", "0.0.0.0", 22),
			makeListener("tcp", "0.0.0.0", 80),
		},
	}

	if err := baseline.Save(path, b); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Listeners) != 2 {
		t.Errorf("expected 2 listeners, got %d", len(loaded.Listeners))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	b, err := baseline.Load("/nonexistent/path/baseline.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(b.Listeners) != 0 {
		t.Errorf("expected empty baseline, got %v", b.Listeners)
	}
}

func TestLoad_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o644)

	_, err := baseline.Load(path)
	if err == nil {
		t.Fatal("expected error for corrupt JSON")
	}
}

func TestDiff(t *testing.T) {
	base := &baseline.Baseline{
		Listeners: []ports.Listener{
			makeListener("tcp", "0.0.0.0", 22),
			makeListener("tcp", "0.0.0.0", 80),
		},
	}

	current := []ports.Listener{
		makeListener("tcp", "0.0.0.0", 80),
		makeListener("tcp", "0.0.0.0", 443),
	}

	added, removed := baseline.Diff(base, current)

	if len(added) != 1 || added[0].Port != 443 {
		t.Errorf("expected added=[443], got %v", added)
	}
	if len(removed) != 1 || removed[0].Port != 22 {
		t.Errorf("expected removed=[22], got %v", removed)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	listeners := []ports.Listener{makeListener("tcp", "0.0.0.0", 22)}
	base := &baseline.Baseline{Listeners: listeners}

	added, removed := baseline.Diff(base, listeners)
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", added, removed)
	}
}
