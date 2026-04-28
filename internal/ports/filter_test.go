package ports

import (
	"testing"
)

func makeListener(proto, addr string, port uint16) Listener {
	return Listener{Proto: proto, Addr: addr, Port: port}
}

func TestApplyFilter_ByProtocol(t *testing.T) {
	input := []Listener{
		makeListener("tcp", "0.0.0.0", 80),
		makeListener("udp", "0.0.0.0", 53),
		makeListener("tcp", "127.0.0.1", 8080),
	}

	opts := FilterOptions{Protocols: []string{"tcp"}}
	got := ApplyFilter(input, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2 tcp listeners, got %d", len(got))
	}
	for _, l := range got {
		if l.Proto != "tcp" {
			t.Errorf("expected tcp, got %s", l.Proto)
		}
	}
}

func TestApplyFilter_LoopbackOnly(t *testing.T) {
	input := []Listener{
		makeListener("tcp", "0.0.0.0", 80),
		makeListener("tcp", "127.0.0.1", 8080),
		makeListener("udp", "::1", 53),
		makeListener("tcp", "192.168.1.5", 9090),
	}

	opts := FilterOptions{Protocols: []string{"tcp", "udp"}, LoopbackOnly: true}
	got := ApplyFilter(input, opts)
	if len(got) != 3 {
		t.Fatalf("expected 3 loopback listeners, got %d", len(got))
	}
}

func TestApplyFilter_EmptyProtocols_AllowsAll(t *testing.T) {
	input := []Listener{
		makeListener("tcp", "0.0.0.0", 80),
		makeListener("udp", "0.0.0.0", 53),
	}

	opts := FilterOptions{Protocols: []string{}}
	got := ApplyFilter(input, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2 listeners, got %d", len(got))
	}
}

func TestApplyFilter_NoMatch(t *testing.T) {
	input := []Listener{
		makeListener("tcp", "0.0.0.0", 80),
	}

	opts := FilterOptions{Protocols: []string{"udp"}}
	got := ApplyFilter(input, opts)
	if len(got) != 0 {
		t.Fatalf("expected 0 listeners, got %d", len(got))
	}
}

func TestDefaultFilterOptions(t *testing.T) {
	opts := DefaultFilterOptions()
	if len(opts.Protocols) != 2 {
		t.Errorf("expected 2 default protocols, got %d", len(opts.Protocols))
	}
	if opts.LoopbackOnly {
		t.Error("expected LoopbackOnly to be false by default")
	}
}

func TestIsLoopback(t *testing.T) {
	cases := []struct {
		addr string
		want bool
	}{
		{"127.0.0.1", true},
		{"127.0.0.53", true},
		{"::1", true},
		{"0.0.0.0", true},
		{"::", true},
		{"192.168.1.1", false},
		{"10.0.0.1", false},
	}
	for _, tc := range cases {
		got := isLoopback(tc.addr)
		if got != tc.want {
			t.Errorf("isLoopback(%q) = %v, want %v", tc.addr, got, tc.want)
		}
	}
}
