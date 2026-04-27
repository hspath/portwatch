package config

import "github.com/user/portwatch/internal/ports"

// FilterIgnored removes any listeners whose port appears in the ignore list.
// It returns a new slice and does not modify the original.
func FilterIgnored(listeners []ports.Listener, ignorePorts []uint16) []ports.Listener {
	if len(ignorePorts) == 0 {
		return listeners
	}

	ignoreSet := make(map[uint16]struct{}, len(ignorePorts))
	for _, p := range ignorePorts {
		ignoreSet[p] = struct{}{}
	}

	filtered := make([]ports.Listener, 0, len(listeners))
	for _, l := range listeners {
		if _, skip := ignoreSet[l.Port]; !skip {
			filtered = append(filtered, l)
		}
	}
	return filtered
}
