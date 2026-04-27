package ports

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Listener represents an open TCP/UDP port with its associated process info.
type Listener struct {
	Protocol string
	Address  string
	Port     uint16
	PID      int
	Process  string
}

// String returns a human-readable representation of a Listener.
func (l Listener) String() string {
	return fmt.Sprintf("%s %s:%d (pid=%d, process=%s)", l.Protocol, l.Address, l.Port, l.PID, l.Process)
}

// ScanListeners reads /proc/net/tcp and /proc/net/tcp6 to return all
// currently listening TCP ports on Linux.
func ScanListeners() ([]Listener, error) {
	var listeners []Listener

	for _, path := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		results, err := parseProcNet(path, "tcp")
		if err != nil {
			// file may not exist on all systems; skip gracefully
			continue
		}
		listeners = append(listeners, results...)
	}

	return listeners, nil
}

// parseProcNet parses a /proc/net/tcp or /proc/net/tcp6 file.
func parseProcNet(path, proto string) ([]Listener, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var listeners []Listener
	scanner := bufio.NewScanner(f)

	// skip header line
	scanner.Scan()

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}

		// state 0A = TCP_LISTEN
		if fields[3] != "0A" {
			continue
		}

		addr, port, err := parseHexAddr(fields[1])
		if err != nil {
			continue
		}

		listeners = append(listeners, Listener{
			Protocol: proto,
			Address:  addr,
			Port:     port,
		})
	}

	return listeners, scanner.Err()
}

// parseHexAddr converts a hex-encoded "address:port" field from /proc/net/tcp.
func parseHexAddr(hexAddr string) (string, uint16, error) {
	parts := strings.SplitN(hexAddr, ":", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid address field: %s", hexAddr)
	}

	portVal, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return "", 0, err
	}

	return parts[0], uint16(portVal), nil
}
