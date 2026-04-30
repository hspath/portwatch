package ports

// Severity represents the urgency level of a port event.
type Severity int

const (
	SeverityInfo    Severity = iota // expected or baseline-matching change
	SeverityWarning                 // unexpected listener on a high port
	SeverityCritical                // unexpected listener on a privileged or well-known port
)

// String returns a human-readable label for the severity.
func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// SeverityFor derives the appropriate Severity for a Listener based on its
// classification and whether the event represents an addition or removal.
//
// Rules:
//   - Removed listeners are always INFO (expected churn).
//   - Added well-known or privileged-unknown ports → CRITICAL.
//   - Added high/unknown ports → WARNING.
func SeverityFor(l Listener, added bool) Severity {
	if !added {
		return SeverityInfo
	}

	switch Classify(l) {
	case ClassWellKnown, ClassPrivilegedUnknown:
		return SeverityCritical
	default:
		return SeverityWarning
	}
}
