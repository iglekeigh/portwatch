package scanner

import (
	"net"
	"testing"
	"time"
)

// startTestServer starts a TCP listener on a random port and returns the port.
func startTestServer(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return ln.Addr().(*net.TCPAddr).Port
}

func TestScan_OpenPort(t *testing.T) {
	port := startTestServer(t)
	s := New(ScanOptions{
		Host:    "127.0.0.1",
		Ports:   []int{port},
		Timeout: time.Second,
	})
	open, err := s.Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(open) != 1 || open[0] != port {
		t.Errorf("expected port %d to be open, got %v", port, open)
	}
}

func TestScan_ClosedPort(t *testing.T) {
	s := New(ScanOptions{
		Host:    "127.0.0.1",
		Ports:   []int{1},
		Timeout: 200 * time.Millisecond,
	})
	open, err := s.Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(open) != 0 {
		t.Errorf("expected no open ports, got %v", open)
	}
}

func TestScan_DefaultTimeout(t *testing.T) {
	s := New(ScanOptions{Host: "127.0.0.1", Ports: []int{}})
	if s.opts.Timeout != 2*time.Second {
		t.Errorf("expected default timeout 2s, got %v", s.opts.Timeout)
	}
}
