package config

import "github.com/user/portwatch/internal/ports"

// IgnoreRule describes a single port+protocol combination to suppress from alerts.
type IgnoreRule struct {
	Proto string `json:"proto"` // "tcp" or "udp"
	Port  uint16 `json:"port"`
}

// FilterIgnored removes listeners that match any ignore rule in the config.
func FilterIgnored(listeners []ports.Listener, rules []IgnoreRule) []ports.Listener {
	if len(rules) == 0 {
		return listeners
	}

	type key struct {
		proto string
		port  uint16
	}
	ignored := make(map[key]struct{}, len(rules))
	for _, r := range rules {
		ignored[key{proto: r.Proto, port: r.Port}] = struct{}{}
	}

	var result []ports.Listener
	for _, l := range listeners {
		if _, skip := ignored[key{proto: l.Proto, port: l.Port}]; !skip {
			result = append(result, l)
		}
	}
	return result
}
