package ports

// WellKnownPorts maps port numbers to common service names.
var WellKnownPorts = map[uint16]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	465:  "smtps",
	587:  "submission",
	993:  "imaps",
	995:  "pop3s",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Classification holds the result of classifying a listener.
type Classification struct {
	ServiceName string
	WellKnown   bool
	Privileged  bool
}

// Classify returns a Classification for the given Listener.
func Classify(l Listener) Classification {
	name, known := WellKnownPorts[l.Port]
	if !known {
		name = "unknown"
	}
	return Classification{
		ServiceName: name,
		WellKnown:   known,
		Privileged:  l.Port < 1024,
	}
}

// ClassifyAll returns a map of Listener to Classification for a slice of listeners.
func ClassifyAll(listeners []Listener) map[Listener]Classification {
	result := make(map[Listener]Classification, len(listeners))
	for _, l := range listeners {
		result[l] = Classify(l)
	}
	return result
}
