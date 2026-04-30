package ports

import (
	"testing"
)

func enrichListener(ip string, port uint16, proto string) Listener {
	return Listener{IP: ip, Port: port, Protocol: proto}
}

func TestEnrich_KnownService(t *testing.T) {
	l := enrichListener("0.0.0.0", 22, "tcp")
	el := Enrich(l, true, DefaultEnricherOptions())
	if el.ServiceName != "ssh" {
		t.Errorf("expected ssh, got %q", el.ServiceName)
	}
}

func TestEnrich_UnknownService(t *testing.T) {
	l := enrichListener("0.0.0.0", 9999, "tcp")
	el := Enrich(l, true, DefaultEnricherOptions())
	if el.ServiceName == "" {
		t.Error("expected non-empty service name")
	}
}

func TestEnrich_SeveritySet(t *testing.T) {
	l := enrichListener("0.0.0.0", 80, "tcp")
	el := Enrich(l, true, DefaultEnricherOptions())
	if el.Severity == 0 {
		t.Error("expected non-zero severity")
	}
}

func TestEnrich_ClassificationSet(t *testing.T) {
	l := enrichListener("0.0.0.0", 443, "tcp")
	el := Enrich(l, true, DefaultEnricherOptions())
	if el.Classification == "" {
		t.Error("expected non-empty classification")
	}
}

func TestEnrich_NoDNS_HostnameEmpty(t *testing.T) {
	l := enrichListener("127.0.0.1", 8080, "tcp")
	opts := DefaultEnricherOptions()
	opts.ResolveDNS = false
	el := Enrich(l, true, opts)
	if el.Hostname != "" {
		t.Errorf("expected empty hostname without DNS, got %q", el.Hostname)
	}
}

func TestEnrichAll_ReturnsCorrectLength(t *testing.T) {
	listeners := []Listener{
		enrichListener("0.0.0.0", 22, "tcp"),
		enrichListener("0.0.0.0", 80, "tcp"),
		enrichListener("0.0.0.0", 443, "tcp"),
	}
	result := EnrichAll(listeners, true, DefaultEnricherOptions())
	if len(result) != len(listeners) {
		t.Errorf("expected %d enriched listeners, got %d", len(listeners), len(result))
	}
}

func TestDefaultEnricherOptions_DNSDisabled(t *testing.T) {
	opts := DefaultEnricherOptions()
	if opts.ResolveDNS {
		t.Error("expected ResolveDNS to be false by default")
	}
}

func TestKnownService_AllWellKnown(t *testing.T) {
	cases := map[uint16]string{
		22:    "ssh",
		80:    "http",
		443:   "https",
		3306:  "mysql",
		5432:  "postgres",
		6379:  "redis",
		27017: "mongodb",
	}
	for port, want := range cases {
		got := knownService(port)
		if got != want {
			t.Errorf("port %d: expected %q, got %q", port, want, got)
		}
	}
}
