package ports

import "strings"

// Listener represents an open port on the system.
type FilterOptions struct {
	Protocols []string // e.g. ["tcp", "udp"]
	LoopbackOnly bool
}

// DefaultFilterOptions returns sensible defaults: both tcp and udp, all addresses.
func DefaultFilterOptions() FilterOptions {
	return FilterOptions{
		Protocols:    []string{"tcp", "udp"},
		LoopbackOnly: false,
	}
}

// ApplyFilter returns only the listeners that match the given options.
func ApplyFilter(listeners []Listener, opts FilterOptions) []Listener {
	protoSet := make(map[string]struct{}, len(opts.Protocols))
	for _, p := range opts.Protocols {
		protoSet[strings.ToLower(p)] = struct{}{}
	}

	var result []Listener
	for _, l := range listeners {
		if len(protoSet) > 0 {
			if _, ok := protoSet[strings.ToLower(l.Proto)]; !ok {
				continue
			}
		}
		if opts.LoopbackOnly && !isLoopback(l.Addr) {
			continue
		}
		result = append(result, l)
	}
	return result
}

// isLoopback returns true if the address is a loopback address.
func isLoopback(addr string) bool {
	return strings.HasPrefix(addr, "127.") ||
		addr == "::1" ||
		addr == "0.0.0.0" ||
		addr == "::"
}
