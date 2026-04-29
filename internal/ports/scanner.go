package ports

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
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
// The address portion is stored in little-endian byte order by the kernel and
// is decoded into a human-readable IP string.
func parseHexAddr(hexAddr string) (string, uint16, error) {
	parts := strings.SplitN(hexAddr, ":", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid address field: %s", hexAddr)
	}

	portVal, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return "", 0, err
	}

	ip, err := hexToIP(parts[0])
	if err != nil {
		return parts[0], uint16(portVal), nil
	}

	return ip, uint16(portVal), nil
}

// hexToIP decodes a hex-encoded, little-endian IP address as found in
// /proc/net/tcp (4 bytes for IPv4, 16 bytes for IPv6) into a string.
func hexToIP(hexIP string) (string, error) {
	b, err := hex.DecodeString(hexIP)
	if err != nil {
		return "", err
	}
	switch len(b) {
	case 4:
		// IPv4 is stored in little-endian order
		v := binary.LittleEndian.Uint32(b)
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, v)
		return ip.String(), nil
	case 16:
		// IPv6 is stored as four little-endian 32-bit words
		ip := make(net.IP, 16)
		for i := 0; i < 4; i++ {
			v := binary.LittleEndian.Uint32(b[i*4 : i*4+4])
			binary.BigEndian.PutUint32(ip[i*4:], v)
		}
		return ip.String(), nil
	default:
		return "", fmt.Errorf("unexpected IP length: %d bytes", len(b))
	}
}
