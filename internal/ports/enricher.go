package ports

import (
	"fmt"
	"net"
)

// EnrichedListener wraps a Listener with additional metadata.
type EnrichedListener struct {
	Listener
	Hostname  string
	ServiceName string
	Severity  Severity
	Classification Classification
}

// EnricherOptions controls enrichment behaviour.
type EnricherOptions struct {
	// ResolveDNS enables reverse-DNS hostname resolution.
	ResolveDNS bool
}

// DefaultEnricherOptions returns sensible defaults.
func DefaultEnricherOptions() EnricherOptions {
	return EnricherOptions{
		ResolveDNS: false,
	}
}

// Enrich adds metadata to a Listener.
func Enrich(l Listener, added bool, opts EnricherOptions) EnrichedListener {
	el := EnrichedListener{
		Listener:       l,
		Classification: Classify(l),
		Severity:       SeverityFor(l, added),
	}

	if opts.ResolveDNS {
		el.Hostname = resolveHostname(l.IP)
	}

	el.ServiceName = knownService(l.Port)
	return el
}

// EnrichAll enriches a slice of listeners.
func EnrichAll(listeners []Listener, added bool, opts EnricherOptions) []EnrichedListener {
	out := make([]EnrichedListener, len(listeners))
	for i, l := range listeners {
		out[i] = Enrich(l, added, opts)
	}
	return out
}

func resolveHostname(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		return ip
	}
	return names[0]
}

func knownService(port uint16) string {
	switch port {
	case 22:
		return "ssh"
	case 80:
		return "http"
	case 443:
		return "https"
	case 3306:
		return "mysql"
	case 5432:
		return "postgres"
	case 6379:
		return "redis"
	case 27017:
		return "mongodb"
	default:
		return fmt.Sprintf("port/%d", port)
	}
}
